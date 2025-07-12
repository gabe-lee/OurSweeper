package server_world

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gabe-lee/OurSweeper/internal/cascade_queue"
	C "github.com/gabe-lee/OurSweeper/internal/common"
	"github.com/gabe-lee/OurSweeper/internal/coord"
	"github.com/gabe-lee/OurSweeper/internal/lockset"
	"github.com/gabe-lee/OurSweeper/internal/sweep_result"
	"github.com/gabe-lee/OurSweeper/internal/tile"
	"github.com/gabe-lee/OurSweeper/internal/xmath"
)

type (
	Coord        = coord.Coord
	CascadeQueue = cascade_queue.CascadeQueue
	CascadeTile  = cascade_queue.CascadeTile
	Bounds4      = coord.Bounds4
	QuadLock     = lockset.QuadLock
)

type World struct {
	Id            atomic.Uint32
	Tiles         [C.WORLD_TILE_COUNT]tile.Tile
	Locks         [C.WORLD_LOCK_COUNT]sync.Mutex
	TotalMines    uint32
	ExplodedMines atomic.Uint32
	SweptTiles    atomic.Uint32
	Ended         atomic.Bool
	Expires       time.Time
}

func mineChance(pos Coord) float64 {
	fx := float64(pos.X)
	fy := float64(pos.Y)
	var dx float64
	var dy float64
	if pos.X > C.WORLD_HALF_WIDTH {
		dx = fx - C.WORLD_CENTER_X
	} else {
		dx = C.WORLD_CENTER_X - fx
	}
	if pos.Y > C.WORLD_HALF_HEIGHT {
		dy = fy - C.WORLD_CENTER_Y
	} else {
		dy = C.WORLD_CENTER_Y - fy
	}
	dist := math.Sqrt((dx * dx) + (dy * dy))
	percent := (C.MAX_DIST_FROM_CENTER - dist) / C.MAX_DIST_FROM_CENTER
	return xmath.Lerp(C.MIN_MINE_CHANCE, C.MAX_MINE_CHANCE, percent)
}

func (w *World) LockEntireWorld() {
	for i := range C.WORLD_LOCK_COUNT {
		w.Locks[i].Lock()
	}
}
func (w *World) UnlockEntireWorld() {
	for i := range C.WORLD_LOCK_COUNT {
		w.Locks[i].Unlock()
	}
}

func (w *World) AquireQuadLock(bounds Bounds4) QuadLock {
	lockBounds := bounds.ShiftDownScalar(C.T_TO_L)
	topLeftIdx := lockBounds.TopLeft.ToIndex(C.LY_SHIFT)
	topRightIdx := lockBounds.TopRight.ToIndex(C.LY_SHIFT)
	botLeftIdx := lockBounds.BotLeft.ToIndex(C.LY_SHIFT)
	botRightIdx := lockBounds.BotRight.ToIndex(C.LY_SHIFT)
	locks := lockset.QuadLock{}
	w.Locks[topLeftIdx].Lock()
	locks.LockIndexes[0] = topLeftIdx
	locks.Len = 1
	if topRightIdx != topLeftIdx {
		w.Locks[topRightIdx].Lock()
		locks.LockIndexes[1] = topRightIdx
		locks.Len = 2
	}
	if botLeftIdx != topLeftIdx && botLeftIdx != topRightIdx {
		w.Locks[botLeftIdx].Lock()
		locks.LockIndexes[locks.Len] = botLeftIdx
		locks.Len += 1
	}
	if botRightIdx != topLeftIdx && botRightIdx != topRightIdx && botRightIdx != botLeftIdx {
		w.Locks[botRightIdx].Lock()
		locks.LockIndexes[locks.Len] = botRightIdx
		locks.Len += 1
	}
	return locks
}

func (w *World) ReleaseQuadLock(locks lockset.QuadLock) {
	for _, idx := range locks.LockIndexes[:locks.Len] {
		w.Locks[idx].Unlock()
	}
}

// func (w *World) TryLockTile(set *C.WORLD_LOCK_COUNTet.C.WORLD_LOCK_COUNTet, x, y int) bool {
// 	lockX, lockY := tileCoordToLockCoord(x, y)
// 	lockIdx := GetLockIndex(lockX, lockY)
// 	if set.AlreadyLocked(lockIdx) {
// 		return true
// 	}
// 	didLock := w.Locks[lockIdx].TryLock()
// 	if didLock {
// 		set.AddLock(lockIdx, lockX, lockY)
// 	}
// 	return didLock
// }

