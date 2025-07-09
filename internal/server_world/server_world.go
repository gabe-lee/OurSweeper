package server_world

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/gabe-lee/OurSweeper/internal/ansi"
	"github.com/gabe-lee/OurSweeper/internal/lockset"
	"github.com/gabe-lee/OurSweeper/internal/sweep_result"
	"github.com/gabe-lee/OurSweeper/internal/tile"
	"github.com/gabe-lee/OurSweeper/internal/utils"
	"github.com/gabe-lee/OurSweeper/internal/xmath"
)

const (
	WIDTH          int = 64
	HEIGHT         int = 64
	HALF_WIDTH     int = WIDTH / 2
	HALF_HEIGHT    int = WIDTH / 2
	TILES          int = WIDTH * HEIGHT
	LOCKS          int = 64
	TILES_PER_LOCK int = TILES / LOCKS
	AXIS_PER_LOCK  int = 8 // square root of TILES_PER_LOCK
	T_TO_L         int = 3 // bits to shift DOWN to turn Tile X or Y into Lock X or Y
	TY_SHIFT       int = 6 // bits to shift tile `y` value up/down when calculating index
	LY_SHIFT       int = 3 // bits to shift lock `y` value up/down when calculating index

	LAST_ROW = TILES - WIDTH

	MIN_MINE_CHANCE float64 = 0.15
	MAX_MINE_CHANCE float64 = 0.30

	INITIAL_OPAQUE_TILES int = TILES - INITIAL_SWEPT_TILES
	INITIAL_SWEPT_TILES  int = (2 * WIDTH) + ((2 * HEIGHT) - 4)

	CENTER_X float64 = float64(HALF_WIDTH)
	CENTER_Y float64 = float64(HALF_HEIGHT)

	ICON_0      = " "
	ICON_1      = ansi.FG_BLU + "1" + ansi.CLEAR
	ICON_2      = ansi.FG_CYA + "2" + ansi.CLEAR
	ICON_3      = ansi.FG_GRN + "3" + ansi.CLEAR
	ICON_4      = ansi.FG_YEL + "4" + ansi.CLEAR
	ICON_5      = ansi.FG_RED + "5" + ansi.CLEAR
	ICON_6      = ansi.FG_MAG + "6" + ansi.CLEAR
	ICON_7      = ansi.FG_WHT + "7" + ansi.CLEAR
	ICON_8      = ansi.FG_BLK + "8" + ansi.CLEAR
	ICON_BOMB   = ansi.INV_RED + "X" + ansi.CLEAR
	ICON_SKULL  = ansi.FG_RED + "💀" + ansi.CLEAR
	ICON_FLAG   = "█"
	ICON_OPAQUE = "▒"

	MAX_TILE_CASCADE int = 36

	RISC_EXPONENT_ADD   float64 = 1.125
	RISC_EXPONENT_MULT  float64 = 1.0
	RISC_FULL_BLIND     float64 = 1.050 // Calculated so a full blind results in 150 pts
	IDX_FULL_BLIND_BASE uint8   = 9

	WORLD_TIME_LIMIT time.Duration = time.Hour * 24
)

var BOMB_NEAR_BASE_SCORE = [10]float64{
	//0  1    2    3     4     5     6     7     8    Full Blind
	1.0, 5.0, 7.0, 10.0, 15.0, 21.0, 34.0, 50.0, 1.0, 10.0,
}

var MAX_DIST_FROM_CENTER float64 = math.Sqrt((CENTER_X * CENTER_X) + (CENTER_Y + CENTER_Y))

