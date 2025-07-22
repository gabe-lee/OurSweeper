package common

import (
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/wire"
)

type (
	Writer    = io.Writer
	Reader    = io.Reader
	ByteOrder = binary.ByteOrder
)

type SweepRequest struct {
	Pos ByteCoord
}

func (s *SweepRequest) WireRead(w *wire.Incoming) {
	s.Pos.WireRead(w)
}

func (s *SweepRequest) WireWrite(w *wire.Outgoing) {
	s.Pos.WireWrite(w)
}

func NewSweepRequest(pos Coord) SweepRequest {
	return SweepRequest{
		Pos: pos.ToCoordByte(),
	}
}

var _ wire.WireReader = (*SweepRequest)(nil)
var _ wire.WireWriter = (*SweepRequest)(nil)
