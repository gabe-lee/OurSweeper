package wire

import "errors"

type WireReader interface {
	// Read the serialized type data from the provided IncomingWire
	//
	// Any Errors should be attatched to the `IncomingWire`
	WireRead(wire *IncomingWire)
}

type WireWriter interface {
	// Write the serialized type data into the provided OutgoingWire
	//
	// Any Errors should be attatched to the `IncomingWire`
	WireWrite(wire *OutgoingWire)
}

var ErrVarintOverflow64 = errors.New("reading varint causes overflow of target 64-bit integer")

const maxVarInt64Len = 9
const maxVarInt64LastIdx = maxVarInt64Len - 1

var ErrVarintOverflow32 = errors.New("reading varint causes overflow of target 32-bit integer")

const maxVarInt32Len = 5
const maxVarInt32LastIdx = maxVarInt32Len - 1
const maxVarint32LastByte = 0b00001111

var ErrVarintOverflow16 = errors.New("reading varint causes overflow of target 16-bit integer")

const maxVarInt16Len = 3
const maxVarInt16LastIdx = maxVarInt16Len - 1
const maxVarint16LastByte = 0b00000011

type Order interface {
	ReadU16(src [2]byte, dst *uint16)
	ReadU32(src [4]byte, dst *uint32)
	ReadU64(src [8]byte, dst *uint64)
	WriteU16(src uint16, dst *[2]byte)
	WriteU32(src uint32, dst *[4]byte)
	WriteU64(src uint64, dst *[8]byte)
}

var LE le

type le struct{}

func (le) ReadU16(src [2]byte, dst *uint16) {
	*dst = uint16(src[0]) | uint16(src[1])<<8
}

func (le) ReadU32(src [4]byte, dst *uint32) {
	*dst = uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16 | uint32(src[3])<<24
}

func (le) ReadU64(src [8]byte, dst *uint64) {
	*dst = uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32 | uint64(src[5])<<40 | uint64(src[6])<<48 | uint64(src[7])<<56
}

func (le) WriteU16(src uint16, dst *[2]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
}

func (le) WriteU32(src uint32, dst *[4]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
}

func (le) WriteU64(src uint64, dst *[8]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
	dst[4] = byte(src >> 32)
	dst[5] = byte(src >> 40)
	dst[6] = byte(src >> 48)
	dst[7] = byte(src >> 56)
}

var _ Order = le{}

var BE be

type be struct{}

func (be) ReadU16(src [2]byte, dst *uint16) {
	*dst = uint16(src[1]) | uint16(src[0])<<8
}

func (be) ReadU32(src [4]byte, dst *uint32) {
	*dst = uint32(src[3]) | uint32(src[2])<<8 | uint32(src[1])<<16 | uint32(src[0])<<24
}

func (be) ReadU64(src [8]byte, dst *uint64) {
	*dst = uint64(src[7]) | uint64(src[6])<<8 | uint64(src[5])<<16 | uint64(src[4])<<24 |
		uint64(src[3])<<32 | uint64(src[2])<<40 | uint64(src[1])<<48 | uint64(src[0])<<56
}

func (be) WriteU16(src uint16, dst *[2]byte) {
	dst[0] = byte(src >> 8)
	dst[1] = byte(src)
}

func (be) WriteU32(src uint32, dst *[4]byte) {
	dst[0] = byte(src >> 24)
	dst[1] = byte(src >> 16)
	dst[2] = byte(src >> 8)
	dst[3] = byte(src)
}

func (be) WriteU64(src uint64, dst *[8]byte) {
	dst[0] = byte(src >> 56)
	dst[1] = byte(src >> 48)
	dst[2] = byte(src >> 40)
	dst[3] = byte(src >> 32)
	dst[4] = byte(src >> 24)
	dst[5] = byte(src >> 16)
	dst[6] = byte(src >> 8)
	dst[7] = byte(src)
}

var _ Order = be{}
