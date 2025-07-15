package internal

const (
	INIT_FULL_U64 uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111
	INIT_NEAR_8   uint64 = bit0 | bit1 | bit2 | bit3 | bit4 | bit5 | bit6 | bit7
	INIT_NEAR_4   uint64 = bit1 | bit3 | bit5 | bit7
)
const (
	idx0 uint64 = iota
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
	idx_a
	idx_b
	idx_c
	idx_d
	idx_e
	idx_f
	idx_g
	idx_h
	idx_i
	idx_j
	idx_k
	idx_l
	idx_m
	idx_n
	idx_o
	idx_p
	idx_q
	idx_r
	idx_s
	idx_t
	idx_u
	idx_v
	idx_w
	idx_x
	idx_y
	idx_z
	idx_aa
	idx_bb
)

// +----------+
// |     y
// |    mno
// |  OPQRSTU
// |  l89ABCV
// | xkN012DWp
// bbwjM7.3EXqz
// | viL654FYr
// |  hKJIHGZ
// |  gfedcba
// |    uts
// |     aa
// +----------+

// const (
// 	idxNorthWest = idx0
// 	idxNorth     = idx1
// 	idxNorthEast = idx2
// 	idxEast      = idx3
// 	idxSouthEast = idx4
// 	idxSouth     = idx5
// 	idxSouthWest = idx6
// 	idxWest      = idx7

// 	bitNorthWest = bit0
// 	bitNorth     = bit1
// 	bitNorthEast = bit2
// 	bitEast      = bit3
// 	bitSouthEast = bit4
// 	bitSouth     = bit5
// 	bitSouthWest = bit6
// 	bitWest      = bit7
// )

// +----------+
// |     y
// |    mno
// |  OPQRSTU
// |  l89ABCV
// | xkN012DWp
// bbwjM7.3EXqz
// | viL654FYr
// |  hKJIHGZ
// |  gfedcba
// |    uts
// |     aa
// +----------+

const (
	bit_   uint64 = 0
	bit0   uint64 = 1 << idx0
	bit1   uint64 = 1 << idx1
	bit2   uint64 = 1 << idx2
	bit3   uint64 = 1 << idx3
	bit4   uint64 = 1 << idx4
	bit5   uint64 = 1 << idx5
	bit6   uint64 = 1 << idx6
	bit7   uint64 = 1 << idx7
	bit8   uint64 = 1 << idx8
	bit9   uint64 = 1 << idx9
	bitA   uint64 = 1 << idxA
	bitB   uint64 = 1 << idxB
	bitC   uint64 = 1 << idxC
	bitD   uint64 = 1 << idxD
	bitE   uint64 = 1 << idxE
	bitF   uint64 = 1 << idxF
	bitG   uint64 = 1 << idxG
	bitH   uint64 = 1 << idxH
	bitI   uint64 = 1 << idxI
	bitJ   uint64 = 1 << idxJ
	bitK   uint64 = 1 << idxK
	bitL   uint64 = 1 << idxL
	bitM   uint64 = 1 << idxM
	bitN   uint64 = 1 << idxN
	bitO   uint64 = 1 << idxO
	bitP   uint64 = 1 << idxP
	bitQ   uint64 = 1 << idxQ
	bitR   uint64 = 1 << idxR
	bitS   uint64 = 1 << idxS
	bitT   uint64 = 1 << idxT
	bitU   uint64 = 1 << idxU
	bitV   uint64 = 1 << idxV
	bitW   uint64 = 1 << idxW
	bitX   uint64 = 1 << idxX
	bitY   uint64 = 1 << idxY
	bitZ   uint64 = 1 << idxZ
	bit_a  uint64 = 1 << idx_a
	bit_b  uint64 = 1 << idx_b
	bit_c  uint64 = 1 << idx_c
	bit_d  uint64 = 1 << idx_d
	bit_e  uint64 = 1 << idx_e
	bit_f  uint64 = 1 << idx_f
	bit_g  uint64 = 1 << idx_g
	bit_h  uint64 = 1 << idx_h
	bit_i  uint64 = 1 << idx_i
	bit_j  uint64 = 1 << idx_j
	bit_k  uint64 = 1 << idx_k
	bit_l  uint64 = 1 << idx_l
	bit_m  uint64 = 1 << idx_m
	bit_n  uint64 = 1 << idx_n
	bit_o  uint64 = 1 << idx_o
	bit_p  uint64 = 1 << idx_p
	bit_q  uint64 = 1 << idx_q
	bit_r  uint64 = 1 << idx_r
	bit_s  uint64 = 1 << idx_s
	bit_t  uint64 = 1 << idx_t
	bit_u  uint64 = 1 << idx_u
	bit_v  uint64 = 1 << idx_v
	bit_w  uint64 = 1 << idx_w
	bit_x  uint64 = 1 << idx_x
	bit_y  uint64 = 1 << idx_y
	bit_z  uint64 = 1 << idx_z
	bit_aa uint64 = 1 << idx_aa
	bit_bb uint64 = 1 << idx_bb
)

