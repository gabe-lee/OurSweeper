package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"

	"github.com/gabe-lee/OurSweeper/data_buffer"
)

const HMAC_LEN = sha256.Size
const MIN_TOKEN_LEN = HMAC_LEN + 1

var HASHER = sha256.New

type (
	WriteBuffer     = data_buffer.WriteBuffer
	WriteBufferPool = data_buffer.WriteBufferPool
	ReadBuffer      = data_buffer.ReadBuffer
	Readable        = data_buffer.Readable
	Writable        = data_buffer.Writable
	Sizable         = data_buffer.Sizable
	SizeWritable    = data_buffer.SizeWritable
	Hash            = hash.Hash
)

func Create(secret []byte, inputPayload Writable, outputBuffer *WriteBuffer) []byte {
	start := outputBuffer.Len()
	outputBuffer.Writable(inputPayload)
	end := outputBuffer.Len()
	written := outputBuffer.BytesRef()[start:end]
	HMAC := hmac.New(HASHER, secret)
	HMAC.Write(written)
	outputBuffer.Hash(HMAC)
	end = outputBuffer.Len()
	return outputBuffer.BytesRef()[start:end]
}

func OpenAndValidate(secret []byte, inputBuffer *ReadBuffer, outputPayload Readable) (valid bool, err error) {
	start := inputBuffer.Pos()
	err = inputBuffer.Readable(outputPayload)
	if err != nil {
		return false, err
	}
	end := inputBuffer.Pos()
	sentHMAC, err := inputBuffer.Discard(HMAC_LEN)
	if err != nil {
		return false, err
	}
	HMAC := hmac.New(HASHER, secret)
	_, err = HMAC.Write(inputBuffer.BytesRef()[start:end])
	if err != nil {
		return false, err
	}
	expectedHMAC := make([]byte, 0, HMAC_LEN)
	expectedHMAC = HMAC.Sum(expectedHMAC)
	valid = hmac.Equal(expectedHMAC, sentHMAC)
	return valid, nil
}

func Open(inputBuffer *ReadBuffer, outputPayload Readable) (tokenData, sentHMAC []byte, err error) {
	start := inputBuffer.Pos()
	err = inputBuffer.Readable(outputPayload)
	if err != nil {
		return nil, nil, err
	}
	end := inputBuffer.Pos()
	tokenData = inputBuffer.BytesRef()[start:end]
	sentHMAC, err = inputBuffer.Discard(HMAC_LEN)
	return
}
