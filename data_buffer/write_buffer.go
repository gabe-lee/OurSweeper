package data_buffer

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sync"
	"unicode/utf8"
	"unsafe"

	"github.com/gabe-lee/OurSweeper/lock"
)

type (
	//Implements lock.Locker, but is defined here as a concrete type
	Lock = lock.MiniLock
)

type (
	Writer   = io.Writer
	Stringer = fmt.Stringer
	Locker   = sync.Locker
)

const (
	smallestGrow = 64
)

type WriteBuffer struct {
	data []byte
}

func (buf *WriteBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	nn := len(buf.data)
	buf.EnsureSpaceFrom(nil, int(off), n)
	endOff := off + int64(nn)
	n = copy(buf.data[off:endOff], p)
	return
}

func NewWriteBuffer(initCap int) WriteBuffer {
	return WriteBuffer{
		data: make([]byte, 0, initCap),
	}
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
		buf.EnsureSpace(nil, smallestGrow)
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

func (buf *WriteBuffer) ReadFromWithMax(r io.Reader, max int64) (n int64, err error) {
	origLen := len(buf.data)
	rem := cap(buf.data) - origLen
	truemax := min(int(max), rem)
	newLen := origLen
	var nn int
	nn, err = r.Read(buf.data[newLen : newLen+truemax])
	newLen += nn
	buf.data = buf.data[:newLen]
	for err == nil && nn > 0 {
		buf.EnsureSpace(nil, smallestGrow)
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

func (buf *WriteBuffer) ReadFromUntilByte(r io.Reader, terminalByte byte) (n int, err error) {
	var nextByteArr [1]byte
	var nn int
	oldLen := len(buf.data)
	nn, err = r.Read(nextByteArr[:])
	for nn > 0 && err == nil && nextByteArr[0] != terminalByte {
		buf.data = append(buf.data, nextByteArr[0])
	}
	if nextByteArr[0] == terminalByte {
		buf.data = append(buf.data, nextByteArr[0])
	}
	n = len(buf.data) - oldLen
	return
}

func (buf WriteBuffer) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(buf.data)
	return int64(nn), err
}

func (buf *WriteBuffer) WriteString(str string) {
	buf.data = append(buf.data, str...)
}

func (buf *WriteBuffer) WriteStrings(strs ...string) {
	for _, str := range strs {
		buf.data = append(buf.data, str...)
	}
}

func (buf *WriteBuffer) WriteSlices(slices ...[]byte) {
	for _, slice := range slices {
		buf.data = append(buf.data, slice...)
	}
}

func (buf *WriteBuffer) WriteByte(b byte) error {
	buf.data = append(buf.data, b)
	return nil
}

func (buf *WriteBuffer) WriteBytes(b ...byte) {
	buf.data = append(buf.data, b...)
}

func (buf *WriteBuffer) WriteRune(r rune) {
	var arr [4]byte
	n := utf8.EncodeRune(arr[:], r)
	buf.data = append(buf.data, arr[:n]...)
}

func (buf *WriteBuffer) WriteRunes(runes ...rune) {
	var arr [4]byte
	for _, r := range runes {
		n := utf8.EncodeRune(arr[:], r)
		buf.data = append(buf.data, arr[:n]...)
	}
}

func (buf *WriteBuffer) WriteRuneSlice(runes []rune) {
	var arr [4]byte
	for _, r := range runes {
		n := utf8.EncodeRune(arr[:], r)
		buf.data = append(buf.data, arr[:n]...)
	}
}

func (buf *WriteBuffer) WriteAny(vals ...any) {
	for _, v := range vals {
		switch vv := v.(type) {
		case string:
			buf.WriteString(vv)
		case *string:
			if vv == nil {
				return
			}
			buf.WriteString(*vv)
		case []string:
			for _, vvv := range vv {
				buf.WriteString(vvv)
			}
		case *[]string:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.WriteString(vvv)
			}
		case byte:
			buf.WriteByte(vv)
		case *byte:
			if vv == nil {
				return
			}
			buf.WriteByte(*vv)
		case []byte:
			buf.WriteBytes(vv...)
		case *[]byte:
			if vv == nil {
				return
			}
			buf.WriteBytes(*vv...)
		case [][]byte:
			for _, vvv := range vv {
				buf.WriteBytes(vvv...)
			}
		case *[][]byte:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.WriteBytes(vvv...)
			}
		case rune:
			buf.WriteRune(vv)
		case *rune:
			if vv == nil {
				return
			}
			buf.WriteRune(*vv)
		case []rune:
			buf.WriteRuneSlice(vv)
		case *[]rune:
			if vv == nil {
				return
			}
			buf.WriteRuneSlice(*vv)
		case [][]rune:
			for _, vvv := range vv {
				buf.WriteRuneSlice(vvv)
			}
		case *[][]rune:
			if vv == nil {
				return
			}
			for _, vvv := range *vv {
				buf.WriteRuneSlice(vvv)
			}
		default:
			panic(fmt.Sprintf("invalid type `%s` for `StringBuffer.WriteAny()`", reflect.TypeOf(v).Name()))
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

func (buf WriteBuffer) Equals(other WriteBuffer) bool {
	return bytes.Equal(buf.data, other.data)
}

func (buf *WriteBuffer) Len() int { return len(buf.data) }

func (buf *WriteBuffer) Cap() int { return cap(buf.data) }

func (buf *WriteBuffer) EnsureSpace(pool *WriteBufferPool, space int) {
	buf.EnsureSpaceFrom(pool, len(buf.data), space)
}

func (buf *WriteBuffer) EnsureSpaceFrom(pool *WriteBufferPool, offset int, space int) {
	newSize := int(offset) + space
	currSize := cap(buf.data)
	if newSize <= currSize {
		return
	}
	if pool != nil {
		currClass := getSizeClass(currSize)
		newClass := getSizeClass(newSize)
		if newClass > currClass {
			swapBuf := pool.getBufferByClassInternal(newClass)
			copy(swapBuf.data, buf.data)
			buf.Swap(swapBuf)
			pool.ReleaseBuffer(swapBuf)
			return
		}
	}
	newData := append([]byte(nil), make([]byte, newSize)...)
	i := copy(newData, buf.data)
	oldData := buf.data
	buf.data = newData[:i]
	if pool != nil {
		newBuf := WriteBuffer{}
		newBuf.data = oldData
		pool.ReleaseBuffer(&newBuf)
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
