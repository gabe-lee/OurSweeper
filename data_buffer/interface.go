package data_buffer

import (
	"errors"
	"io"
)

type Writable interface {
	WriteToBuffer(buf *WriteBuffer)
}

type Readable interface {
	ReadFromBuffer(buf *ReadBuffer) error
}

type Sizable interface {
	SizeOnBuffer() int
}

type SizeReadable interface {
	Readable
	Sizable
}

type SizeWritable interface {
	Writable
	Sizable
}

type ReadWritable interface {
	Writable
	Readable
}

type SizeReadWritable interface {
	Writable
	Readable
	Sizable
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

const (
	smallestGrow = 64

	maxVarInt64Len      = 9
	maxVarInt64LastIdx  = maxVarInt64Len - 1
	maxVarInt32Len      = 5
	maxVarInt32LastIdx  = maxVarInt32Len - 1
	maxVarint32LastByte = 0b00001111
	maxVarInt16Len      = 3
	maxVarInt16LastIdx  = maxVarInt16Len - 1
	maxVarint16LastByte = 0b00000011
)

const (
	SeekStart   = io.SeekStart
	SeekCurrent = io.SeekCurrent
	SeekEnd     = io.SeekEnd
)

var (
	ErrInvalidWhence    = errors.New("invalid whence parameter: only 0 (SeekStart), 1 (SeekCurrent), or 2 (SeekEnd) are allowed")
	ErrInvalidSeek      = errors.New("seek operation would put write position outside buffer range")
	ErrInvalidUTF8      = errors.New("invalid utf8 encoding")
	ErrVarintOverflow16 = errors.New("reading varint causes overflow of target 16-bit integer")
	ErrVarintOverflow32 = errors.New("reading varint causes overflow of target 32-bit integer")
	ErrVarintOverflow64 = errors.New("reading varint causes overflow of target 64-bit integer")
)

type (
	Writer = io.Writer
)
