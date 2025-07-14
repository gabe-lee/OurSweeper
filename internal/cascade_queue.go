package internal

import (
	"math/bits"
)

// type RelativeCoordBlock struct {
// 	Bits   uint64
// 	Center Coord
// }

// func NewRelativeCoordBlock(center Coord) RelativeCoordBlock {
// 	return RelativeCoordBlock{
// 		Center: center,
// 	}
// }

// func (r *RelativeCoordBlock) InitNear8() {
// 	r.Bits = INIT_NEAR_8
// }

// func (r *RelativeCoordBlock) InitNear4() {
// 	r.Bits = INIT_NEAR_4
// }

// func (r *RelativeCoordBlock) AddCascade(idx uint64) {
// 	r.Bits |= NearBitsTable[idx]
// }

// func (r *RelativeCoordBlock) AddBit(bit uint64) {
// 	r.Bits |= bit
// }

// func (r *RelativeCoordBlock) GetCoord(idx int) Coord {
// 	return r.Center.Add(CoordTable[idx])
// }

type RelativeCoord struct {
	Pos Coord
	Bit uint64
	Idx uint64
}

type CascadeQueue struct {
	ToCheckList  uint64
	DidSweepList uint64
	Idx          uint64
	Center       Coord
}

func NewCascadeQueue(center Coord) CascadeQueue {
	queue := CascadeQueue{
		ToCheckList:  INIT_NEAR_8,
		DidSweepList: 0,
		Center:       center,
		Idx:          0,
	}
	return queue
}

type CascadeCoord struct {
	Pos         Coord
	RelativeIdx uint64
}

func (q *CascadeQueue) NextToCheck() (coord CascadeCoord, ok bool) {
	nextMask := INIT_FULL_U64 << q.Idx
	if q.ToCheckList&nextMask == 0 {
		return coord, false
	}
	coord.RelativeIdx = uint64(bits.TrailingZeros64(q.ToCheckList))
	bit := uint64(1) << coord.RelativeIdx
	q.ToCheckList &= ^bit
	coord.Pos = q.Center.Add(NearCoordTable[coord.RelativeIdx])
	return coord, true
}

func (q *CascadeQueue) AddSweep(coord CascadeCoord) {
	q.DidSweepList |= uint64(1) << uint64(coord.RelativeIdx)
}

func (q *CascadeQueue) AddCascade(coord CascadeCoord) {
	checkBits := NearBitsTable[coord.RelativeIdx]
	q.ToCheckList |= checkBits
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
