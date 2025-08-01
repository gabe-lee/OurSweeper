package common

import (
	"encoding/binary"
	"io"
)

type (
	Writer    = io.Writer
	Reader    = io.Reader
	ByteOrder = binary.ByteOrder
)

type SweepRequest struct {
	Pos ByteCoord
}

func NewSweepRequest(pos Coord) SweepRequest {
	return SweepRequest{
		Pos: pos.ToCoordByte(),
	}
}
