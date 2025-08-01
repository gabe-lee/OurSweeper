package data_buffer

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type WriteBuffer struct {
	noCopy

	data []byte
	pool *WriteBufferPool
}

func NewWriteBuffer(initCap int) WriteBuffer {
	return WriteBuffer{
		data: make([]byte, 0, initCap),
	}
}

func NewWriteBufferFilled(data []byte) WriteBuffer {
	return WriteBuffer{
		data: data,
	}
}

func (buf *WriteBuffer) WriteString(s string) (n int, err error) {
	buf.String(s)
	return len(s), nil
}

func (buf *WriteBuffer) Close() error {
	buf.Reset()
	if buf.pool != nil {
		buf.pool.ReleaseBuffer(buf)
	}
	return nil
}

func (buf *WriteBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	nn := len(buf.data)
	buf.EnsureSpaceFrom(int(off), n)
	endOff := off + int64(nn)
	n = copy(buf.data[off:endOff], p)
	return
}

func (buf *WriteBuffer) Reset() {
	buf.data = buf.data[:0]
}

// CAUTION: This uses [unsafe] to reinterpret the underlying data
// slice as a [string], but strings in Go are expected to be immutable.
// Any changes to the underlying byte slice will invalidate the string
func (buf *WriteBuffer) StringRef() string {
	return unsafe.String(unsafe.SliceData(buf.data), len(buf.data))
}

func (buf *WriteBuffer) StringCopy() string {
	return string(buf.data)
}

func (buf *WriteBuffer) SetData(data []byte) {
	buf.data = data
}

func (buf *WriteBuffer) SetLenRelative(offset int) {
	if offset > 0 {
		buf.EnsureSpace(offset)
	}
	buf.data = buf.data[:len(buf.data)+offset]
}

func (buf *WriteBuffer) BytesRef() []byte {
	return buf.data
}

func (buf *WriteBuffer) BytesCopy() []byte {
	newBuf := make([]byte, buf.Len())
	copy(newBuf, buf.BytesRef())
	return newBuf
}

// Never returns non-nil error
func (buf *WriteBuffer) Write(p []byte) (n int, err error) {
	buf.EnsureSpace(len(p))
	buf.data = append(buf.data, p...)
	return len(p), nil
}

func (buf *WriteBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	origLen := len(buf.data)
	newLen := origLen
	var nn int
	nn, err = r.Read(buf.data[newLen:cap(buf.data)])
	newLen += nn
	buf.data = buf.data[:newLen]
	for err == nil && nn > 0 {
		buf.EnsureSpace(smallestGrow)
		nn, err = r.Read(buf.data[newLen:cap(buf.data)])
		newLen += nn
		buf.data = buf.data[:newLen]
	}
	if err == io.EOF {
		err = nil
	}
	n = int64(newLen) - int64(origLen)
	return
}

func (buf *WriteBuffer) ReadNFrom(r io.Reader, n int) (nn int, err error) {
	start := buf.Len()
	buf.Skip(n)
	end := buf.Len()
	nn, err = r.Read(buf.data[start:end])
	if n != nn {
		err = io.ErrShortBuffer
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func (buf *WriteBuffer) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(buf.data)
	return int64(nn), err
}

func (buf *WriteBuffer) WriteByte(b byte) error {
	buf.data = append(buf.data, b)
	return nil
}

func (buf *WriteBuffer) WriteBytes(b ...byte) error {
	buf.data = append(buf.data, b...)
	return nil
}

func (buf *WriteBuffer) WriteText(vals ...any) {
	for _, v := range vals {
		switch vv := v.(type) {
		case string:
			buf.String(vv)
		case *string:
			if vv == nil {
				return
			}
			buf.String(*vv)
		case []string:
			for _, vvv := range vv {
				buf.String(vvv)
			}
		case *[]string:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.String(vvv)
			}
		case byte:
			buf.U8(vv)
		case *byte:
			if vv == nil {
				return
			}
			buf.U8(*vv)
		case []byte:
			buf.U8_Slice(vv)
		case *[]byte:
			if vv == nil {
				return
			}
			buf.U8_Slice(*vv)
		case [][]byte:
			for _, vvv := range vv {
				buf.U8_Slice(vvv)
			}
		case *[][]byte:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.U8_Slice(vvv)
			}
		case rune:
			buf.Rune(vv)
		case *rune:
			if vv == nil {
				return
			}
			buf.Rune(*vv)
		case []rune:
			buf.Rune_Slice(vv)
		case *[]rune:
			if vv == nil {
				return
			}
			buf.Rune_Slice(*vv)
		case [][]rune:
			for _, vvv := range vv {
				buf.Rune_Slice(vvv)
			}
		case *[][]rune:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.Rune_Slice(vvv)
			}
		default:
			panic(fmt.Sprintf("invalid type `%s` for `StringBuffer.WriteText()`", reflect.TypeOf(v).Name()))
		}
	}
}

func (buf *WriteBuffer) TrimEndWhitespace() {
	i := len(buf.data)
	if i == 0 {
		return
	}
	for i > 0 && buf.data[i-1] <= byte(' ') {
		i -= 1
	}
	buf.data = buf.data[:i]
}

func (buf *WriteBuffer) Swap(other *WriteBuffer) {
	temp := buf.data
	buf.data = other.data
	other.data = temp
}

func (buf *WriteBuffer) Equals(other *WriteBuffer) bool {
	return bytes.Equal(buf.data, other.data)
}

func (buf *WriteBuffer) Len() int { return len(buf.data) }

func (buf *WriteBuffer) Pos() int { return len(buf.data) }

func (buf *WriteBuffer) Cap() int { return cap(buf.data) }

func (buf *WriteBuffer) EnsureSpace(space int) {
	buf.EnsureSpaceFrom(len(buf.data), space)
}

func (buf *WriteBuffer) EnsureSpaceFrom(offset int, space int) {
	newSize := int(offset) + space
	currSize := cap(buf.data)
	if newSize <= currSize {
		return
	}
	if buf.pool != nil {
		currClass := getSizeClass(currSize)
		newClass := getSizeClass(newSize)
		if newClass > currClass {
			swapBuf := buf.pool.getBufferByClassInternal(newClass)
			copy(swapBuf.data, buf.data)
			buf.Swap(swapBuf)
			buf.pool.ReleaseBuffer(swapBuf)
			return
		}
	}
	oldData := buf.data
	newData := append([]byte(nil), make([]byte, newSize)...)
	i := copy(newData, oldData)
	buf.data = newData[:i]
	if buf.pool != nil {
		newBuf := WriteBuffer{
			data: oldData,
			pool: buf.pool,
		}
		buf.pool.ReleaseBuffer(&newBuf)
	}
}

func (buf *WriteBuffer) ReaderRef() ReadBuffer {
	return ReadBuffer{
		data: buf.data,
	}
}

func (buf *WriteBuffer) ReaderCopy() ReadBuffer {
	b := ReadBuffer{
		data: make([]byte, len(buf.data)),
	}
	copy(b.data, buf.data)
	return b
}

var _ io.Writer = (*WriteBuffer)(nil)
var _ io.WriterTo = (*WriteBuffer)(nil)
var _ io.ReaderFrom = (*WriteBuffer)(nil)
var _ io.WriterAt = (*WriteBuffer)(nil)
var _ io.ByteWriter = (*WriteBuffer)(nil)
var _ io.Closer = (*WriteBuffer)(nil)
var _ io.StringWriter = (*WriteBuffer)(nil)
