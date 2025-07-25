package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"

	"github.com/gabe-lee/OurSweeper/data_buffer"
	"github.com/gabe-lee/OurSweeper/wire"
)

const HMAC_LEN = sha256.Size
const MIN_TOKEN_LEN = HMAC_LEN + 1

var HASHER = sha256.New
var BYTE_ORDER = wire.LE

type (
	WireWriter      = wire.WireWriter
	WireSizeWriter  = wire.WireSizeWriter
	WireReader      = wire.WireReader
	WireSizeReader  = wire.WireSizeReader
	IncomingWire    = wire.Incoming
	OutgoingWire    = wire.Outgoing
	WriteBuffer     = data_buffer.WriteBuffer
	WriteBufferPool = data_buffer.WriteBufferPool
	ReadBuffer      = data_buffer.ReadBuffer
	Hash            = hash.Hash
)

// type tokenWriter struct {
// 	buf  WriteBuffer
// 	err  error
// 	hash Hash
// }

// func newWriter(secret []byte, payloadSize int) tokenWriter {
// 	t := tokenWriter{
// 		buf:  data_buffer.NewWriteBuffer(payloadSize + HMAC_LEN),
// 		hash: hmac.New(HASHER, secret),
// 	}
// 	return t
// }

// func (t *tokenWriter) Write(p []byte) (n int, err error) {
// 	n, _ = t.buf.Write(p)
// 	_, err = t.hash.Write(p)
// 	if t.err == nil {
// 		t.err = err
// 	}
// 	return
// }

// func (t *tokenWriter) WriteNoHash(p []byte) {
// 	t.buf.Write(p)
// }

// type tokenValidReader struct {
// 	buf  ReadBuffer
// 	err  error
// 	hash Hash
// }

// func newValidReader(secret []byte, data []byte) tokenValidReader {
// 	t := tokenValidReader{
// 		buf:  data_buffer.NewReadBuffer(data),
// 		hash: hmac.New(HASHER, secret),
// 	}
// 	return t
// }

// func (t *tokenValidReader) Read(p []byte) (n int, err error) {
// 	n, err = t.buf.Read(p)
// 	if t.err == nil {
// 		t.err = err
// 	}
// 	_, err = t.hash.Write(p)
// 	if t.err == nil {
// 		t.err = err
// 	}
// 	return
// }

// var _ io.Reader = (*tokenValidReader)(nil)

// type tokenReader struct {
// 	buf ReadBuffer
// 	err error
// }

// func newReader(data []byte) tokenReader {
// 	t := tokenReader{
// 		buf: data_buffer.NewReadBuffer(data),
// 	}
// 	return t
// }

// func (t *tokenReader) Read(p []byte) (n int, err error) {
// 	n, err = t.buf.Read(p)
// 	if t.err == nil {
// 		t.err = err
// 	}
// 	return
// }

// var _ io.Reader = (*tokenReader)(nil)

func Create(secret []byte, inputPayload WireSizeWriter, outputBuffer *WriteBuffer) error {
	size := inputPayload.WireSize()
	fullSize := size + HMAC_LEN
	write := wire.NewOutgoingAdv(outputBuffer, BYTE_ORDER, hmac.New(HASHER, secret), nil, fullSize)
	write.Struct(inputPayload)
	write.OwnHash()
	return write.Err
}

func OpenAndValidate(secret []byte, inputBuffer *ReadBuffer, outputPayload WireReader) (valid bool, err error) {
	read := wire.NewIncomingAdv(inputBuffer, BYTE_ORDER, hmac.New(HASHER, secret), nil, inputBuffer.Len())
	read.Struct(outputPayload)
	sentHMAC := inputBuffer.BytesRef()[:HMAC_LEN]
	expectedHMAC := make([]byte, 0, HMAC_LEN)
	expectedHMAC = read.Hasher.Sum(expectedHMAC)
	valid = hmac.Equal(expectedHMAC, sentHMAC)
	return valid, read.Err
}

func Open(inputBuffer *ReadBuffer, outputPayload WireReader) error {
	read := wire.NewIncomingAdv(inputBuffer, BYTE_ORDER, nil, nil, inputBuffer.Len())
	read.Struct(outputPayload)
	return read.Err
}
