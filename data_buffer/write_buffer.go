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

const ()

type WriteBuffer struct {
	data []byte
}

func NewWriteBuffer(initCap int) WriteBuffer {
	return WriteBuffer{
		data: make([]byte, 0, initCap),
	}
}

func (buf *WriteBuffer) Reset() {
	buf.data = buf.data[:0]
}

func (buf *WriteBuffer) StringRef() string {
	return unsafe.String(unsafe.SliceData(buf.data), len(buf.data))
}

func (buf *WriteBuffer) StringCopy() string {
	return string(buf.data)
}

func (buf *WriteBuffer) BytesRef() []byte {
	return buf.data
}

func (buf *WriteBuffer) BytesCopy() []byte {
	newBuf := make([]byte, buf.Len())
	copy(newBuf, buf.BytesRef())
	return newBuf
}

func (buf *WriteBuffer) Write(p []byte) (n int, err error) {
	buf.data = append(buf.data, p...)
	return len(p), nil
}

func (buf *WriteBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	newLen := len(buf.data)
	nn, err := r.Read(buf.data[newLen:cap(buf.data)])
	newLen += nn
	buf.data = buf.data[:newLen]
	if err == nil && newLen == cap(buf.data) {
		var rem []byte
		rem, err = io.ReadAll(r)
		buf.data = append(buf.data, rem...)
		nn += len(rem)
	}
	return int64(nn), err
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

func (buf *WriteBuffer) WriteByte(b byte) {
	buf.data = append(buf.data, b)
}

func (buf *WriteBuffer) WriteBytes(b ...byte) {
	buf.data = append(buf.data, b...)
}

func (buf *WriteBuffer) WriteRune(r rune) {
	var arr [4]byte
	n := utf8.EncodeRune(arr[:], r)
	buf.data = append(buf.data, arr[:n]...)
}
func (buf *WriteBuffer) WriteRunes(runes []rune) {
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
		case byte:
			buf.WriteByte(vv)
		case *byte:
			if vv == nil {
				return
			}
			buf.WriteByte(*vv)
		case []byte:
			buf.WriteBytes(vv...)
		case [][]byte:
			for _, vvv := range vv {
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
			buf.WriteRunes(vv)
		case [][]rune:
			for _, vvv := range vv {
				buf.WriteRunes(vvv)
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

func (buf *WriteBuffer) EnsureSpace(space int) {
	curr := cap(buf.data) - len(buf.data)
	need := space - curr
	if need <= 0 {
		return
	}
	newTotal := cap(buf.data) + need
	newData := append([]byte(nil), make([]byte, newTotal)...)
	i := copy(newData, buf.data)
	buf.data = newData[:i]
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
