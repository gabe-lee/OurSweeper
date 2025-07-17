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
	IncomingWire = wire.IncomingWire
	OutgoingWire = wire.OutgoingWire
)

func Create(secret []byte, inputPayload WireWriter) (token []byte, err error) {
	buf := bytes.Buffer{}
	w := wire.NewOutgoing(&buf, wire.LE)
	w.TryWrite_WireWriter(inputPayload)
	if w.HasErr() {
		return nil, w.Err()
	}
	if w.Len() <= 0 {
		return nil, fmt.Errorf("token payload must be at least 1 byte of data")
	}
	hasher := hmac.New(HASHER, secret)
	_, err = hasher.Write(buf.Bytes())
	if err != nil {
		return nil, err
	}
	HMAC := make([]byte, 0, HMAC_LEN)
	HMAC = hasher.Sum(HMAC)
	w.TryWrite_SliceU8(HMAC)
	return buf.Bytes(), w.Err()
}

func OpenAndValidate(secret []byte, token []byte, outputPayload WireReader) (valid bool, err error) {
	tlen := len(token)
	if tlen < MIN_TOKEN_LEN {
		return false, fmt.Errorf("invalid token: must be at least %d bytes long (%d bytes for HMAC and at least 1 byte of data), got %d bytes", MIN_TOKEN_LEN, HMAC_LEN, tlen)
	}
	hmacStart := tlen - HMAC_LEN
	sentHMAC := token[hmacStart:]
	sentPayload := token[:hmacStart]
	hasher := hmac.New(HASHER, secret)
	_, err = hasher.Write([]byte(sentPayload))
	if err != nil {
		return false, err
	}
	expectedHMAC := make([]byte, 0, HMAC_LEN)
	expectedHMAC = hasher.Sum(expectedHMAC)
	valid = hmac.Equal(expectedHMAC, sentHMAC)
	w := wire.NewIncomingSlice(token, wire.LE)
	w.TryRead_WireReader(outputPayload)
	return valid, w.Err()
}

func Open(token []byte, outputPayload WireReader) error {
	w := wire.NewIncomingSlice(token, wire.LE)
	w.TryRead_WireReader(outputPayload)
	return w.Err()
}
