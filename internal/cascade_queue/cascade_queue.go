package cascade_queue

import (
	"github.com/gabe-lee/OurSweeper/internal/coord"
)

type Coord = coord.Coord

const (
	// NORTH = iota
	// EAST
	// SOUTH
	// WEST
	INIT_CONT uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111
	INIT_NEXT uint64 = 0b1111
)
const (
	idxNorth uint64 = iota
	idxEast
	idxSouth
	idxWest
	idx0
	idx1
	idx2
	idx3
	idx4
	idx5
	idx6
	idx7
	idx8
	idx9
	idxA
	idxB
	idxC
	idxD
	idxE
	idxF
	idxG
	idxH
	idxI
	idxJ
	idxK
	idxL
	idxM
	idxN
	idxO
	idxP
	idxQ
	idxR
	idxS
	idxT
	idxU
	idxV
	idxW
	idxX
	idxY
	idxZ
)
const (
	bitNorth uint64 = 1 << idxNorth
	bitEast  uint64 = 1 << idxEast
	bitSouth uint64 = 1 << idxSouth
	bitWest  uint64 = 1 << idxWest
	bit0     uint64 = 1 << idx0
	bit1     uint64 = 1 << idx1
	bit2     uint64 = 1 << idx2
	bit3     uint64 = 1 << idx3
	bit4     uint64 = 1 << idx4
	bit5     uint64 = 1 << idx5
	bit6     uint64 = 1 << idx6
	bit7     uint64 = 1 << idx7
	bit8     uint64 = 1 << idx8
	bit9     uint64 = 1 << idx9
	bitA     uint64 = 1 << idxA
	bitB     uint64 = 1 << idxB
	bitC     uint64 = 1 << idxC
	bitD     uint64 = 1 << idxD
	bitE     uint64 = 1 << idxE
	bitF     uint64 = 1 << idxF
	bitG     uint64 = 1 << idxG
	bitH     uint64 = 1 << idxH
	bitI     uint64 = 1 << idxI
	bitJ     uint64 = 1 << idxJ
	bitK     uint64 = 1 << idxK
	bitL     uint64 = 1 << idxL
	bitM     uint64 = 1 << idxM
	bitN     uint64 = 1 << idxN
	bitO     uint64 = 1 << idxO
	bitP     uint64 = 1 << idxP
	bitQ     uint64 = 1 << idxQ
	bitR     uint64 = 1 << idxR
	bitS     uint64 = 1 << idxS
	bitT     uint64 = 1 << idxT
	bitU     uint64 = 1 << idxU
	bitV     uint64 = 1 << idxV
	bitW     uint64 = 1 << idxW
	bitX     uint64 = 1 << idxX
	bitY     uint64 = 1 << idxY
	bitZ     uint64 = 1 << idxZ
)
const (
	addNorth uint64 = bit7 | bit0 | bit1
	addEast  uint64 = bit1 | bit2 | bit3
	addSouth uint64 = bit3 | bit4 | bit5
	addWest  uint64 = bit5 | bit6 | bit7
)

