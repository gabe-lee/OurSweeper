package wire

import (
	"errors"
)

type WireReader interface {
	// Read the serialized type data from the provided IncomingWire
	//
	// No errors are returned/cached by the Incoming wire. If errors are needed,
	// they must be cached by the underlying io.Reader
	WireRead(read *Incoming)
}

type WireWriter interface {
	// Write the serialized type data into the provided OutgoingWire
	//
	// No errors are returned/cached by the Outgoing wire. If errors are needed,
	// they must be cached by the underlying io.Writer
	WireWrite(write *Outgoing)
}

type WireSizer interface {
	// Report the calculated total length of a Write/Read of this type
	WireSize() int
}

type WireSizeReader interface {
	WireReader
	WireSizer
}

type WireSizeWriter interface {
	WireWriter
	WireSizer
}

type WireReadWriter interface {
	WireWriter
	WireReader
}

type WireSizeReadWriter interface {
	WireWriter
	WireReader
	WireSizer
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
	Read16(src *[2]byte, dst *uint16)
	Read24(src *[3]byte, dst *uint32)
	Read32(src *[4]byte, dst *uint32)
	Read40(src *[5]byte, dst *uint64)
	Read48(src *[6]byte, dst *uint64)
	Read56(src *[7]byte, dst *uint64)
	Read64(src *[8]byte, dst *uint64)
	Write16(src uint16, dst *[2]byte)
	Write24(src uint32, dst *[3]byte)
	Write32(src uint32, dst *[4]byte)
	Write40(src uint64, dst *[5]byte)
	Write48(src uint64, dst *[6]byte)
	Write56(src uint64, dst *[7]byte)
	Write64(src uint64, dst *[8]byte)
}

var LE le

type le struct{}

func (bin le) Read16(src *[2]byte, dst *uint16) {
	*dst = uint16(src[0])
	*dst |= uint16(src[1]) << 8
}

func (le) Read24(src *[3]byte, dst *uint32) {
	*dst = uint32(src[0])
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2]) << 16
}

func (le) Read32(src *[4]byte, dst *uint32) {
	*dst = uint32(src[0])
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2]) << 16
	*dst |= uint32(src[3]) << 24
}

func (le) Read40(src *[5]byte, dst *uint64) {
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
}

func (le) Read48(src *[6]byte, dst *uint64) {
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
}

func (le) Read56(src *[7]byte, dst *uint64) {
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
	*dst |= uint64(src[6]) << 48
}

func (le) Read64(src *[8]byte, dst *uint64) {
	*dst = uint64(src[0])
	*dst |= uint64(src[1]) << 8
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 32
	*dst |= uint64(src[5]) << 40
	*dst |= uint64(src[6]) << 48
	*dst |= uint64(src[7]) << 56
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

func (le) Write56(src uint64, dst *[7]byte) {
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

func (be) Read16(src *[2]byte, dst *uint16) {
	*dst = uint16(src[0]) << 8
	*dst |= uint16(src[1])
}

func (be) Read24(src *[3]byte, dst *uint32) {
	*dst = uint32(src[0]) << 16
	*dst |= uint32(src[1]) << 8
	*dst |= uint32(src[2])
}

func (be) Read32(src *[4]byte, dst *uint32) {
	*dst = uint32(src[0]) << 24
	*dst |= uint32(src[1]) << 16
	*dst |= uint32(src[2]) << 8
	*dst |= uint32(src[3])
}

func (be) Read40(src *[5]byte, dst *uint64) {
	*dst = uint64(src[0]) << 32
	*dst |= uint64(src[1]) << 24
	*dst |= uint64(src[2]) << 16
	*dst |= uint64(src[3]) << 8
	*dst |= uint64(src[4])
}

func (be) Read48(src *[6]byte, dst *uint64) {
	*dst = uint64(src[0]) << 40
	*dst |= uint64(src[1]) << 32
	*dst |= uint64(src[2]) << 24
	*dst |= uint64(src[3]) << 16
	*dst |= uint64(src[4]) << 8
	*dst |= uint64(src[5])
}

func (be) Read56(src *[7]byte, dst *uint64) {
	*dst = uint64(src[0]) << 48
	*dst |= uint64(src[1]) << 40
	*dst |= uint64(src[2]) << 32
	*dst |= uint64(src[3]) << 24
	*dst |= uint64(src[4]) << 16
	*dst |= uint64(src[5]) << 8
	*dst |= uint64(src[6])
}

func (be) Read64(src *[8]byte, dst *uint64) {
	*dst = uint64(src[0]) << 56
	*dst |= uint64(src[1]) << 48
	*dst |= uint64(src[2]) << 40
	*dst |= uint64(src[3]) << 32
	*dst |= uint64(src[4]) << 24
	*dst |= uint64(src[5]) << 16
	*dst |= uint64(src[6]) << 8
	*dst |= uint64(src[7])
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

func (be) Write56(src uint64, dst *[7]byte) {
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
