package sweep_request

import (
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/internal/coord"
	"github.com/gabe-lee/OurSweeper/internal/wire_serializer"
)

type (
	Coord     = coord.Coord
	ByteCoord = coord.ByteCoord
	ByteOrder = binary.ByteOrder
	Writer    = io.Writer
	Reader    = io.Reader
)

type SweepRequest struct {
	Pos ByteCoord
}

// ReadWire implements wire_serializer.WireSerializer.
func (s *SweepRequest) ReadWire(r Reader, order ByteOrder) error {
	return s.Pos.ReadWire(r, order)
}

// WriteWire implements wire_serializer.WireSerializer.
func (s *SweepRequest) WriteWire(w Writer, order ByteOrder, code uint32) error {
	return s.Pos.WriteWire(w, order, 0)
}

func NewSweepRequest(pos coord.Coord) SweepRequest {
	return SweepRequest{
		Pos: pos.ToByteCoord(),
	}
}

var _ wire_serializer.WireSerializer = (*SweepRequest)(nil)