// +----------+
// |     y
// |    mno
// |  OPQRSTU
// |  l89ABCV
// | xkN012DWp
// bbwjM7.3EXqz
// | viL654FYr
// |  hKJIHGZ
// |  gfedcba
// |    uts
// |     aa
// +----------+

var (
	NearCoordTable = [64]Coord{
		idx0: {X: -1, Y: -1},
		idx1: {X: +0, Y: -1},
		idx2: {X: +1, Y: -1},
		idx3: {X: +1, Y: +0},
		idx4: {X: +1, Y: +1},
		idx5: {X: +0, Y: +1},
		idx6: {X: -1, Y: +1},
		idx7: {X: -1, Y: +0},

		idx8: {X: -2, Y: -2},
		idx9: {X: -1, Y: -2},
		idxA: {X: +0, Y: -2},
		idxB: {X: +1, Y: -2},
		idxC: {X: +2, Y: -2},
		idxD: {X: +2, Y: -1},
		idxE: {X: +2, Y: +0},
		idxF: {X: +2, Y: +1},
		idxG: {X: +2, Y: +2},
		idxH: {X: +1, Y: +2},
		idxI: {X: +0, Y: +2},
		idxJ: {X: -1, Y: +2},
		idxK: {X: -2, Y: +2},
		idxL: {X: -2, Y: +1},
		idxM: {X: -2, Y: +0},
		idxN: {X: -2, Y: -1},

		idxO:  {X: -3, Y: -3},
		idxP:  {X: -2, Y: -3},
		idxQ:  {X: -1, Y: -3},
		idxR:  {X: +0, Y: -3},
		idxS:  {X: +1, Y: -3},
		idxT:  {X: +2, Y: -3},
		idxU:  {X: +3, Y: -3},
		idxV:  {X: +3, Y: -2},
		idxW:  {X: +3, Y: -1},
		idxX:  {X: +3, Y: +0},
		idxY:  {X: +3, Y: +1},
		idxZ:  {X: +3, Y: +2},
		idx_a: {X: +3, Y: +3},
		idx_b: {X: +2, Y: +3},
		idx_c: {X: +1, Y: +3},
		idx_d: {X: +0, Y: +3},
		idx_e: {X: -1, Y: +3},
		idx_f: {X: -2, Y: +3},
		idx_g: {X: -3, Y: +3},
		idx_h: {X: -3, Y: +2},
		idx_i: {X: -3, Y: +1},
		idx_j: {X: -3, Y: +0},
		idx_k: {X: -3, Y: -1},
		idx_l: {X: -3, Y: -2},

		// +----------+
		// |     y
		// |    mno
		// |  OPQRSTU
		// |  l89ABCV
		// | xkN012DWp
		// bbwjM7.3EXqz
		// | viL654FYr
		// |  hKJIHGZ
		// |  gfedcba
		// |    uts
		// |     aa
		// +----------+

		idx_m: {X: -1, Y: -4},
		idx_n: {X: +0, Y: -4},
		idx_o: {X: +1, Y: -4},

		idx_p: {X: +4, Y: -1},
		idx_q: {X: +4, Y: +0},
		idx_r: {X: +4, Y: +1},

		idx_s: {X: +1, Y: +4},
		idx_t: {X: +0, Y: +4},
		idx_u: {X: -1, Y: +4},

		idx_v: {X: -4, Y: +1},
		idx_w: {X: -4, Y: +0},
		idx_x: {X: -4, Y: -1},

		idx_y:  {X: +0, Y: -5},
		idx_z:  {X: +5, Y: +0},
		idx_aa: {X: +0, Y: +5},
		idx_bb: {X: -5, Y: +0},
	}
	// NearCoordTable8 = [8]Coord{
	// 	NearCoordTable[idxNorthWest], NearCoordTable[idxNorth], NearCoordTable[idxNorthEast],
	// 	NearCoordTable[idxWest] /*                          */, NearCoordTable[idxEast],
	// 	NearCoordTable[idxSouthWest], NearCoordTable[idxSouth], NearCoordTable[idxSouthEast],
	// }
	// NearBitTable8 = [8]uint64{
	// 	bitNorthWest, bitNorth, bitNorthEast,
	// 	bitWest /*          */, bitEast,
	// 	bitSouthWest, bitSouth, bitSouthEast,
	// }
	// +----------+
	// |     y
	// |    mno
	// |  OPQRSTU
	// |  l89ABCV
	// | xkN012DWp
	// bbwjM7.3EXqz
	// | viL654FYr
	// |  hKJIHGZ
	// |  gfedcba
	// |    uts
	// |     aa
	// +----------+

	NearBitsTable = [64]uint64{
		idx0: bit8 | bit9 | bitA | bit1 | bit7 | bitM | bitN,
		idx1: bit9 | bitA | bitB | bit2 | bit3 | bit7 | bit0,
		idx2: bitA | bitB | bitC | bitD | bitE | bit3 | bit1,
		idx3: bit1 | bit2 | bitD | bitE | bitF | bit4 | bit5,
		idx4: bit3 | bitE | bitF | bitG | bitH | bitI | bit5,
		idx5: bit3 | bit4 | bitH | bitI | bitJ | bit6 | bit7,
		idx6: bitM | bit7 | bit5 | bitI | bitJ | bitK | bitL,
		idx7: bitN | bit0 | bit1 | bit5 | bit6 | bitL | bitM,
		idx8: bitO | bitP | bitQ | bit9 | bit0 | bitN | bit_k | bit_l,
		idx9: bitP | bitQ | bitR | bitA | bit1 | bit0 | bitN | bit8,
		idxA: bitQ | bitR | bitS | bitB | bit2 | bit1 | bit0 | bit9,
		idxB: bitR | bitS | bitT | bitC | bitD | bit2 | bit1 | bitA,
		idxC: bitS | bitT | bitU | bitV | bitW | bitD | bit2 | bitB,
		idxD: bitB | bitC | bitV | bitW | bitX | bitE | bit3 | bit2,
		idxE: bit2 | bitD | bitW | bitX | bitY | bitF | bit4 | bit3,
		idxF: bit3 | bitE | bitX | bitY | bitZ | bitG | bitH | bit4,
		idxG: bit4 | bitF | bitY | bitZ | bit_a | bit_b | bit_c | bitH,
		idxH: bit4 | bit5 | bitF | bitG | bit_b | bit_c | bit_d | bitI,
		idxI: bit6 | bit5 | bit4 | bitH | bit_c | bit_d | bit_e | bitJ,
		idxJ: bitL | bit6 | bit5 | bitI | bit_d | bit_e | bit_f | bitK,
		idxK: bit_i | bitL | bit6 | bitJ | bit_e | bit_f | bit_g | bit_h,
		idxL: bit_j | bitM | bit7 | bit6 | bitJ | bitK | bit_h | bit_i,
		idxM: bit_k | bitN | bit0 | bit7 | bit6 | bitL | bit_i | bit_j,
		idxN: bit_l | bit8 | bit9 | bit0 | bit7 | bitM | bit_j | bit_k,

		// +----------+
		// |     y
		// |    mno
		// |  OPQRSTU
		// |  l89ABCV
		// | xkN012DWp
		// bbwjM7.3EXqz
		// | viL654FYr
		// |  hKJIHGZ
		// |  gfedcba
		// |    uts
		// |     aa
		// +----------+

		idxO:   bitP | bit8 | bit_l,
		idxP:   bit_m | bitQ | bit9 | bit8 | bit_l | bitO,
		idxQ:   bit_m | bit_n | bitR | bitA | bit9 | bit8 | bitP,
		idxR:   bit_m | bit_n | bit_o | bitS | bitB | bitA | bit9 | bitQ,
		idxS:   bit_n | bit_o | bitT | bitC | bitB | bitA | bitR,
		idxT:   bit_o | bitU | bitV | bitC | bitB | bitS,
		idxU:   bitV | bitC | bitT,
		idxV:   bitT | bitU | bit_p | bitW | bitD | bitC,
		idxW:   bitC | bitV | bit_p | bit_q | bitX | bitE | bitD,
		idxX:   bitD | bitW | bit_p | bit_q | bit_r | bitY | bitF | bitE,
		idxY:   bitE | bitX | bit_q | bit_r | bitZ | bitG | bitF,
		idxZ:   bitF | bitY | bit_r | bit_a | bit_b | bitG,
		idx_a:  bitG | bitZ | bit_b,
		idx_b:  bitH | bitG | bitZ | bit_a | bit_s | bit_c,
		idx_c:  bitI | bitH | bitG | bit_b | bit_s | bit_t | bit_d,
		idx_d:  bitJ | bitI | bitH | bit_c | bit_s | bit_t | bit_u | bit_e,
		idx_e:  bitK | bitJ | bitI | bit_d | bit_t | bit_u | bit_f,
		idx_f:  bit_h | bitK | bitJ | bit_e | bit_u | bit_g,
		idx_g:  bit_h | bitK | bit_f,
		idx_h:  bit_v | bit_i | bitL | bitK | bit_f | bit_g,
		idx_i:  bit_w | bit_j | bitM | bitL | bitK | bit_h | bit_v,
		idx_j:  bit_x | bit_k | bitN | bitM | bitL | bit_i | bit_v | bit_w,
		idx_k:  bit_l | bit8 | bitN | bitM | bit_j | bit_w | bit_x,
		idx_l:  bitO | bitP | bit8 | bitN | bit_k | bit_x,
		idx_m:  bit_y | bit_n | bitR | bitQ | bitP,
		idx_n:  bit_y | bit_o | bitS | bitR | bitQ | bit_m,
		idx_o:  bit_y | bitT | bitS | bitR | bit_n,
		idx_p:  bitV | bit_z | bit_q | bitX | bitW,
		idx_q:  bitW | bit_p | bit_z | bit_r | bitY | bitX,
		idx_r:  bitX | bit_q | bit_z | bitZ | bitY,
		idx_s:  bit_d | bit_c | bit_b | bit_aa | bit_t,
		idx_t:  bit_e | bit_d | bit_c | bit_s | bit_aa | bit_u,
		idx_u:  bit_f | bit_e | bit_d | bit_t | bit_aa,
		idx_v:  bit_bb | bit_w | bit_j | bit_i | bit_h,
		idx_w:  bit_x | bit_k | bit_j | bit_i | bit_v | bit_bb,
		idx_x:  bit_l | bit_k | bit_j | bit_w | bit_bb,
		idx_y:  bit_m | bit_n | bit_o,
		idx_z:  bit_p | bit_q | bit_r,
		idx_aa: bit_u | bit_t | bit_s,
		idx_bb: bit_x | bit_w | bit_v,
	}
	// +----------+
	// |     y
	// |    mno
	// |  OPQRSTU
	// |  l89ABCV
	// | xkN012DWp
	// bbwjM7.3EXqz
	// | viL654FYr
	// |  hKJIHGZ
	// |  gfedcba
	// |    uts
	// |     aa
	// +----------+

)

// type DebugRelCoords struct {
// 	Table [11][11]byte
// }

// func (d *DebugRelCoords) Init() {
// 	for y := range 11 {
// 		for x := range 11 {
// 			d.Table[y][x] = byte('_')
// 		}
// 	}
// }

func DoActionOn8NearbyCoords(center Coord, action func(nearPos Coord)) {
	queue := NewCascadeQueue(center)
	next, hasMoreNext := queue.NextToCheck()
	for hasMoreNext {
		action(next.Pos)
		next, hasMoreNext = queue.NextToCheck()
	}
}

func DoActionOn8NearbyCoordsInRange(center Coord, minX, maxX, minY, maxY int, action func(nearPos Coord, nearBit uint64)) {
	queue := NewCascadeQueue(center)
	next, hasMoreNext := queue.NextToCheck()
	for hasMoreNext {
		if next.Pos.IsInRange(minX, maxX, minY, maxY) {
			action(next.Pos, uint64(1)<<next.RelativeIdx)
		}
		next, hasMoreNext = queue.NextToCheck()
	}
}
