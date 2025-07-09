package errorresult

import (
	"encoding/binary"
	"io"

	msg "github.com/gabe-lee/OurSweeper/internal/messages"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/gabe-lee/OurSweeper/internal/wire_serializer"
)

type ErrorResult uint32

func (e *ErrorResult) ReadWire(r io.Reader, order binary.ByteOrder, code uint32) error {
	ec := utils.ErrorChecker{}
	switch code {
	case msg.ERROR:
		if ec.IsErr(binary.Read(r, order, (*uint32)(e))) {
			return ec.Err
		}
	default:
		return wire_serializer.MakeError(code, "ErrorResult", []uint32{msg.ERROR})
	}
	return nil
}

func (e *ErrorResult) WriteWire(w io.Writer, order binary.ByteOrder, code uint32) error {
	ec := utils.ErrorChecker{}
	switch code {
	case msg.ERROR:
		if ec.IsErr(binary.Write(w, order, code)) {
			return ec.Err
		}
		if ec.IsErr(binary.Write(w, order, (*uint32)(e))) {
			return ec.Err
		}
	default:
		return wire_serializer.MakeError(code, "ErrorResult", []uint32{msg.ERROR})
	}
	return nil
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
