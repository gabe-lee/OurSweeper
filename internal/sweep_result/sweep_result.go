package sweep_result

import (
	"encoding/binary"
	"io"

	"github.com/gabe-lee/OurSweeper/internal/coord"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/gabe-lee/OurSweeper/internal/wire_serializer"
)

const (
	OWN_SWEEP uint32 = iota
	OTHER_SWEEP
)

type (
	Coord     = coord.Coord
	ByteCoord = coord.ByteCoord
)

type SweepResult struct {
	Score  uint16
	Coords [36]ByteCoord
	Icons  [36]byte
	Len    byte
}

func (s *SweepResult) AddTile(score uint16, icon byte, pos Coord) {
	s.Score += score
	s.Coords[s.Len] = pos.ToByteCoord()
	s.Icons[s.Len] = icon
	s.Len += 1
}

func (s *SweepResult) WriteWire(w io.Writer, order binary.ByteOrder, code uint32) error {
	var e utils.ErrorChecker
	switch code {
	case OWN_SWEEP:
		if e.IsErr(binary.Write(w, order, OWN_SWEEP)) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Score)) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Len)) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Coords[:(s.Len<<1)])) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Icons[:s.Len])) {
			return e.Err
		}
	case OTHER_SWEEP:
		if e.IsErr(binary.Write(w, order, OTHER_SWEEP)) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Len)) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Coords[:(s.Len<<1)])) {
			return e.Err
		}
		if e.IsErr(binary.Write(w, order, s.Icons[:s.Len])) {
			return e.Err
		}
	default:
		return wire_serializer.MakeError(code, "SweepResult", []uint32{OWN_SWEEP, OTHER_SWEEP})
	}
	return nil
}
func (s *SweepResult) ReadWire(r io.Reader, order binary.ByteOrder) error {
	var e utils.ErrorChecker
	var code uint32
	if e.IsErr(binary.Read(r, order, &code)) {
		return e.Err
	}
	switch code {
	case OWN_SWEEP:
		if e.IsErr(binary.Read(r, order, s.Score)) {
			return e.Err
		}
		if e.IsErr(binary.Read(r, order, s.Len)) {
			return e.Err
		}
		if e.IsErr(binary.Read(r, order, s.Coords[:(s.Len<<1)])) {
			return e.Err
		}
		if e.IsErr(binary.Read(r, order, s.Icons[:s.Len])) {
			return e.Err
		}
	case OTHER_SWEEP:
		if e.IsErr(binary.Read(r, order, s.Len)) {
			return e.Err
		}
		if e.IsErr(binary.Read(r, order, s.Coords[:(s.Len<<1)])) {
			return e.Err
		}
		if e.IsErr(binary.Read(r, order, s.Icons[:s.Len])) {
			return e.Err
		}
	default:
		return wire_serializer.MakeError(code, "SweepResult", []uint32{OWN_SWEEP, OTHER_SWEEP})
	}
	return nil
}

var _ wire_serializer.WireSerializer = (*SweepResult)(nil)
