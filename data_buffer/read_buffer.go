package data_buffer

import "io"

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

var _ io.Reader = (*ReadBuffer)(nil)
var _ io.WriterTo = (*ReadBuffer)(nil)