// func (w *World) LockTile(set *C.WORLD_LOCK_COUNTet.C.WORLD_LOCK_COUNTet, lockIdx, lx, ly int) {
// 	if set.AlreadyLocked(lockIdx) {
// 		return
// 	}
// 	w.Locks[lockIdx].Lock()
// 	set.AddLock(lockIdx, lx, ly)
// }

// func (w *World) UnlockTiles(set *C.WORLD_LOCK_COUNTet.C.WORLD_LOCK_COUNTet) {
// 	for y := set.YMin; y <= set.YMax; y++ {
// 		for x := set.XMin; x <= set.XMax; x++ {
// 			idx := (uint64(y) * 8) + uint64(x)
// 			bit := uint64(1) << idx
// 			if set.C.WORLD_LOCK_COUNT&bit > 0 {
// 				w.Locks[idx].Unlock()
// 			}
// 		}
// 	}
// 	*set = C.WORLD_LOCK_COUNTet.C.WORLD_LOCK_COUNTet{}
// }

func (w *World) initMine(idx int, pos Coord) bool {
	thresh := mineChance(pos)
	randVal := rand.Float64()
	if randVal <= thresh {
		w.Tiles[idx].SetMine()
		return true
	}
	return false
}

// func (w *World) initNearby(idx, x, y int) {
// 	var total uint8 = 0
// 	utils.DoFuncOnNearbyCoords()
// 	yMin := max(y-1, 0)
// 	yMax := min(y+1, int(C.C.WORLD_TILE_WIDTH-1))
// 	xMin := max(y-1, 0)
// 	xMax := min(y+1, int(C.WORLD_TILE_HEIGHT-1))
// 	for yy := yMin; yy <= yMax; yy++ {
// 		for xx := xMin; xx <= xMax; xx++ {
// 			if xx != x || yy != y {
// 				nearIdx := GetIndex(xx, yy)
// 				if w.Tiles[nearIdx].IsMine() {
// 					total++
// 				}
// 			}
// 		}
// 	}
// 	w.Tiles[idx].SetNearby(total)
// }

func (w *World) InitNew(id uint32) {
	w.Id.Store(id)
	w.SweptTiles.Store(uint32(C.INITIAL_SWEPT_TILES))
	w.ExplodedMines.Store(0)
	w.Ended.Store(false)
	w.Expires = time.Now().Add(C.WORLD_TIME_LIMIT)
	for idx := range C.WORLD_TILE_COUNT {
		w.Tiles[idx] = tile.Tile(0)
	}
	for idx := range C.WORLD_LOCK_COUNT {
		w.Locks[idx] = sync.Mutex{}
	}

	for idx := range C.WORLD_TILE_COUNT {
		pos := coord.FromIndex(idx, C.TY_SHIFT, C.TX_MASK)
		if pos.IsInRangeExcludeEdges(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y) {
			mine := w.initMine(idx, pos)
			if mine {
				w.TotalMines += 1
				nears := pos.GetNearbyCoords(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y)
				for _, nearPos := range nears.Coords[:nears.Len] {
					nearIdx := nearPos.ToIndex(C.TY_SHIFT)
					w.Tiles[nearIdx].IncrNearbyMineCount()
				}
			}
		} else {
			w.Tiles[idx].SetVizSweptEmpty()
		}
	}
}

func (w *World) SweepTile(pos Coord) sweep_result.SweepResult {
	result := sweep_result.SweepResult{}
	tileIdx := pos.ToIndex(C.TY_SHIFT)
	t := w.Tiles[tileIdx]
	if t.IsSwept() {
		return result
	}
	bounds := pos.GetBounds4(Coord{X: C.MAX_CASCDE_DIST, Y: C.MAX_CASCDE_DIST}, 0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y)
	locks := w.AquireQuadLock(bounds)
	defer w.ReleaseQuadLock(locks)
	isMine := t.IsMine()
	if isMine {
		t.SetVizSweptBomb()
		w.ExplodedMines.Add(1)
	} else {
		t.SetVizSweptEmpty()
	}
	result.AddTile(w.getScore(pos), t.GetIconServer(), pos)
	if isMine {
		w.reduceNearbyBombCounts(&result, pos)
	} else if t.GetNearby() == 0 {
		w.cascade(&result, pos)
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
	} else if swept := w.SweptTiles.Load(); swept == uint32(C.WORLD_TILE_COUNT) {
		w.Ended.Store(true)
	} else if time.Now().After(w.Expires) {
		w.Ended.Store(true)
	}
}

