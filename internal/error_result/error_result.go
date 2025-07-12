package errorresult

import (
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/internal/wire_serializer"
)

type ErrorResult uint32

func (e *ErrorResult) ReadWire(r io.Reader, order binary.ByteOrder) error {
	return binary.Read(r, order, (*uint32)(e))
}

func (e *ErrorResult) WriteWire(w io.Writer, order binary.ByteOrder, code uint32) error {
	return binary.Write(w, order, (*uint32)(e))
}

const (
	INVALID_USERNAME_PASS uint32 = iota
	TILE_ALREADY_SWEPT
	msgCount
)

var MSG = [msgCount]string{
	INVALID_USERNAME_PASS: "Invalid username or password",
	TILE_ALREADY_SWEPT:    "The tile is already swept",
}

var _ wire_serializer.WireSerializer = (*ErrorResult)(nil)
