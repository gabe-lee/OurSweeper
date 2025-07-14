package internal

const (
	INIT_FULL_U64 uint64 = 0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111
	INIT_NEAR_8   uint64 = bit7 | bitNorth | bit1 | bitEast | bit3 | bitSouth | bit5 | bitWest
	INIT_NEAR_4   uint64 = bitNorth | bitEast | bitSouth | bitWest
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
)
const (
	idxNorthWest = idx7
	idxNorthEast = idx1
	idxSouthWest = idx5
	idxSouthEast = idx3

	bitNorthWest = bit7
	bitNorthEast = bit1
	bitSouthWest = bit5
	bitSouthEast = bit3
)

// +----------+
// |     a
// |    tKb
// |  xsZ8Lcu
// |  rYJ09Md
// | qXI7_1ANe
// |pWH6_._2BOf
// | oVG5_3CPg
// |  nUF4DQh
// |  wmTERiv
// |    lSj
// |     k
// +----------+

const (
	bitNone  uint64 = 0
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
	bit_a    uint64 = 1 << idx_a
	bit_b    uint64 = 1 << idx_b
	bit_c    uint64 = 1 << idx_c
	bit_d    uint64 = 1 << idx_d
	bit_e    uint64 = 1 << idx_e
	bit_f    uint64 = 1 << idx_f
	bit_g    uint64 = 1 << idx_g
	bit_h    uint64 = 1 << idx_h
	bit_i    uint64 = 1 << idx_i
	bit_j    uint64 = 1 << idx_j
	bit_k    uint64 = 1 << idx_k
	bit_l    uint64 = 1 << idx_l
	bit_m    uint64 = 1 << idx_m
	bit_n    uint64 = 1 << idx_n
	bit_o    uint64 = 1 << idx_o
	bit_p    uint64 = 1 << idx_p
	bit_q    uint64 = 1 << idx_q
	bit_r    uint64 = 1 << idx_r
	bit_s    uint64 = 1 << idx_s
	bit_t    uint64 = 1 << idx_t
	bit_u    uint64 = 1 << idx_u
	bit_v    uint64 = 1 << idx_v
	bit_w    uint64 = 1 << idx_w
	bit_x    uint64 = 1 << idx_x
)

// +----------+
// |     a
// |    tKb
// |  xsZ8Lcu
// |  rYJ09Md
// | qXI7_1ANe
// |pWH6_._2BOf
// | oVG5_3CPg
// |  nUF4DQh
// |  wmTERiv
// |    lSj
// |     k
// +----------+

var (
	NearCoordTable = [64]Coord{
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
	NearCoordTable8 = [8]Coord{
		NearCoordTable[idxNorthWest], NearCoordTable[idxNorth], NearCoordTable[idxNorthEast],
		NearCoordTable[idxWest] /*                          */, NearCoordTable[idxEast],
		NearCoordTable[idxSouthWest], NearCoordTable[idxSouth], NearCoordTable[idxSouthEast],
	}
	NearBitTable8 = [8]uint64{
		bitNorthWest, bitNorth, bitNorthEast,
		bitWest /*          */, bitEast,
		bitSouthWest, bitSouth, bitSouthEast,
	}
	// +----------+
	// |     a
	// |    tKb
	// |  xsZ8Lcu
	// |  rYJ09Md
	// | qXI7_1ANe
	// |pWH6_._2BOf
	// | oVG5_3CPg
	// |  nUF4DQh
	// |  wmTERiv
	// |    lSj
	// |     k
	// +----------+

	NearBitsTable = [64]uint64{
		idxNorth: bit7 | bit0 | bit1 | bitJ | bit9 | bitEast | bitWest | bitNone,
		idxEast:  bit1 | bit2 | bit3 | bitA | bitC | bitNorth | bitSouth | bitNone,
		idxSouth: bit3 | bit4 | bit5 | bitD | bitF | bitEast | bitWest | bitNone,
		idxWest:  bit5 | bit6 | bit7 | bitG | bitI | bitSouth | bitNorth | bitNone,
		idx0:     bitJ | bit8 | bit9 | bitZ | bitL | bit7 | bit1 | bitNorth,
		idx1:     bit9 | bitA | bitM | bit2 | bit0 | bitNorth | bitEast | bitNone,
		idx2:     bitA | bitB | bitC | bitN | bitP | bit1 | bit3 | bitEast,
		idx3:     bitC | bitD | bitQ | bit2 | bit4 | bitEast | bitSouth | bitNone,
		idx4:     bitD | bitE | bitF | bitR | bitT | bit3 | bit5 | bitSouth,
		idx5:     bitF | bitG | bitU | bit4 | bit6 | bitSouth | bitWest | bitNone,
		idx6:     bitG | bitH | bitI | bitV | bitX | bit7 | bit5 | bitWest,
		idx7:     bitI | bitJ | bitY | bit6 | bit0 | bitWest | bitNorth | bitNone,
		idx8:     bitZ | bitK | bitL | bitJ | bit0 | bit9,
		idx9:     bitL | bitM | bit0 | bit8 | bitA | bit1 | bitNorth,
		idxA:     bitM | bitN | bit1 | bit9 | bit2 | bitB | bitEast,
		idxB:     bitN | bitO | bitP | bitA | bit2 | bitC,
		idxC:     bitP | bitQ | bit2 | bitB | bitD | bit3 | bitEast,
		idxD:     bitQ | bitR | bit3 | bitC | bit4 | bitE | bitSouth,
		idxE:     bitR | bitS | bitT | bitF | bit4 | bitD,
		idxF:     bitT | bitU | bitG | bit5 | bit4 | bitE | bitSouth,
		idxG:     bitU | bitV | bitH | bit6 | bit5 | bitF | bitEast,
		idxH:     bitV | bitW | bitX | bitI | bit6 | bitG,
		idxI:     bitX | bitY | bitH | bit6 | bitJ | bit7 | bitWest,
		idxJ:     bitY | bitZ | bit8 | bit0 | bitI | bit7 | bitNorth,
		idxK:     bitZ | bit8 | bitL,
		idxL:     bitK | bit8 | bit0 | bit9 | bitM,
		idxM:     bitL | bit9 | bit1 | bitA | bitN,
		idxN:     bitM | bitA | bit2 | bitB | bitO,
		idxO:     bitN | bitB | bitP,
		idxP:     bit2 | bitB | bitO | bitC | bitQ,
		idxQ:     bit3 | bitC | bitP | bitD | bitR,
		idxR:     bit4 | bitD | bitQ | bitE | bitS,
		idxS:     bitT | bitE | bitR,
		idxT:     bitS | bitE | bit4 | bitU | bitF,
		idxU:     bitT | bitF | bit5 | bitG | bitV,
		idxV:     bitU | bitG | bit6 | bitE | bitH,
		idxW:     bitX | bitH | bitV,
		idxX:     bitW | bitH | bit6 | bitI | bitY,
		idxY:     bitX | bitI | bit7 | bitJ | bitZ,
		idxZ:     bitY | bitJ | bit0 | bit8 | bitK,
	}
	// +----------+
	// |     a
	// |    tKb
	// |  xsZ8Lcu
	// |  rYJ09Md
	// | qXI7_1ANe
	// |pWH6_._2BOf
	// | oVG5_3CPg
	// |  nUF4DQh
	// |  wmTERiv
	// |    lSj
	// |     k
	// +----------+

)
