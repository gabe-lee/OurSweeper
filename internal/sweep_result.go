package internal

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/bits"

	"github.com/gabe-lee/OurSweeper/serializer"
	"github.com/gabe-lee/OurSweeper/utils"
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
	s.Center = pos.ToByteCoord()
	s.Icons[0] = icon
	s.Len = 1
}

func (s *SweepResult) AddScoreAndIcon(score uint16, icon byte) {
	s.Score += score
	i := s.Len >> 1
	o := s.Len & 1
	o |= o << 1
	o |= o << 2
	s.Icons[i] |= icon << o
	s.Len += 1
}

func (s *SweepResult) DoActionOnAllTiles(action func(pos Coord, icon byte)) {

	if s.Len == 0 {
		return
	}
	remainingBits := s.RelativeBits
	center := s.Center.ToCoord()
	icon := s.Icons[0] & ICON_MASK
	action(center, icon)
	var idx byte = 1
	for idx < s.Len {
		bitIdx := bits.TrailingZeros64(remainingBits)
		bit := uint64(1) << bitIdx
		remainingBits &= ^bit
		pos := center.Add(NearCoordTable[bitIdx])
		iconIdx := idx >> 1
		iconOff := idx & 1
		iconOff |= iconOff << 1
		iconOff |= iconOff << 2
		icon = (s.Icons[iconIdx] >> iconOff) & ICON_MASK
		action(pos, icon)
		idx += 1
	}
}

// func (s SweepResult) Iter() SweepResultIter {
// 	return SweepResultIter{
// 		Center:        s.Center.ToCoord(),
// 		RemainingBits: s.RelativeBits,
// 		RetCenter:     s.Len > 0,
// 	}
// }

func (s *SweepResult) CodedSerialize(order binary.ByteOrder) (data []byte, err error) {
	e := utils.ErrorCollector{}
	data = make([]byte, 0, 44)
	w := bytes.NewBuffer(data)
	e.Do(binary.Write(w, order, SERVER_SWEEP))
	e.Do(binary.Write(w, order, s.Score))
	e.Do(s.Center.Serialize(w, order))
	e.Do(binary.Write(w, order, s.RelativeBits))
	e.Do(binary.Write(w, order, s.Len))
	iconLen := (s.Len + 1) >> 1
	e.Do(binary.Write(w, order, s.Icons[:iconLen]))
	return w.Bytes(), e.Err
}
func (s *SweepResult) Deserialize(r io.Reader, order binary.ByteOrder) error {
	e := utils.ErrorCollector{}
	e.Do(binary.Read(r, order, &s.Score))
	e.Do(s.Center.Deserialize(r, order))
	e.Do(binary.Read(r, order, &s.RelativeBits))
	e.Do(binary.Read(r, order, &s.Len))
	iconLen := (s.Len + 1) >> 1
	e.Do(binary.Read(r, order, s.Icons[:iconLen]))
	return e.Err
}

var _ serializer.CodedSerializer = (*SweepResult)(nil)
var _ serializer.Deserializer = (*SweepResult)(nil)

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
