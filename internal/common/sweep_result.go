package common

import (
	"math/bits"

	"github.com/gabe-lee/OurSweeper/wire"
)

const (
	OWN_SWEEP uint32 = iota
	OTHER_SWEEP
)

type SweepResult struct {
	Score        uint16
	Center       ByteCoord
	RelativeBits uint64
	Icons        [MAX_ICON_LEN]byte
	Len          byte
}

func (s *SweepResult) InitSweep(pos Coord, score uint16, icon byte) {
	s.Score = score
	s.Center = pos.ToCoordByte()
	s.Icons[0] = icon
	s.Len = 1
}

func (s *SweepResult) AddCascadeSweep(icon byte, bit uint64) {
	s.Score += uint16(BOMB_NEAR_BASE_SCORE[0])
	i := s.Len >> 1
	o := (s.Len & 1) << byte(ICON_BITS_SHIFT)
	s.Icons[i] |= icon << o
	s.Len += 1
	s.RelativeBits |= bit
}

func (s *SweepResult) AddBombUpdate(icon byte, bit uint64) {
	i := s.Len >> 1
	o := (s.Len & 1) << byte(ICON_BITS_SHIFT)
	s.Icons[i] |= icon << o
	s.Len += 1
	s.RelativeBits |= bit
}

func (s *SweepResult) DoActionOnAllTiles(action func(pos Coord, icon byte)) {
	if s.Len == 0 {
		return
	}
	remainingBits := s.RelativeBits
	center := s.Center.ToCoordInt()
	icon := s.Icons[0] & ICON_MASK
	action(center, icon)
	var idx byte = 1
	for idx < s.Len {
		bitIdx := bits.TrailingZeros64(remainingBits)
		bit := uint64(1) << bitIdx
		remainingBits &= ^bit
		pos := center.Add(NearCoordTable[bitIdx])
		iconIdx := idx >> 1
		iconOff := (idx & 1) << byte(ICON_BITS_SHIFT)
		icon = (s.Icons[iconIdx] >> iconOff) & ICON_MASK
		action(pos, icon)
		idx += 1
	}
}

func (s *SweepResult) WireWrite(w *wire.OutgoingWire) {
	w.TryWrite_U16(s.Score)
	s.Center.WireWrite(w)
	w.TryWrite_U64(s.RelativeBits)
	w.TryWrite_U8(s.Len)
	iconLen := (s.Len + 1) >> 1
	w.TryWrite_SliceU8(s.Icons[:iconLen])
}
func (s *SweepResult) WireRead(w *wire.IncomingWire) {
	w.TryRead_U16(&s.Score)
	s.Center.WireRead(w)
	w.TryRead_U64(&s.RelativeBits)
	w.TryRead_U8(&s.Len)
	iconLen := (s.Len + 1) >> 1
	w.TryRead_SliceU8(s.Icons[:iconLen])
}

var _ wire.WireReader = (*SweepResult)(nil)
var _ wire.WireWriter = (*SweepResult)(nil)

// type SweepResultIter struct {
// 	RemainingBits uint64
// 	Center        Coord
// 	RetCenter     bool
// }

// func (r *SweepResultIter) Next() (coord Coord, ok bool) {
// 	if !r.RetCenter {
// 		r.RetCenter = true
// 		return r.Center, true
// 	}
// 	if r.RemainingBits == 0 {
// 		return coord, false
// 	}
// 	idx := bits.TrailingZeros64(r.RemainingBits)
// 	bit := uint64(1) << idx
// 	r.RemainingBits &= ^bit
// 	coord = r.Center.Add(NearCoordTable[idx])
// 	return coord, true
// }
