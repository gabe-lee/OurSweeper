package internal

import (
	"encoding/binary"
	"math"
	"time"
)

const (
	WORLD_TILE_WIDTH  int = 64
	WORLD_TILE_HEIGHT int = 64
	WORLD_MAX_X       int = WORLD_TILE_WIDTH - 1
	WORLD_MAX_Y       int = WORLD_TILE_HEIGHT - 1
	WORLD_HALF_WIDTH  int = WORLD_TILE_WIDTH / 2
	WORLD_HALF_HEIGHT int = WORLD_TILE_HEIGHT / 2
	WORLD_TILE_COUNT  int = WORLD_TILE_WIDTH * WORLD_TILE_HEIGHT
	WORLD_LOCK_COUNT  int = 64
	TILES_PER_LOCK    int = WORLD_TILE_COUNT / WORLD_LOCK_COUNT
	AXIS_PER_LOCK     int = 8
	T_TO_L            int = 3 // bits to shift DOWN to turn Tile X or Y into Lock X or Y
	TY_SHIFT          int = 6 // bits to shift tile `y` value up/down when calculating index
	LY_SHIFT          int = 3 // bits to shift lock `y` value up/down when calculating index
	TX_MASK           int = WORLD_TILE_WIDTH - 1
	LX_MASK           int = AXIS_PER_LOCK - 1

	WORLD_LAST_ROW = WORLD_TILE_COUNT - WORLD_TILE_WIDTH

	MIN_MINE_CHANCE float64 = 0.15
	MAX_MINE_CHANCE float64 = 0.30

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

	TILE_SIZE        int = 32
	TILE_SIZE_SCALED int = TILE_SIZE / DISPLAY_SCALE_DOWN

	TILE_SHEET_WIDTH  int = 12
	TILE_SHEET_HEIGHT int = 4

	WINDOW_WIDTH         int     = 800
	WINDOW_HEIGHT        int     = 800
	BOARD_WIDTH          int     = TILE_SIZE_SCALED * WORLD_TILE_WIDTH
	BOARD_HEIGHT         int     = TILE_SIZE_SCALED * WORLD_TILE_HEIGHT
	BOARD_OVERFLOW_X     int     = BOARD_WIDTH - WINDOW_WIDTH
	BOARD_OVERFLOW_Y     int     = BOARD_HEIGHT - WINDOW_HEIGHT
	MIN_BOARD_POS_X      float64 = float64(-BOARD_OVERFLOW_X)
	MIN_BOARD_POS_Y      float64 = float64(-BOARD_OVERFLOW_Y)
	MAX_BOARD_POS_X      float64 = 0
	MAX_BOARD_POS_Y      float64 = 0
	DISPLAY_SCALE_DOWN   int     = 2
	DISPLAY_SCALE_DOWN_F float64 = float64(DISPLAY_SCALE_DOWN)
	WHEEL_SPEED          float64 = 6.0

	MAX_CASCADE_LEN int = 40
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

	// Cascade Position Chart
	// +----------+
	// |
	// |  * 6IU *
	// |   5TAJV
	// |  4SH_BKW
	// |  RG_*_CL
	// |  3QF_DMX
	// |   2PENY
	// |  * 1OZ *
	// |
	// +----------+
	NW_OFF = [2]int{-3, -3}
	NE_OFF = [2]int{3, -3}
	SW_OFF = [2]int{-3, 3}
	SE_OFF = [2]int{3, 3}

	OFFSET_NORTH = Coord{X: 0, Y: -1}
	OFFSET_EAST  = Coord{X: 1, Y: 0}
	OFFSET_SOUTH = Coord{X: 0, Y: 1}
	OFFSET_WEST  = Coord{X: -1, Y: 0}

	CASCADE_NORTH = [3]Coord{
		OFFSET_WEST,
		OFFSET_NORTH,
		OFFSET_EAST,
	}
	CASCADE_EAST = [3]Coord{
		OFFSET_NORTH,
		OFFSET_EAST,
		OFFSET_SOUTH,
	}
	CASCADE_SOUTH = [3]Coord{
		OFFSET_EAST,
		OFFSET_SOUTH,
		OFFSET_WEST,
	}
	CASCADE_WEST = [3]Coord{
		OFFSET_SOUTH,
		OFFSET_WEST,
		OFFSET_NORTH,
	}

	BOARD_TILES = [16][2]int{
		ICON_CODE_0:      {0, 0},
		ICON_CODE_1:      {1 * TILE_SIZE, 0},
		ICON_CODE_2:      {2 * TILE_SIZE, 0},
		ICON_CODE_3:      {3 * TILE_SIZE, 0},
		ICON_CODE_4:      {4 * TILE_SIZE, 0},
		ICON_CODE_5:      {5 * TILE_SIZE, 0},
		ICON_CODE_6:      {6 * TILE_SIZE, 0},
		ICON_CODE_7:      {7 * TILE_SIZE, 0},
		ICON_CODE_8:      {8 * TILE_SIZE, 0},
		ICON_CODE_FLAG:   {9 * TILE_SIZE, 0},
		ICON_CODE_BOMB:   {10 * TILE_SIZE, 0},
		ICON_CODE_OPAQUE: {11 * TILE_SIZE, 0},
	}
)
