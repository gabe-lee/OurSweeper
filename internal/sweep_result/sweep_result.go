package sweep_result

import (
	"math/bits"

	"github.com/gabe-lee/OurSweeper/coord"
	"github.com/gabe-lee/OurSweeper/data_buffer"
	C "github.com/gabe-lee/OurSweeper/internal/consts"
	"github.com/gabe-lee/OurSweeper/utils"
)

type (
	Writable    = data_buffer.Writable
	Readable    = data_buffer.Readable
	Sizable     = data_buffer.Sizable
	WriteBuffer = data_buffer.WriteBuffer
	ReadBuffer  = data_buffer.ReadBuffer
	ByteCoord   = coord.Coord[byte]
	Coord       = coord.Coord[int]
)

type SweepResult struct {
	Score        uint16
	Center       ByteCoord
	RelativeBits uint64
	Icons        [C.MAX_ICON_LEN]byte
	Len          byte
}

func (s *SweepResult) InitSweep(pos Coord, score uint16, icon byte) {
	s.Score = score
	s.Center = pos.ToCoordByte()
	s.Icons[0] = icon
	s.Len = 1
}

func (s *SweepResult) AddCascadeSweep(icon byte, bit uint64) {
	s.Score += uint16(C.BOMB_NEAR_BASE_SCORE[0])
	i := s.Len >> 1
	o := (s.Len & 1) << byte(C.ICON_BITS_SHIFT)
	s.Icons[i] |= icon << o
	s.Len += 1
	s.RelativeBits |= bit
}

func (s *SweepResult) AddBombUpdate(icon byte, bit uint64) {
	i := s.Len >> 1
	o := (s.Len & 1) << byte(C.ICON_BITS_SHIFT)
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
	icon := s.Icons[0] & C.ICON_MASK
	action(center, icon)
	var idx byte = 1
	for idx < s.Len {
		bitIdx := bits.TrailingZeros64(remainingBits)
		bit := uint64(1) << bitIdx
		remainingBits &= ^bit
		pos := center.Add(C.NearCoordTable[bitIdx])
		iconIdx := idx >> 1
		iconOff := (idx & 1) << byte(C.ICON_BITS_SHIFT)
		icon = (s.Icons[iconIdx] >> iconOff) & C.ICON_MASK
		action(pos, icon)
		idx += 1
	}
}

func (s *SweepResult) WriteToBuffer(buf *WriteBuffer) {
	buf.U16_LE(s.Score)
	buf.U8(s.Center.X)
	buf.U8(s.Center.Y)
	buf.U64_LE(s.RelativeBits)
	buf.U8(s.Len)
	iconLen := (s.Len + 1) >> 1
	buf.U8_Slice(s.Icons[:iconLen])
}

func (s *SweepResult) ReadFromBuffer(buf *ReadBuffer) error {
	e := utils.FirstError{}
	e.Add(buf.U16_LE(&s.Score))
	e.Add(buf.U8(&s.Center.X))
	e.Add(buf.U8(&s.Center.Y))
	e.Add(buf.U64_LE(&s.RelativeBits))
	e.Add(buf.U8(&s.Len))
	iconLen := (s.Len + 1) >> 1
	e.Add(buf.U8_Slice(s.Icons[:iconLen]))
	return e.Err
}

func (s *SweepResult) SizeOnBuffer() int {
	return 13 + int((s.Len+1)>>1)
}

var _ Readable = (*SweepResult)(nil)
var _ Writable = (*SweepResult)(nil)
var _ Sizable = (*SweepResult)(nil)