var (
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

	CARDINAL_ADD = [4]uint64{
		0b01_01_00_00_00_00_00_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //NORTH
		0b00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //EAST
		0b00_00_00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //SOUTH
		0b00_00_00_00_00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //WEST
	}

	CARDINAL_OFF = [4][2]int{
		{0, -1}, //NORTH
		{1, 0},  //EAST
		{0, 1},  //SOUTH
		{-1, 0}, //WEST
	}

	CASCADE_ADD = [20]uint64{
		0b00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_01_00_00_00_00_00_00_00_00_00_00_00_00, //A
		0b00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //B
		0b00_00_00_00_00_00_00_00_00_00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //C
		0b00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //D
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //E
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00, //F
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_01_00_00_00_00_00_00_00_00_00_00_00_00_00, //G
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_00_00_00, //H
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_00_00_00_00_00_00_00_00_00_00_01, //I
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00_00, //J
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00_00, //K
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00_00, //L
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00_00, //M
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00_00, //N
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00_00, //O
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00_00, //P
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00_00, //Q
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00_00, //R
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01_00, //S
		0b00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_00_01_01, //T
	}
	CASCADE_OFF = [32][2]int{
		{0, -2},  //A=0
		{1, -1},  //B=1
		{2, 0},   //C=2
		{1, 1},   //D=3
		{0, 2},   //E=4
		{-1, 1},  //F=5
		{-2, 0},  //G=6
		{-1, 1},  //H=7
		{0, -3},  //I=8
		{1, -2},  //J=9
		{2, -1},  //K=10
		{3, 0},   //L=11
		{2, 1},   //M=12
		{1, 2},   //N=13
		{0, 3},   //O=14
		{-1, 2},  //P=15
		{-2, 1},  //Q=16
		{-3, 0},  //R=17
		{-2, -1}, //S=18
		{-1, -2}, //T=19
		{1, -3},  //U=20
		{2, -2},  //V=21
		{3, -1},  //W=22
		{3, 1},   //X=23
		{2, 2},   //Y=24
		{1, 3},   //Z=25
		{-1, 3},  //1=26
		{-2, 2},  //2=27
		{-3, 1},  //3=28
		{-3, -1}, //4=29
		{-2, -2}, //5=30
		{-1, -3}, //6=31
	}
)

type World struct {
	Id            atomic.Uint32
	Tiles         [TILES]tile.Tile
	Locks         [LOCKS]sync.Mutex
	TotalMines    uint32
	ExplodedMines atomic.Uint32
	SweptTiles    atomic.Uint32
	Ended         atomic.Bool
	Expires       time.Time
}

func mineChance(x int, y int) float64 {
	fx := float64(x)
	fy := float64(y)
	var dx float64
	var dy float64
	if x > HALF_WIDTH {
		dx = fx - CENTER_X
	} else {
		dx = CENTER_X - fx
	}
	if y > HALF_HEIGHT {
		dy = fy - CENTER_Y
	} else {
		dy = CENTER_Y - fy
	}
	dist := math.Sqrt((dx * dx) + (dy * dy))
	percent := (MAX_DIST_FROM_CENTER - dist) / MAX_DIST_FROM_CENTER
	return xmath.Lerp(MIN_MINE_CHANCE, MAX_MINE_CHANCE, percent)
}

func GetIndex(x int, y int) int {
	return (y << TY_SHIFT) | x
}

func GetLockIndex(xl int, yl int) int {
	return (yl << LY_SHIFT) | xl
}

func tileCoordToLockCoord(x, y int) (lx, ly int) {
	lx = x >> T_TO_L
	ly = y >> T_TO_L
	return
}

func GetMineChance(x int, y int) float64 {
	if x == 0 || x == WIDTH-1 || y == 0 || y == HEIGHT-1 {
		return 0.0
	}
	return float64(mineChance(x, y))
}

func (w *World) LockEntireWorld() {
	for i := range LOCKS {
		w.Locks[i].Lock()
	}
}
func (w *World) UnlockEntireWorld() {
	for i := range LOCKS {
		w.Locks[i].Unlock()
	}
}

