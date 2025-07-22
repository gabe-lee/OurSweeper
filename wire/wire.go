package wire

import (
	"errors"
	"hash"
	"io"
	"sync"
)

type (
	Hash            = hash.Hash
	Locker          = sync.Locker
	ReadWriteSeeker = io.ReadWriteSeeker
)

type WireReader interface {
	// Read the serialized type data from the provided IncomingWire
	//
	// Any Errors should be attatched to the `IncomingWire`
	WireRead(read *Incoming)
}

type WireWriter interface {
	// Write the serialized type data into the provided OutgoingWire
	//
	// Any Errors should be attatched to the `IncomingWire`
	WireWrite(write *Outgoing)
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
	Read16(src [2]byte, dst *uint16)
	Read24(src [3]byte, dst *uint32)
	Read32(src [4]byte, dst *uint32)
	Read40(src [5]byte, dst *uint64)
	Read48(src [6]byte, dst *uint64)
	Read54(src [7]byte, dst *uint64)
	Read64(src [8]byte, dst *uint64)
	Write16(src uint16, dst *[2]byte)
	Write24(src uint32, dst *[3]byte)
	Write32(src uint32, dst *[4]byte)
	Write40(src uint64, dst *[5]byte)
	Write48(src uint64, dst *[6]byte)
	Write54(src uint64, dst *[7]byte)
	Write64(src uint64, dst *[8]byte)
}

var LE le

type le struct{}

func (le) Read16(src [2]byte, dst *uint16) {
	*dst = uint16(src[0]) | uint16(src[1])<<8
}

func (le) Read24(src [3]byte, dst *uint32) {
	*dst = uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16
}

func (le) Read32(src [4]byte, dst *uint32) {
	*dst = uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16 | uint32(src[3])<<24
}

func (le) Read40(src [5]byte, dst *uint64) {
	*dst = uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32
}

func (le) Read48(src [6]byte, dst *uint64) {
	*dst = uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32 | uint64(src[5])<<40
}

func (le) Read54(src [7]byte, dst *uint64) {
	*dst = uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32 | uint64(src[5])<<40 | uint64(src[6])<<48
}

func (le) Read64(src [8]byte, dst *uint64) {
	*dst = uint64(src[0]) | uint64(src[1])<<8 | uint64(src[2])<<16 | uint64(src[3])<<24 |
		uint64(src[4])<<32 | uint64(src[5])<<40 | uint64(src[6])<<48 | uint64(src[7])<<56
}

func (le) Write16(src uint16, dst *[2]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
}

func (le) Write24(src uint32, dst *[3]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
}

func (le) Write32(src uint32, dst *[4]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
}

func (le) Write40(src uint64, dst *[5]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
	dst[4] = byte(src >> 32)
}
func (le) Write48(src uint64, dst *[6]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
	dst[4] = byte(src >> 32)
	dst[5] = byte(src >> 40)
}
func (le) Write54(src uint64, dst *[7]byte) {
	dst[0] = byte(src)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 24)
	dst[4] = byte(src >> 32)
	dst[5] = byte(src >> 40)
	dst[6] = byte(src >> 48)
}

func (le) Write64(src uint64, dst *[8]byte) {
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

func (be) Read16(src [2]byte, dst *uint16) {
	*dst = uint16(src[1]) | uint16(src[0])<<8
}

func (be) Read24(src [3]byte, dst *uint32) {
	*dst = uint32(src[2]) | uint32(src[1])<<8 | uint32(src[0])<<16
}

func (be) Read32(src [4]byte, dst *uint32) {
	*dst = uint32(src[3]) | uint32(src[2])<<8 | uint32(src[1])<<16 | uint32(src[0])<<24
}

func (be) Read40(src [5]byte, dst *uint64) {
	*dst = uint64(src[4]) | uint64(src[3])<<8 | uint64(src[2])<<16 | uint64(src[1])<<24 |
		uint64(src[0])<<32
}

func (be) Read48(src [6]byte, dst *uint64) {
	*dst = uint64(src[5]) | uint64(src[4])<<8 | uint64(src[3])<<16 | uint64(src[2])<<24 |
		uint64(src[1])<<32 | uint64(src[0])<<40
}

func (be) Read54(src [7]byte, dst *uint64) {
	*dst = uint64(src[6]) | uint64(src[5])<<8 | uint64(src[4])<<16 | uint64(src[3])<<24 |
		uint64(src[2])<<32 | uint64(src[1])<<40 | uint64(src[0])<<48
}

func (be) Read64(src [8]byte, dst *uint64) {
	*dst = uint64(src[7]) | uint64(src[6])<<8 | uint64(src[5])<<16 | uint64(src[4])<<24 |
		uint64(src[3])<<32 | uint64(src[2])<<40 | uint64(src[1])<<48 | uint64(src[0])<<56
}

func (be) Write16(src uint16, dst *[2]byte) {
	dst[0] = byte(src >> 8)
	dst[1] = byte(src)
}

func (be) Write24(src uint32, dst *[3]byte) {
	dst[0] = byte(src >> 16)
	dst[1] = byte(src >> 8)
	dst[2] = byte(src)
}

func (be) Write32(src uint32, dst *[4]byte) {
	dst[0] = byte(src >> 24)
	dst[1] = byte(src >> 16)
	dst[2] = byte(src >> 8)
	dst[3] = byte(src)
}

func (be) Write40(src uint64, dst *[5]byte) {
	dst[0] = byte(src >> 32)
	dst[1] = byte(src >> 24)
	dst[2] = byte(src >> 16)
	dst[3] = byte(src >> 8)
	dst[4] = byte(src)
}

func (be) Write48(src uint64, dst *[6]byte) {
	dst[0] = byte(src >> 40)
	dst[1] = byte(src >> 32)
	dst[2] = byte(src >> 24)
	dst[3] = byte(src >> 16)
	dst[4] = byte(src >> 8)
	dst[5] = byte(src)
}

func (be) Write54(src uint64, dst *[7]byte) {
	dst[0] = byte(src >> 48)
	dst[1] = byte(src >> 40)
	dst[2] = byte(src >> 32)
	dst[3] = byte(src >> 24)
	dst[4] = byte(src >> 16)
	dst[5] = byte(src >> 8)
	dst[6] = byte(src)
}

func (be) Write64(src uint64, dst *[8]byte) {
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
