package sweep_result

import (
	"encoding/binary"
	"io"

	msg "github.com/gabe-lee/OurSweeper/internal/messages"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/gabe-lee/OurSweeper/internal/wire_serializer"
)

type SweepResult struct {
	Score  uint32
	Coords [74]byte
	Icons  [36]byte
	Len    byte
}

func (s *SweepResult) AddTile(score uint32, icon byte, x, y int) {
	s.Score += score
	c := s.Len * 2
	s.Coords[c] = byte(x)
	s.Coords[c+1] = byte(y)
	s.Icons[s.Len] = icon
	s.Len += 1
}

func (s *SweepResult) WriteWire(w io.Writer, order binary.ByteOrder, code uint32) error {
	var e utils.ErrorChecker
	switch code {
	case msg.SERVER_SWEEP_OWN:
		if e.IsErr(binary.Write(w, order, msg.SERVER_SWEEP_OWN)) {
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
	case msg.SERVER_SWEEP_OTHER:
		if e.IsErr(binary.Write(w, order, msg.SERVER_SWEEP_OTHER)) {
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
		return wire_serializer.MakeError(code, "SweepResult", []uint32{msg.SERVER_SWEEP_OWN, msg.SERVER_SWEEP_OTHER})
	}
	return nil
}
func (s *SweepResult) ReadWire(r io.Reader, order binary.ByteOrder, code uint32) error {
	var e utils.ErrorChecker
	switch code {
	case msg.SERVER_SWEEP_OWN:
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
	case msg.SERVER_SWEEP_OTHER:
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
		return wire_serializer.MakeError(code, "SweepResult", []uint32{msg.SERVER_SWEEP_OWN, msg.SERVER_SWEEP_OTHER})
	}
	return nil
}

var _ wire_serializer.WireSerializer = (*SweepResult)(nil)