func (w *World) getScore(pos Coord) uint16 {
	var lowestNearBombChance float64 = C.RISC_FULL_BLIND
	var lowestNearBombChanceBombs uint8 = C.IDX_FULL_BLIND_BASE
	nearbyCoords := pos.GetNearbyCoords(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y)
	for _, nearPos := range nearbyCoords.Coords[:nearbyCoords.Len] {
		nearIdx := nearPos.ToIndex(C.TY_SHIFT)
		if w.Tiles[nearIdx].IsSwept() {
			thisBombChance, thisBombs := w.getBombProbabilityAndNearby(nearPos)
			if thisBombChance < lowestNearBombChance {
				lowestNearBombChance = thisBombChance
				lowestNearBombChanceBombs = thisBombs
			}
		}
	}
	scoreFloat := C.BOMB_NEAR_BASE_SCORE[lowestNearBombChanceBombs]
	if lowestNearBombChance < 1.0 {
		exp := C.RISC_EXPONENT_ADD + (C.RISC_EXPONENT_MULT * lowestNearBombChance)
		scoreFloat = math.Ceil(math.Pow(scoreFloat, exp))
	}
	return uint16(scoreFloat)
}

func (w *World) getBombProbabilityAndNearby(pos Coord) (safe float64, near uint8) {
	var opaques float64
	bombs := w.Tiles[pos.ToIndex(C.TY_SHIFT)].GetNearby()
	nearbyCoords := pos.GetNearbyCoords(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y)
	for _, nearPos := range nearbyCoords.Coords[:nearbyCoords.Len] {
		nearIdx := nearPos.ToIndex(C.TY_SHIFT)
		if !w.Tiles[nearIdx].IsSwept() {
			opaques += 1.0
		}
	}
	return float64(bombs) / opaques, bombs
}

func (w *World) reduceNearbyBombCounts(result *sweep_result.SweepResult, pos Coord) {
	nearbyCoords := pos.GetNearbyCoords(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y)
	for _, nearPos := range nearbyCoords.Coords[:nearbyCoords.Len] {
		nearIdx := nearPos.ToIndex(C.TY_SHIFT)
		bombs := w.Tiles[nearIdx].GetNearby()
		w.Tiles[nearIdx].SetNearby(bombs - 1)
		if w.Tiles[nearIdx].IsSwept() {
			result.AddTile(0, w.Tiles[nearIdx].GetIconServer(), nearPos)
		}
	}
}

func (w *World) checkCascade(result *sweep_result.SweepResult, pos Coord) (didCascade bool) {
	if !pos.IsInRange(0, C.WORLD_MAX_X, 0, C.WORLD_MAX_Y) {
		return false
	}
	thisIdx := pos.ToIndex(C.TY_SHIFT)
	if w.Tiles[thisIdx].IsSwept() {
		return false
	}
	w.Tiles[thisIdx].SetVizSweptEmpty()
	if w.Tiles[thisIdx].GetNearby() == 0 {
		didCascade = true
	}
	result.AddTile(uint16(C.BOMB_NEAR_BASE_SCORE[0]), w.Tiles[thisIdx].GetIconServer(), pos)
	return didCascade
}

func (w *World) cascade(result *sweep_result.SweepResult, pos Coord) {
	queue := cascade_queue.New(pos)
	next, ok := queue.Next()
	for {
		if !ok {
			break
		}
		didCascade := w.checkCascade(result, next.Pos)
		if didCascade {
			queue.Cascade(next)
		}
		next, ok = queue.Next()
	}
}