func (w *World) AquireQuadLock(topLeftX, topLeftY, botRightX, botRightY int) lockset.QuadLock {
	nwX, nwY := tileCoordToLockCoord(topLeftX, topLeftY)
	neX, neY := tileCoordToLockCoord(botRightX, topLeftY)
	swX, swY := tileCoordToLockCoord(topLeftX, botRightY)
	seX, seY := tileCoordToLockCoord(botRightX, botRightY)
	locks := lockset.QuadLock{}
	locks.LockIndexes[0] = GetLockIndex(nwX, nwY)
	w.Locks[locks.LockIndexes[0]].Lock()
	i := 1
	if neX != nwX {
		locks.LockIndexes[i] = GetLockIndex(neX, neY)
		w.Locks[locks.LockIndexes[i]].Lock()
		i += 1
	}
	if swY != nwY {
		locks.LockIndexes[i] = GetLockIndex(swX, swY)
		w.Locks[locks.LockIndexes[i]].Lock()
		i += 1
	}
	if seX != nwX || seY != nwY {
		locks.LockIndexes[i] = GetLockIndex(seX, seY)
		w.Locks[locks.LockIndexes[i]].Lock()
		i += 1
	}
	locks.Len = i
	return locks
}

func (w *World) ReleaseQuadLock(locks lockset.QuadLock) {
	for idx := range locks.LockIndexes[:locks.Len] {
		w.Locks[idx].Unlock()
	}
}

func (w *World) TryLockTile(set *lockset.LockSet, x, y int) bool {
	lockX, lockY := tileCoordToLockCoord(x, y)
	lockIdx := GetLockIndex(lockX, lockY)
	if set.AlreadyLocked(lockIdx) {
		return true
	}
	didLock := w.Locks[lockIdx].TryLock()
	if didLock {
		set.AddLock(lockIdx, lockX, lockY)
	}
	return didLock
}

// func (w *World) LockTile(set *lockset.LockSet, lockIdx, lx, ly int) {
// 	if set.AlreadyLocked(lockIdx) {
// 		return
// 	}
// 	w.Locks[lockIdx].Lock()
// 	set.AddLock(lockIdx, lx, ly)
// }

// func (w *World) UnlockTiles(set *lockset.LockSet) {
// 	for y := set.YMin; y <= set.YMax; y++ {
// 		for x := set.XMin; x <= set.XMax; x++ {
// 			idx := (uint64(y) * 8) + uint64(x)
// 			bit := uint64(1) << idx
// 			if set.Locks&bit > 0 {
// 				w.Locks[idx].Unlock()
// 			}
// 		}
// 	}
// 	*set = lockset.LockSet{}
// }

func (w *World) initMine(idx, x, y int) uint32 {
	thresh := GetMineChance(x, y)
	randVal := rand.Float64()
	var result uint32 = 0
	if randVal <= thresh {
		w.Tiles[idx].SetMine()
		result = 1
	}
	return result
}

func (w *World) initNearby(idx, x, y int) {
	var total uint8 = 0
	yMin := max(y-1, 0)
	yMax := min(y+1, int(WIDTH-1))
	xMin := max(y-1, 0)
	xMax := min(y+1, int(HEIGHT-1))
	for yy := yMin; yy <= yMax; yy++ {
		for xx := xMin; xx <= xMax; xx++ {
			if xx != x || yy != y {
				nearIdx := GetIndex(xx, yy)
				if w.Tiles[nearIdx].IsMine() {
					total++
				}
			}
		}
	}
	w.Tiles[idx].SetNearby(total)
}

func (w *World) InitNew(id uint32) {
	w.Id.Store(id)
	w.SweptTiles.Store(uint32(INITIAL_SWEPT_TILES))
	w.ExplodedMines.Store(0)
	w.Ended.Store(false)
	w.Expires = time.Now().Add(WORLD_TIME_LIMIT)
	for y := range HEIGHT {
		for x := range WIDTH {
			idx := GetIndex(x, y)
			w.Tiles[idx] = tile.Tile(0)
			if x > 0 && x < WIDTH-1 && y > 0 && y < HEIGHT-1 {
				w.TotalMines += w.initMine(idx, x, y)
			}
		}
	}
	for y := range HEIGHT {
		for x := range WIDTH {
			idx := GetIndex(x, y)
			w.initNearby(idx, x, y)
		}
	}
	for idx := range WIDTH {
		w.Tiles[idx].SetViz(tile.VIZ_EMPTY)
		w.Tiles[idx+LAST_ROW].SetViz(tile.VIZ_EMPTY)
	}
	for idx := range HEIGHT - 2 {
		w.Tiles[WIDTH+(WIDTH*idx)].SetViz(tile.VIZ_EMPTY)
		w.Tiles[WIDTH+(WIDTH*idx)+WIDTH-1].SetViz(tile.VIZ_EMPTY)
	}
}

