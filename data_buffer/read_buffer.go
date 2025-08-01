package data_buffer

import (
	"io"
	"unicode/utf8"
	"unsafe"
)

type ReadBuffer struct {
	data []byte
	pos  int
}

func (buf *ReadBuffer) Pos() int { return buf.pos }

func (r *ReadBuffer) Close() error {
	r.data = nil
	r.pos = 0
	return nil
}

func (r *ReadBuffer) ReadAt(p []byte, off int64) (n int, err error) {
	if off > int64(len(r.data)) {
		return 0, io.ErrShortBuffer
	}
	subSlice := r.data[off:]
	n = copy(p, subSlice)
	if n < len(p) {
		err = io.ErrShortBuffer
	}
	return
}

func NewReadBuffer(data []byte) ReadBuffer {
	return ReadBuffer{
		data: data,
	}
}

func (r *ReadBuffer) ReadRune() (rn rune, size int, err error) {
	rn, size = utf8.DecodeRune(r.data[r.pos:])
	if rn == utf8.RuneError {
		if size == 0 {
			err = io.ErrShortBuffer
		} else {
			err = ErrInvalidUTF8
		}
	}
	r.pos += size
	return rn, size, err
}

func (r *ReadBuffer) ReadByte() (byte, error) {
	if len(r.data)-r.pos < 1 {
		return 0, io.ErrShortBuffer
	}
	b := r.data[r.pos]
	r.pos += 1
	return b, nil
}

func (r *ReadBuffer) CheckReadSpace(n int) (nn int, err error) {
	nn = min(n, len(r.data)-r.pos)
	if nn != n {
		err = io.ErrShortBuffer
	}
	return
}

func (r *ReadBuffer) Read(p []byte) (n int, err error) {
	n = copy(p[:], r.data[r.pos:])
	if n <= 0 {
		return 0, io.EOF
	}
	r.pos += n
	return n, nil
}

func (r *ReadBuffer) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := w.Write(r.data[r.pos:])
	r.pos += nn
	return int64(nn), err
}

func (r *ReadBuffer) Len() int {
	return len(r.data)
}

func (r *ReadBuffer) UnreadBytesRef() []byte {
	return r.data[r.pos:]
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
var _ io.ByteReader = (*ReadBuffer)(nil)
var _ io.RuneReader = (*ReadBuffer)(nil)
var _ io.ReaderAt = (*ReadBuffer)(nil)
var _ io.Closer = (*ReadBuffer)(nil)