// +----------+
// |     a
// |    tKb
// |   sZ8Lc
// |  rYJ09Md
// | qXI7_1ANe
// |pWH6_._2BOf
// | oVG5_3CPg
// |  nUF4DQh
// |   mTERi
// |    lSj
// |     k
// +----------+
var (
	CoordTable = [40]Coord{
		idxNorth: {X: 0, Y: -1},
		idxEast:  {X: 1, Y: 0},
		idxSouth: {X: 0, Y: 1},
		idxWest:  {X: -1, Y: 0},
		idx0:     {X: 0, Y: -2},
		idx1:     {X: 1, Y: -1},
		idx2:     {X: 2, Y: 0},
		idx3:     {X: 1, Y: 1},
		idx4:     {X: 0, Y: 2},
		idx5:     {X: -1, Y: 1},
		idx6:     {X: -2, Y: 0},
		idx7:     {X: -1, Y: -1},
		idx8:     {X: 0, Y: -3},
		idx9:     {X: 1, Y: -2},
		idxA:     {X: 2, Y: -1},
		idxB:     {X: 3, Y: 0},
		idxC:     {X: 2, Y: 1},
		idxD:     {X: 1, Y: 2},
		idxE:     {X: 0, Y: 3},
		idxF:     {X: -1, Y: 2},
		idxG:     {X: -2, Y: 1},
		idxH:     {X: -3, Y: 0},
		idxI:     {X: -2, Y: -1},
		idxJ:     {X: -1, Y: -2},
		idxK:     {X: 0, Y: -4},
		idxL:     {X: 1, Y: -3},
		idxM:     {X: 2, Y: -2},
		idxN:     {X: 3, Y: -1},
		idxO:     {X: 4, Y: 0},
		idxP:     {X: 3, Y: 1},
		idxQ:     {X: 2, Y: 2},
		idxR:     {X: 1, Y: 3},
		idxS:     {X: 0, Y: 4},
		idxT:     {X: -1, Y: 3},
		idxU:     {X: -2, Y: 2},
		idxV:     {X: -3, Y: 1},
		idxW:     {X: -4, Y: 0},
		idxX:     {X: -3, Y: -1},
		idxY:     {X: -2, Y: -2},
		idxZ:     {X: -1, Y: -3},
	}
	// +----------+
	// |
	// |     a
	// |    tKb
	// |   sZ8Lc
	// |  rYJ09Md
	// | qXI7_1ANe
	// |pWH6_._2BOf
	// | oVG5_3CPg
	// |  nUF4DQh
	// |   mTERi
	// |    lSj
	// |     k
	// |Layer 0 = 4, TOTAL=4
	// |Layer 1 = 8, TOTAL=12
	// |Layer 2 = 12, TOTAL=24
	// |Layer 3 = 16, TOTAL=40
	// |Layer 4 = 20, TOTAL=60
	// +----------+
	AddTable = [40]uint64{
		idxNorth: bit7 | bit0 | bit1,
		idxEast:  bit1 | bit2 | bit3,
		idxSouth: bit3 | bit4 | bit5,
		idxWest:  bit5 | bit6 | bit7,
		idx0:     bitJ | bit8 | bit9,
		idx1:     bit9 | bitA,
		idx2:     bitA | bitB | bitC,
		idx3:     bitC | bitD,
		idx4:     bitD | bitE | bitF,
		idx5:     bitF | bitG,
		idx6:     bitG | bitH | bitI,
		idx7:     bitI | bitJ,
		idx8:     bitZ | bitK | bitL,
		idx9:     bitL | bitM,
		idxA:     bitM | bitN,
		idxB:     bitN | bitO | bitP,
		idxC:     bitP | bitQ,
		idxD:     bitQ | bitR,
		idxE:     bitR | bitS | bitT,
		idxF:     bitT | bitU,
		idxG:     bitU | bitV,
		idxH:     bitV | bitW | bitX,
		idxI:     bitX | bitY,
		idxJ:     bitY | bitZ,
		idxK:     0,
		idxL:     0,
		idxM:     0,
		idxN:     0,
		idxO:     0,
		idxP:     0,
		idxQ:     0,
		idxR:     0,
		idxS:     0,
		idxT:     0,
		idxU:     0,
		idxV:     0,
		idxW:     0,
		idxX:     0,
		idxY:     0,
		idxZ:     0,
	}
)

type CascadeQueue struct {
	NextList     uint64
	ContinueMask uint64
	Center       Coord
	Idx          int
}

func New(center Coord) CascadeQueue {
	queue := CascadeQueue{
		NextList:     INIT_NEXT,
		ContinueMask: INIT_CONT,
		Center:       center,
	}
	return queue
}

type CascadeTile struct {
	Pos Coord
	Idx int
}

// func (q *CascadeQueue) ShouldContinue() bool {
// 	return q.NextList&q.ContinueMask > 0
// }

func (q *CascadeQueue) Next() (tile CascadeTile, ok bool) {
	anyExists := q.NextList&q.ContinueMask > 0
	if !anyExists {
		return CascadeTile{}, false
	}
	exists := q.NextList&(uint64(1)<<q.Idx) > 0
	for !exists {
		q.Idx += 1
		q.ContinueMask <<= 1
		exists = q.NextList&(uint64(1)<<q.Idx) > 0
	}
	item := CascadeTile{
		Idx: q.Idx,
		Pos: q.Center.Add(CoordTable[q.Idx]),
	}
	q.Idx += 1
	q.ContinueMask <<= 1
	return item, true
}

func (q *CascadeQueue) Cascade(item CascadeTile) {
	q.NextList |= AddTable[item.Idx]
}

// func (q *CascadeQueue) QueueIfNew(pos Coord) {
// 	for i := 0; i < q.End; i += 1 {
// 		oldPos := q.Arr[i]
// 		if pos.Equals(oldPos) {
// 			return
// 		}
// 	}
// 	q.Queue(pos)
// }

// func (q *CascadeQueue) Queue(pos Coord) {
// 	q.Arr[q.End] = pos
// 	q.End += 1
// }

// func (q *CascadeQueue) Dequeue() (pos Coord, ok bool) {
// 	if q.Start == q.End {
// 		return Coord{}, false
// 	}
// 	val := q.Arr[q.Start]
// 	q.Start += 1
// 	q.Start = q.Start % C.CASCADE_BUFFER_LEN
// 	return val, true
// }