func (w *World) SweepTile(x, y int) sweep_result.SweepResult {
	result := sweep_result.SweepResult{}
	nwX, nwY := utils.AddOffset(x, y, NW_OFF)
	seX, seY := utils.AddOffset(x, y, SE_OFF)
	locks := w.AquireQuadLock(max(nwX, 0), max(nwY, 0), min(seX, WIDTH-1), min(seY, HEIGHT-1))
	defer w.ReleaseQuadLock(locks)
	tileIdx := GetIndex(x, y)
	t := w.Tiles[tileIdx]
	if t.IsSwept() {
		return result
	}
	isMine := t.IsMine()
	if isMine {
		t.SetViz(tile.VIZ_BOMB)
		w.ExplodedMines.Add(1)
	} else {
		t.SetViz(tile.VIZ_EMPTY)
	}
	result.AddTile(w.getScore(x, y), t.GetIcon(), x, y)
	if isMine {
		w.reduceNearbyBombCounts(&result, x, y)
	} else if t.GetNearby() == 0 {
		w.cascade(&result, x, y)
	}
	if !isMine {
		w.SweptTiles.Add(uint32(result.Len))
	}
	w.Tiles[tileIdx] = t
	w.checkEndState()
	return result
}

func (w *World) checkEndState() {
	if exploded := w.ExplodedMines.Load(); exploded == w.TotalMines {
		w.Ended.Store(true)
	} else if swept := w.SweptTiles.Load(); swept == uint32(TILES) {
		w.Ended.Store(true)
	} else if time.Now().After(w.Expires) {
		w.Ended.Store(true)
	}
}

func (w *World) getScore(x, y int) uint32 {
	var lowestNearBombChance float64 = RISC_FULL_BLIND
	var lowestNearBombChanceBombs uint8 = IDX_FULL_BLIND_BASE
	xx := max(x-1, 0)
	xxMax := max(x+1, WIDTH-1)
	yy := max(y-1, 0)
	yyMax := max(y+1, HEIGHT-1)
	for yy < yyMax {
		for xx < xxMax {
			if xx != x || yy != y {
				thisIdx := GetIndex(xx, yy)
				if w.Tiles[thisIdx].IsSwept() {
					thisBombChance, thisBombs := w.getBombProbabilityAndNearby(xx, yy)
					if thisBombChance < lowestNearBombChance {
						lowestNearBombChance = thisBombChance
						lowestNearBombChanceBombs = thisBombs
					}
				}
			}
			xx += 1
		}
		yy += 1
	}
	scoreFloat := BOMB_NEAR_BASE_SCORE[lowestNearBombChanceBombs]
	if lowestNearBombChance < 1.0 {
		exp := RISC_EXPONENT_ADD + (RISC_EXPONENT_MULT * lowestNearBombChance)
		scoreFloat = math.Ceil(math.Pow(scoreFloat, exp))
	}
	return uint32(scoreFloat)
}

func (w *World) getBombProbabilityAndNearby(x, y int) (safe float64, near uint8) {
	idx := GetIndex(x, y)
	bombs := w.Tiles[idx].GetNearby()
	xx := max(x-1, 0)
	xxMax := max(x+1, WIDTH-1)
	yy := max(y-1, 0)
	yyMax := max(y+1, HEIGHT-1)
	var opaques float64
	for yy < yyMax {
		for xx < xxMax {
			if xx != x || yy != y {
				nearIdx := GetIndex(xx, yy)
				if !w.Tiles[nearIdx].IsSwept() {
					opaques += 1.0
				}
			}
			xx += 1
		}
		yy += 1
	}
	return float64(bombs) / opaques, bombs
}

