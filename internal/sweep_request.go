package internal

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/serializer"
	"github.com/gabe-lee/OurSweeper/utils"
)

type (
	Writer    = io.Writer
	Reader    = io.Reader
	ByteOrder = binary.ByteOrder
)

type SweepRequest struct {
	Pos ByteCoord
}

func (s *SweepRequest) Deserialize(r Reader, order ByteOrder) error {
	return s.Pos.Deserialize(r, order)
}

func (s *SweepRequest) CodedSerialize(order ByteOrder) (data []byte, err error) {
	e := utils.ErrorCollector{}
	data = make([]byte, 0, 6)
	w := bytes.NewBuffer(data)
	e.Do(binary.Write(w, order, CLIENT_SWEEP))
	e.Do(s.Pos.Serialize(w, order))
	return w.Bytes(), e.Err
}

func NewSweepRequest(pos Coord) SweepRequest {
	return SweepRequest{
		Pos: pos.ToCoordByte(),
	}
}

var _ serializer.CodedSerializer = (*SweepRequest)(nil)
var _ serializer.Deserializer = (*SweepRequest)(nil)
