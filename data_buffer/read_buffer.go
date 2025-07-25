package data_buffer

import (
	"io"
	"unsafe"
)

type ReadBuffer struct {
	data []byte
}

func NewReadBuffer(data []byte) ReadBuffer {
	return ReadBuffer{
		data: data,
	}
}

func (r *ReadBuffer) Read(p []byte) (n int, err error) {
	n = copy(p[:], r.data[:])
	if n <= 0 {
		return 0, io.EOF
	}
	r.data = r.data[n:]
	return n, nil
}

func (r *ReadBuffer) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(r.data)
	r.data = r.data[nn:]
	return int64(nn), err
}

func (r *ReadBuffer) Len() int {
	return len(r.data)
}

func (r *ReadBuffer) BytesRef() []byte {
	return r.data
}

func (r *ReadBuffer) BytesCopy() []byte {
	out := make([]byte, len(r.data))
	copy(out, r.data)
	return out
}

// CAUTION: This uses [unsafe] to reinterpret the underlying data
// slice as a [string], but strings in Go are expected to be immutable.
// Any changes to the underlying byte slice will invalidate the string
func (r *ReadBuffer) StringRef() string {
	return unsafe.String(unsafe.SliceData(r.data), len(r.data))
}

func (r *ReadBuffer) StringCopy() string {
	return string(r.data)
}

var _ io.Reader = (*ReadBuffer)(nil)
var _ io.WriterTo = (*ReadBuffer)(nil)

type ReadErrBuffer struct {
	data []byte
	err  error
}

func NewReadErrBuffer(data []byte) ReadErrBuffer {
	return ReadErrBuffer{
		data: data,
	}
}
func (r *ReadErrBuffer) Read(p []byte) (n int, err error) {
	n = copy(p[:], r.data[:])
	if n <= 0 {
		r.err = io.EOF
		return 0, io.EOF
	}
	r.data = r.data[n:]
	return n, nil
}

func (r *ReadErrBuffer) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(r.data)
	if r.err == nil {
		r.err = err
	}
	r.data = r.data[nn:]
	return int64(nn), err
}

func (r *ReadErrBuffer) Len() int {
	return len(r.data)
}

func (r *ReadErrBuffer) BytesRef() []byte {
	return r.data
}

func (r *ReadErrBuffer) BytesCopy() []byte {
	out := make([]byte, len(r.data))
	copy(out, r.data)
	return out
}

// CAUTION: This uses [unsafe] to reinterpret the underlying data
// slice as a [string], but strings in Go are expected to be immutable.
// Any changes to the underlying byte slice will invalidate the string
func (r *ReadErrBuffer) StringRef() string {
	return unsafe.String(unsafe.SliceData(r.data), len(r.data))
}

func (r *ReadErrBuffer) StringCopy() string {
	return string(r.data)
}

var _ io.Reader = (*ReadErrBuffer)(nil)
var _ io.WriterTo = (*ReadErrBuffer)(nil)
