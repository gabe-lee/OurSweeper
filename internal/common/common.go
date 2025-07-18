package common

import (
	"encoding/binary"
	"math"
	"time"
)

const (
	WORLD_TILE_WIDTH           int = 64
	WORLD_TILE_HEIGHT          int = 64
	WORLD_MAX_X                int = WORLD_TILE_WIDTH - 1
	WORLD_MAX_Y                int = WORLD_TILE_HEIGHT - 1
	WORLD_HALF_WIDTH           int = WORLD_TILE_WIDTH / 2
	WORLD_HALF_HEIGHT          int = WORLD_TILE_HEIGHT / 2
	WORLD_TILE_COUNT           int = WORLD_TILE_WIDTH * WORLD_TILE_HEIGHT
	WORLD_TILES_PER_CHUNK_AXIS int = 8
	WORLD_CHUNK_COUNT          int = WORLD_CHUNK_HEIGHT * WORLD_CHUNK_WIDTH
	WORLD_CHUNK_WIDTH          int = WORLD_TILE_WIDTH / WORLD_TILES_PER_CHUNK_AXIS
	WORLD_CHUNK_HEIGHT         int = WORLD_TILE_HEIGHT / WORLD_TILES_PER_CHUNK_AXIS
	TILES_PER_CHUNK            int = WORLD_TILE_COUNT / WORLD_CHUNK_COUNT
	CY_SHIFT                   int = 3 // trailing zeros of the square root of TILES_PER_CHUNK
	CX_MASK                    int = (1 << CY_SHIFT) - 1
	TILE_TO_LOCK_SHIFT         int = TY_SHIFT - LY_SHIFT // bits to shift DOWN to turn Tile X or Y into Lock X or Y
	TY_SHIFT                   int = 6                   // bits to shift tile `y` value up/down when calculating index
	LY_SHIFT                   int = 3                   // bits to shift lock `y` value up/down when calculating index
	TX_MASK                    int = WORLD_TILE_WIDTH - 1
	LX_MASK                    int = WORLD_TILES_PER_CHUNK_AXIS - 1

	WORLD_LAST_ROW = WORLD_TILE_COUNT - WORLD_TILE_WIDTH

	MIN_MINE_CHANCE float64 = 0.15
	MAX_MINE_CHANCE float64 = 0.25

	INITIAL_OPAQUE_TILES int = WORLD_TILE_COUNT - INITIAL_SWEPT_TILES
	INITIAL_SWEPT_TILES  int = (2 * WORLD_TILE_WIDTH) + ((2 * WORLD_TILE_HEIGHT) - 4)

	WORLD_CENTER_X float64 = float64(WORLD_HALF_WIDTH)
	WORLD_CENTER_Y float64 = float64(WORLD_HALF_HEIGHT)

	RISC_EXPONENT_ADD   float64 = 1.125
	RISC_EXPONENT_MULT  float64 = 1.0
	RISC_FULL_BLIND     float64 = 1.050 // Calculated so a full blind sweep results in 150 pts
	IDX_FULL_BLIND_BASE uint8   = 9

	WORLD_TIME_LIMIT time.Duration = time.Hour * 24

	MAX_CASCDE_DIST    int = 4
	CASCADE_BUFFER_LEN int = 41

	

	

	MAX_CASCADE_LEN int = 64
	MAX_SWEEP_LEN   int = MAX_CASCADE_LEN + 1
	MAX_ICON_LEN    int = (MAX_SWEEP_LEN + 1) >> 1

	MAX_CLIENT_MSG_LEN int = 128
)

var BYTE_ORDER = binary.LittleEndian

var (
	MAX_DIST_FROM_CENTER float64 = math.Sqrt((WORLD_CENTER_X * WORLD_CENTER_X) + (WORLD_CENTER_Y + WORLD_CENTER_Y))

	BOMB_NEAR_BASE_SCORE = [10]float64{
		//0  1    2    3     4     5     6     7     8    Full Blind
		1.0, 5.0, 7.0, 10.0, 15.0, 21.0, 34.0, 50.0, 1.0, 10.0,
	}

	
)
