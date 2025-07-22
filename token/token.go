package token

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/gabe-lee/OurSweeper/wire"
)

const HMAC_LEN = sha256.Size
const MIN_TOKEN_LEN = HMAC_LEN + 1

var HASHER = sha256.New
var BYTE_ORDER = binary.LittleEndian

type (
	WireWriter   = wire.WireWriter
	WireReader   = wire.WireReader
	IncomingWire = wire.Incoming
	OutgoingWire = wire.Outgoing
)

func Create(secret []byte, inputPayload WireWriter) (token []byte, err error) {
	buf := bytes.Buffer{}
	hasher := hmac.New(HASHER, secret)
	write := wire.NewOutgoing(&buf, wire.LE)
	write.OnWrite = func(writtenBytes []byte) {
		_, e := hasher.Write(writtenBytes)
		write.ReplaceErr(e)
	}
	write.Struct(inputPayload)
	if write.HasErr() {
		return nil, write.Err()
	}
	if buf.Len() <= 0 {
		return nil, fmt.Errorf("token payload must be at least 1 byte of data")
	}
	_, err = hasher.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	HMAC := make([]byte, 0, HMAC_LEN)
	HMAC = hasher.Sum(HMAC)
	write.OnWrite = nil
	write.U8_Slice(HMAC)
	return buf.Bytes(), write.Err()
}

func OpenAndValidate(secret []byte, token []byte, outputPayload WireReader) (valid bool, err error) {
	tlen := len(token)
	if tlen < MIN_TOKEN_LEN {
		return false, fmt.Errorf("invalid token: must be at least %d bytes long (%d bytes for HMAC and at least 1 byte of data), got %d bytes", MIN_TOKEN_LEN, HMAC_LEN, tlen)
	}
	hmacStart := tlen - HMAC_LEN
	sentHMAC := token[hmacStart:]
	hasher := hmac.New(HASHER, secret)
	read := wire.NewIncomingSlice(token[0:hmacStart], wire.LE)
	read.OnRead = func(readBytes []byte) {
		_, e := hasher.Write(readBytes)
		read.ReplaceErr(e)
	}
	read.Struct(outputPayload)
	expectedHMAC := make([]byte, 0, HMAC_LEN)
	expectedHMAC = hasher.Sum(expectedHMAC)
	valid = hmac.Equal(expectedHMAC, sentHMAC)
	return valid, read.Err()
}

func Open(token []byte, outputPayload WireReader) error {
	read := wire.NewIncomingSlice(token, wire.LE)
	read.Struct(outputPayload)
	return read.Err()
}