// // This is NOT safe for concurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *World) DrawState(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(C.C.WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range C.WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range C.C.WORLD_TILE_WIDTH {
// 			iconCode := w.Tiles[i].GetIcon()
// 			switch iconCode {
// 			case tile.ICON_CODE_BOMB:
// 				wr.Write([]byte(ICON_SKULL))
// 			case tile.ICON_CODE_OPAQUE:
// 				wr.Write([]byte(ICON_OPAQUE))
// 			case tile.ICON_CODE_FLAG:
// 				wr.Write([]byte(ICON_FLAG))
// 			case tile.ICON_CODE_0:
// 				wr.Write([]byte(ICON_0))
// 			case tile.ICON_CODE_1:
// 				wr.Write([]byte(ICON_1))
// 			case tile.ICON_CODE_2:
// 				wr.Write([]byte(ICON_2))
// 			case tile.ICON_CODE_3:
// 				wr.Write([]byte(ICON_3))
// 			case tile.ICON_CODE_4:
// 				wr.Write([]byte(ICON_4))
// 			case tile.ICON_CODE_5:
// 				wr.Write([]byte(ICON_5))
// 			case tile.ICON_CODE_6:
// 				wr.Write([]byte(ICON_6))
// 			case tile.ICON_CODE_7:
// 				wr.Write([]byte(ICON_7))
// 			case tile.ICON_CODE_8:
// 				wr.Write([]byte(ICON_8))
// 			default:
// 				wr.Write([]byte(ICON_0))
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }

// // This is NOT safe for cuncurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *World) DrawMines(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(C.C.WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range C.WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range C.C.WORLD_TILE_WIDTH {
// 			isMine := w.Tiles[i].IsMine()
// 			if isMine {
// 				wr.Write([]byte(ICON_BOMB))
// 			} else {
// 				wr.Write([]byte(ICON_0))
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }

// // This is NOT safe for cuncurrent use when world is being played on,
// // this is only for testing/debugging purposes
// func (w *World) DrawNearby(wr io.Writer) {
// 	var i uint = 0
// 	var buf [4]byte
// 	capFill := strings.Repeat("═", int(C.C.WORLD_TILE_WIDTH))
// 	var cap string = "╔" + capFill + "╗\n"
// 	wr.Write([]byte(cap))
// 	for range C.WORLD_TILE_HEIGHT {
// 		n := utf8.EncodeRune(buf[:], '║')
// 		wr.Write(buf[:n])
// 		for range C.C.WORLD_TILE_WIDTH {
// 			isMine := w.Tiles[i].IsMine()
// 			if isMine {
// 				wr.Write([]byte(ICON_BOMB))
// 			} else {
// 				nearby := w.Tiles[i].GetNearby()
// 				switch nearby {
// 				case 0:
// 					wr.Write([]byte(ICON_0))
// 				case 1:
// 					wr.Write([]byte(ICON_1))
// 				case 2:
// 					wr.Write([]byte(ICON_2))
// 				case 3:
// 					wr.Write([]byte(ICON_3))
// 				case 4:
// 					wr.Write([]byte(ICON_4))
// 				case 5:
// 					wr.Write([]byte(ICON_5))
// 				case 6:
// 					wr.Write([]byte(ICON_6))
// 				case 7:
// 					wr.Write([]byte(ICON_7))
// 				case 8:
// 					wr.Write([]byte(ICON_8))
// 				default:
// 					wr.Write([]byte(ICON_0))
// 				}
// 			}
// 			i++
// 		}
// 		n = utf8.EncodeRune(buf[:], '║')
// 		buf[n] = 0x0A
// 		wr.Write(buf[:n+1])
// 	}
// 	cap = "╚" + capFill + "╝\n"
// 	wr.Write([]byte(cap))
// }

func (w *World) PrintStatus(wr io.Writer) {
	exploed := w.ExplodedMines.Load()
	swept := w.SweptTiles.Load()
	fmt.Fprintf(wr, "World ID: %d\n-Total Mines: %d\n-Exploded Mines: %d\n-Total Tiles: %d\n-SweptTiles: %d\n-Mine Completion: %.2f%%\n-Tile Completion: %.2f%%\n-Time Remaining: %s\n", w.Id.Load(), w.TotalMines, exploed, C.WORLD_TILE_COUNT, swept, float32(exploed)*100.0/float32(w.TotalMines), float32(swept)*100.0/float32(C.WORLD_TILE_COUNT), time.Until(w.Expires))
}

func (w *World) drawSweep(pos Coord) {
	sIdx := pos.ToIndex(C.TY_SHIFT)
	fmt.Print("\n\n")
	lastY := 0
	for idx := range C.WORLD_TILE_COUNT {
		coord := coord.FromIndex(idx, C.TY_SHIFT, C.TX_MASK)
		char := "1"
		if w.Tiles[idx].IsSwept() {
			char = "0"
		}
		if sIdx == idx {
			char = "X"
		}
		if lastY != coord.Y {
			fmt.Print("\n")
		}
		fmt.Print(char)
		lastY = coord.Y
	}
}