func (w *World) reduceNearbyBombCounts(result *sweep_result.SweepResult, x, y int) {
	xx := max(x-1, 0)
	xxMax := max(x+1, WIDTH-1)
	yy := max(y-1, 0)
	yyMax := max(y+1, HEIGHT-1)
	for yy < yyMax {
		for xx < xxMax {
			if xx != x || yy != y {
				nearIdx := GetIndex(xx, yy)
				bombs := w.Tiles[nearIdx].GetNearby()
				w.Tiles[nearIdx].SetNearby(bombs - 1)
				if w.Tiles[nearIdx].IsSwept() {
					result.AddTile(0, w.Tiles[nearIdx].GetIcon(), xx, yy)
				}
			}
			xx += 1
		}
		yy += 1
	}
}

func (w *World) cascade(result *sweep_result.SweepResult, x, y int) {
	var casc uint64
	i := 0
	for i < 4 {
		add := CARDINAL_ADD[i]
		off := CARDINAL_OFF[i]
		xx, yy := utils.AddOffset(x, y, off)
		if xx >= 0 && xx < WIDTH && yy >= 0 && yy < HEIGHT {
			thisIdx := GetIndex(xx, yy)
			w.Tiles[thisIdx].SetViz(tile.VIZ_EMPTY)
			if w.Tiles[thisIdx].GetNearby() == 0 {
				casc += add
			}
			result.AddTile(uint32(BOMB_NEAR_BASE_SCORE[0]), w.Tiles[thisIdx].GetIcon(), xx, yy)
		}
		i += 1
	}
	var mask uint64 = 0b11
	i = 0
	for i < 20 && casc > 0 {
		if casc&mask > 0 {
			add := CASCADE_ADD[i]
			off := CASCADE_OFF[i]
			xx, yy := utils.AddOffset(x, y, off)
			if xx >= 0 && xx < WIDTH && yy >= 0 && yy < HEIGHT {
				thisIdx := GetIndex(xx, yy)
				w.Tiles[thisIdx].SetViz(tile.VIZ_EMPTY)
				if w.Tiles[thisIdx].GetNearby() == 0 {
					casc += add
				}
				result.AddTile(uint32(BOMB_NEAR_BASE_SCORE[0]), w.Tiles[thisIdx].GetIcon(), xx, yy)
			}
			casc &= ^mask
		}
		mask <<= 2
		i += 1
	}
	for i < 32 && casc > 0 {
		if casc&mask > 0 {
			off := CASCADE_OFF[i]
			xx, yy := utils.AddOffset(x, y, off)
			if xx >= 0 && xx < WIDTH && yy >= 0 && yy < HEIGHT {
				thisIdx := GetIndex(xx, yy)
				w.Tiles[thisIdx].SetViz(tile.VIZ_EMPTY)
				result.AddTile(uint32(BOMB_NEAR_BASE_SCORE[0]), w.Tiles[thisIdx].GetIcon(), xx, yy)
			}
			casc &= ^mask
		}
		mask <<= 2
		i += 1
	}
}

// This is NOT safe for concurrent use when world is being played on,
// this is only for testing/debugging purposes
func (w *World) DrawState(wr io.Writer) {
	var i uint = 0
	var buf [4]byte
	capFill := strings.Repeat("═", int(WIDTH))
	var cap string = "╔" + capFill + "╗\n"
	wr.Write([]byte(cap))
	for range HEIGHT {
		n := utf8.EncodeRune(buf[:], '║')
		wr.Write(buf[:n])
		for range WIDTH {
			iconCode := w.Tiles[i].GetIcon()
			switch iconCode {
			case tile.ICON_CODE_BOMB:
				wr.Write([]byte(ICON_SKULL))
			case tile.ICON_CODE_OPAQUE:
				wr.Write([]byte(ICON_OPAQUE))
			case tile.ICON_CODE_FLAG:
				wr.Write([]byte(ICON_FLAG))
			case tile.ICON_CODE_0:
				wr.Write([]byte(ICON_0))
			case tile.ICON_CODE_1:
				wr.Write([]byte(ICON_1))
			case tile.ICON_CODE_2:
				wr.Write([]byte(ICON_2))
			case tile.ICON_CODE_3:
				wr.Write([]byte(ICON_3))
			case tile.ICON_CODE_4:
				wr.Write([]byte(ICON_4))
			case tile.ICON_CODE_5:
				wr.Write([]byte(ICON_5))
			case tile.ICON_CODE_6:
				wr.Write([]byte(ICON_6))
			case tile.ICON_CODE_7:
				wr.Write([]byte(ICON_7))
			case tile.ICON_CODE_8:
				wr.Write([]byte(ICON_8))
			default:
				wr.Write([]byte(ICON_0))
			}
			i++
		}
		n = utf8.EncodeRune(buf[:], '║')
		buf[n] = 0x0A
		wr.Write(buf[:n+1])
	}
	cap = "╚" + capFill + "╝\n"
	wr.Write([]byte(cap))
}

