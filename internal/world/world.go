package world

import (
	"io"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/gabe-lee/OurSweeper/internal/ansi"
	"github.com/gabe-lee/OurSweeper/internal/lockset"
	"github.com/gabe-lee/OurSweeper/internal/tile"
	"github.com/gabe-lee/OurSweeper/internal/xmath"
)

const (
	WIDTH          uint = 64
	HEIGHT         uint = 64
	HALF_WIDTH     uint = WIDTH / 2
	HALF_HEIGHT    uint = WIDTH / 2
	TILES          uint = WIDTH * HEIGHT
	LOCKS          uint = 64
	TILES_PER_LOCK uint = TILES / LOCKS
	AXIS_PER_LOCK  uint = 8 // square root of TILES_PER_LOCK
	T_TO_L         uint = 3 // bits to shift DOWN to turn Tile X or Y into Lock X or Y

	MIN_MINE_CHANCE float64 = 0.15
	MAX_MINE_CHANCE float64 = 0.30

	INITIAL_OPAQUE_TILES uint = TILES - (2 * WIDTH) - ((2 * HEIGHT) - 4)

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
)

var MAX_DIST_FROM_CENTER float64 = math.Sqrt((CENTER_X * CENTER_X) + (CENTER_Y + CENTER_Y))

type World struct {
	Id    atomic.Uint32
	Tiles [TILES]tile.Tile
	Locks [LOCKS]sync.Mutex

	RemainingMines atomic.Uint32
	RemainingSpots atomic.Uint32
}

func mineChance(x uint, y uint) float64 {
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

func GetIndex(x uint, y uint) uint {
	return (y << 6) | x
}

func GetLockIndex(xl uint, yl uint) uint {
	return (yl << T_TO_L) | xl
}

func GetMineChance(x uint, y uint) float64 {
	if x == 0 || x == WIDTH-1 || y == 0 || y == HEIGHT-1 {
		return 0.0
	}
	return float64(mineChance(x, y))
}

func (w *World) LockTile(set *lockset.LockSet, x, y uint) bool {
	lockX, lockY := x>>T_TO_L, y>>T_TO_L
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

func (w *World) UnlockTiles(set *lockset.LockSet) {
	for y := set.YMin; y <= set.YMax; y++ {
		for x := set.XMin; x <= set.XMax; x++ {
			idx := (uint64(y) * 8) + uint64(x)
			bit := uint64(1) << idx
			if set.Locks&bit > 0 {
				w.Locks[idx].Unlock()
			}
		}
	}
	*set = lockset.New()
}

func (w *World) initMine(idx, x, y uint) uint {
	thresh := GetMineChance(x, y)
	randVal := rand.Float64()
	var result uint = 0
	if randVal <= thresh {
		w.Tiles[idx].SetMine()
		result = 1
	}
	return result
}

func (w *World) initNearby(idx, x, y uint) {
	var total uint8 = 0
	yMin := uint(max(int(y)-1, 0))
	yMax := uint(min(int(y)+1, int(WIDTH-1)))
	xMin := uint(max(int(x)-1, 0))
	xMax := uint(min(int(x)+1, int(HEIGHT-1)))
	for yy := yMin; yy <= yMax; yy++ {
		for xx := xMin; xx <= xMax; xx++ {
			if xx == x && yy == y {
				continue
			}
			nearIdx := GetIndex(xx, yy)
			if w.Tiles[nearIdx].IsMine() {
				total++
			}
		}
	}
	w.Tiles[idx].SetNearby(total)
}

func (w *World) InitNew(id uint32) {
	w.Id.Store(id)
	w.RemainingSpots.Store(uint32(INITIAL_OPAQUE_TILES))
	var y uint = 0
	var x uint = 0
	var totalMines uint
	for y < HEIGHT {
		x = 0
		for x < WIDTH {
			idx := GetIndex(x, y)
			w.Tiles[idx] = tile.Tile(0)
			if x > 0 && x < WIDTH-1 && y > 0 && y < HEIGHT-1 {
				totalMines += w.initMine(idx, x, y)
			}
			x++
		}
		y++
	}
	w.RemainingMines.Store(uint32(totalMines))
	y = 0
	x = 0
	for y < HEIGHT {
		x = 0
		for x < WIDTH {
			idx := GetIndex(x, y)
			w.initNearby(idx, x, y)
			x++
		}
		y++
	}
	x = 0
	xx := TILES - WIDTH
	for range WIDTH {
		w.Tiles[x].SetViz(tile.VIZ_SWEPT)
		w.Tiles[xx].SetViz(tile.VIZ_SWEPT)
		x++
		xx++
	}
	y = WIDTH
	yy := WIDTH - 1
	for range HEIGHT - 1 {
		w.Tiles[y].SetViz(tile.VIZ_SWEPT)
		w.Tiles[yy].SetViz(tile.VIZ_SWEPT)
		y += WIDTH
		yy += WIDTH
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
			case tile.ICON_BOMB:
				wr.Write([]byte(ICON_SKULL))
			case tile.ICON_OPAQUE:
				wr.Write([]byte(ICON_OPAQUE))
			case tile.ICON_FLAG:
				wr.Write([]byte(ICON_FLAG))
			case tile.ICON_0:
				wr.Write([]byte(ICON_0))
			case tile.ICON_1:
				wr.Write([]byte(ICON_1))
			case tile.ICON_2:
				wr.Write([]byte(ICON_2))
			case tile.ICON_3:
				wr.Write([]byte(ICON_3))
			case tile.ICON_4:
				wr.Write([]byte(ICON_4))
			case tile.ICON_5:
				wr.Write([]byte(ICON_5))
			case tile.ICON_6:
				wr.Write([]byte(ICON_6))
			case tile.ICON_7:
				wr.Write([]byte(ICON_7))
			case tile.ICON_8:
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
