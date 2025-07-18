package logger

import (
	"bytes"
	"unsafe"
)

type StringBuffer struct {
	data []byte
}

func NewStringBuffer(initCap int) StringBuffer {
	return StringBuffer{
		data: make([]byte, 0, initCap),
	}
}

func (buf *StringBuffer) Reset() {
	buf.data = buf.data[:0]
}

func (buf *StringBuffer) String() string {
	return unsafe.String(unsafe.SliceData(buf.data), len(buf.data))
}

func (buf *StringBuffer) Bytes() []byte {
	return buf.data
}

func (buf *StringBuffer) Write(p []byte) (n int, err error) {
	buf.data = append(buf.data, p...)
	return len(p), nil
}

func (buf *StringBuffer) WriteString(str string) {
	buf.data = append(buf.data, str...)
}

func (buf *StringBuffer) WriteByte(b byte) error {
	buf.data = append(buf.data, b)
	return nil
}

func (buf *StringBuffer) Swap(other *StringBuffer) {
	temp := buf.data
	buf.data = other.data
	other.data = temp
}

func (buf StringBuffer) Equals(other StringBuffer) bool {
	return bytes.Equal(buf.data, other.data)
}

func (buf *StringBuffer) Len() int { return len(buf.data) }

func (buf *StringBuffer) Cap() int { return cap(buf.data) }

func (buf *StringBuffer) EnsureSpace(space int) {
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

var _ Writer = (*StringBuffer)(nil)
var _ Stringer = (*StringBuffer)(nil)