// This is NOT safe for cuncurrent use when world is being played on,
// this is only for testing/debugging purposes
func (w *World) DrawMines(wr io.Writer) {
	var i uint = 0
	var buf [4]byte
	capFill := strings.Repeat("═", int(WIDTH))
	var cap string = "╔" + capFill + "╗\n"
	wr.Write([]byte(cap))
	for range HEIGHT {
		n := utf8.EncodeRune(buf[:], '║')
		wr.Write(buf[:n])
		for range WIDTH {
			isMine := w.Tiles[i].IsMine()
			if isMine {
				wr.Write([]byte(ICON_BOMB))
			} else {
				wr.Write([]byte(ICON_0))
			}
			i++
		}
		n = utf8.EncodeRune(buf[:], '║')
		buf[n] = 0x0A
		wr.Write(buf[:n+1])
	}
	cap = "╚" + capFill + "╝\n"
	wr.Write([]byte(cap))
}

// This is NOT safe for cuncurrent use when world is being played on,
// this is only for testing/debugging purposes
func (w *World) DrawNearby(wr io.Writer) {
	var i uint = 0
	var buf [4]byte
	capFill := strings.Repeat("═", int(WIDTH))
	var cap string = "╔" + capFill + "╗\n"
	wr.Write([]byte(cap))
	for range HEIGHT {
		n := utf8.EncodeRune(buf[:], '║')
		wr.Write(buf[:n])
		for range WIDTH {
			isMine := w.Tiles[i].IsMine()
			if isMine {
				wr.Write([]byte(ICON_BOMB))
			} else {
				nearby := w.Tiles[i].GetNearby()
				switch nearby {
				case 0:
					wr.Write([]byte(ICON_0))
				case 1:
					wr.Write([]byte(ICON_1))
				case 2:
					wr.Write([]byte(ICON_2))
				case 3:
					wr.Write([]byte(ICON_3))
				case 4:
					wr.Write([]byte(ICON_4))
				case 5:
					wr.Write([]byte(ICON_5))
				case 6:
					wr.Write([]byte(ICON_6))
				case 7:
					wr.Write([]byte(ICON_7))
				case 8:
					wr.Write([]byte(ICON_8))
				default:
					wr.Write([]byte(ICON_0))
				}
			}
			i++
		}
		n = utf8.EncodeRune(buf[:], '║')
		buf[n] = 0x0A
		wr.Write(buf[:n+1])
	}
	cap = "╚" + capFill + "╝\n"
	wr.Write([]byte(cap))
}

func (w *World) PrintStatus(wr io.Writer) {
	exploed := w.ExplodedMines.Load()
	swept := w.SweptTiles.Load()
	fmt.Fprintf(wr, "World ID: %d\n-Total Mines: %d\n-Exploded Mines: %d\n-Total Tiles: %d\n-SweptTiles: %d\n-Mine Completion: %.2f%%\n-Tile Completion: %.2f%%\n-Time Remaining: %s\n", w.Id.Load(), w.TotalMines, exploed, TILES, swept, float32(exploed)*100.0/float32(w.TotalMines), float32(swept)*100.0/float32(TILES), time.Until(w.Expires))
}
