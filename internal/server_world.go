package internal

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gabe-lee/OurSweeper/xmath"
)

type (
	Mutex = sync.Mutex
)

type ServerWorld struct {
	Id            atomic.Uint32
	Tiles         [WORLD_TILE_COUNT]Tile
	Locks         [WORLD_LOCK_COUNT]Mutex
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
	if pos.X > WORLD_HALF_WIDTH {
		dx = fx - WORLD_CENTER_X
	} else {
		dx = WORLD_CENTER_X - fx
	}
	if pos.Y > WORLD_HALF_HEIGHT {
		dy = fy - WORLD_CENTER_Y
	} else {
		dy = WORLD_CENTER_Y - fy
	}
	dist := math.Sqrt((dx * dx) + (dy * dy))
	percent := (MAX_DIST_FROM_CENTER - dist) / MAX_DIST_FROM_CENTER
	return xmath.Lerp(MIN_MINE_CHANCE, MAX_MINE_CHANCE, percent)
}

func (w *ServerWorld) LockEntireWorld() {
	for i := range WORLD_LOCK_COUNT {
		w.Locks[i].Lock()
	}
}
func (w *ServerWorld) UnlockEntireWorld() {
	for i := range WORLD_LOCK_COUNT {
		w.Locks[i].Unlock()
	}
}

func (w *ServerWorld) AquireQuadLock(bounds Bounds4) QuadLock {
	lockBounds := bounds.ShiftDownScalar(T_TO_L)
	topLeftIdx := lockBounds.TopLeft.ToIndex(LY_SHIFT)
	topRightIdx := lockBounds.TopRight.ToIndex(LY_SHIFT)
	botLeftIdx := lockBounds.BotLeft.ToIndex(LY_SHIFT)
	botRightIdx := lockBounds.BotRight.ToIndex(LY_SHIFT)
	locks := QuadLock{}
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

func (w *ServerWorld) ReleaseQuadLock(locks QuadLock) {
	for _, idx := range locks.LockIndexes[:locks.Len] {
		w.Locks[idx].Unlock()
	}
}

// func (w *World) TryLockTile(set *WORLD_LOCK_COUNTet.WORLD_LOCK_COUNTet, x, y int) bool {
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

// func (w *World) LockTile(set *WORLD_LOCK_COUNTet.WORLD_LOCK_COUNTet, lockIdx, lx, ly int) {
// 	if set.AlreadyLocked(lockIdx) {
// 		return
// 	}
// 	w.Locks[lockIdx].Lock()
// 	set.AddLock(lockIdx, lx, ly)
// }

// func (w *World) UnlockTiles(set *WORLD_LOCK_COUNTet.WORLD_LOCK_COUNTet) {
// 	for y := set.YMin; y <= set.YMax; y++ {
// 		for x := set.XMin; x <= set.XMax; x++ {
// 			idx := (uint64(y) * 8) + uint64(x)
// 			bit := uint64(1) << idx
// 			if set.WORLD_LOCK_COUNT&bit > 0 {
// 				w.Locks[idx].Unlock()
// 			}
// 		}
// 	}
// 	*set = WORLD_LOCK_COUNTet.WORLD_LOCK_COUNTet{}
// }

func (w *ServerWorld) initMine(idx int, pos Coord) bool {
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
// 	yMax := min(y+1, int(WORLD_TILE_WIDTH-1))
// 	xMin := max(y-1, 0)
// 	xMax := min(y+1, int(WORLD_TILE_HEIGHT-1))
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

func (w *ServerWorld) InitNew(id uint32) {
	w.Id.Store(id)
	w.SweptTiles.Store(uint32(INITIAL_SWEPT_TILES))
	w.ExplodedMines.Store(0)
	w.Ended.Store(false)
	w.Expires = time.Now().Add(WORLD_TIME_LIMIT)
	for idx := range WORLD_TILE_COUNT {
		w.Tiles[idx] = Tile(0)
	}
	for idx := range WORLD_LOCK_COUNT {
		w.Locks[idx] = sync.Mutex{}
	}

	for idx := range WORLD_TILE_COUNT {
		pos := CoordFromIndex(idx, TY_SHIFT, TX_MASK)
		if pos.IsInRangeExcludeEdges(0, WORLD_MAX_X, 0, WORLD_MAX_Y) {
			mine := w.initMine(idx, pos)
			if mine {
				w.TotalMines += 1
				nears := pos.GetNearbyCoords(0, WORLD_MAX_X, 0, WORLD_MAX_Y)
				for _, nearPos := range nears.Coords[:nears.Len] {
					nearIdx := nearPos.ToIndex(TY_SHIFT)
					w.Tiles[nearIdx].IncrNearbyMineCount()
				}
			}
		} else {
			w.Tiles[idx].SetVizSweptEmpty()
		}
	}
}

func (w *ServerWorld) SweepTile(pos Coord) SweepResult {
	result := SweepResult{}
	tileIdx := pos.ToIndex(TY_SHIFT)
	t := w.Tiles[tileIdx]
	if t.IsSwept() {
		return result
	}
	bounds := pos.GetBounds4(Coord{X: MAX_CASCDE_DIST, Y: MAX_CASCDE_DIST}, 0, WORLD_MAX_X, 0, WORLD_MAX_Y)
	locks := w.AquireQuadLock(bounds)
	defer w.ReleaseQuadLock(locks)
	isMine := t.IsMine()
	if isMine {
		t.SetVizSweptBomb()
		w.ExplodedMines.Add(1)
	} else {
		t.SetVizSweptEmpty()
	}
	result.InitSweep(pos, w.getScore(pos), t.GetIconServer())
	if isMine {
		w.reduceNearbyBombCounts(&result, pos)
	} else if t.GetNearby() == 0 {
		casc := w.cascade(&result, pos)
		result.RelativeBits = casc.DidSweepList
	}
	if !isMine {
		w.SweptTiles.Add(uint32(result.Len))
	}
	w.Tiles[tileIdx] = t
	w.checkEndState()

	return result
}

func (w *ServerWorld) checkEndState() {
	if exploded := w.ExplodedMines.Load(); exploded == w.TotalMines {
		w.Ended.Store(true)
	} else if swept := w.SweptTiles.Load(); swept == uint32(WORLD_TILE_COUNT) {
		w.Ended.Store(true)
	} else if time.Now().After(w.Expires) {
		w.Ended.Store(true)
	}
}

func (w *ServerWorld) getScore(pos Coord) uint16 {
	var lowestNearBombChance float64 = RISC_FULL_BLIND
	var lowestNearBombChanceBombs uint8 = IDX_FULL_BLIND_BASE
	nearbyCoords := pos.GetNearbyCoords(0, WORLD_MAX_X, 0, WORLD_MAX_Y)
	for _, nearPos := range nearbyCoords.Coords[:nearbyCoords.Len] {
		nearIdx := nearPos.ToIndex(TY_SHIFT)
		if w.Tiles[nearIdx].IsSwept() {
			thisBombChance, thisBombs := w.getBombProbabilityAndNearby(nearPos)
			if thisBombChance < lowestNearBombChance {
				lowestNearBombChance = thisBombChance
				lowestNearBombChanceBombs = thisBombs
			}
		}
	}
	scoreFloat := BOMB_NEAR_BASE_SCORE[lowestNearBombChanceBombs]
	if lowestNearBombChance < 1.0 {
		exp := RISC_EXPONENT_ADD + (RISC_EXPONENT_MULT * lowestNearBombChance)
		scoreFloat = math.Ceil(math.Pow(scoreFloat, exp))
	}
	return uint16(scoreFloat)
}

func (w *ServerWorld) getBombProbabilityAndNearby(pos Coord) (safe float64, near uint8) {
	var opaques float64
	bombs := w.Tiles[pos.ToIndex(TY_SHIFT)].GetNearby()
	nearbyCoords := pos.GetNearbyCoords(0, WORLD_MAX_X, 0, WORLD_MAX_Y)
	for _, nearPos := range nearbyCoords.Coords[:nearbyCoords.Len] {
		nearIdx := nearPos.ToIndex(TY_SHIFT)
		if !w.Tiles[nearIdx].IsSwept() {
			opaques += 1.0
		}
	}
	return float64(bombs) / opaques, bombs
}

func (w *ServerWorld) reduceNearbyBombCounts(result *SweepResult, pos Coord) {
	nearbyCoords := pos.GetNearbyCoords(0, WORLD_MAX_X, 0, WORLD_MAX_Y)
	for i := range nearbyCoords.Len {
		nearPos := nearbyCoords.Coords[i]
		nearBit := nearbyCoords.Bits[i]
		nearIdx := nearPos.ToIndex(TY_SHIFT)
		bombs := w.Tiles[nearIdx].GetNearby()
		w.Tiles[nearIdx].SetNearby(bombs - 1)
		if w.Tiles[nearIdx].IsSwept() {
			result.RelativeBits |= nearBit
			result.AddScoreAndIcon(0, w.Tiles[nearIdx].GetIconServer())
		}
	}
}

func (w *ServerWorld) checkCascade(result *SweepResult, pos Coord) (didSweep, didCascade bool) {
	if !pos.IsInRange(0, WORLD_MAX_X, 0, WORLD_MAX_Y) {
		return false, false
	}
	thisIdx := pos.ToIndex(TY_SHIFT)
	if w.Tiles[thisIdx].IsSwept() {
		return false, false
	}
	w.Tiles[thisIdx].SetVizSweptEmpty()
	if w.Tiles[thisIdx].GetNearby() == 0 {
		didCascade = true
	}
	result.AddScoreAndIcon(uint16(BOMB_NEAR_BASE_SCORE[0]), w.Tiles[thisIdx].GetIconServer())
	return true, didCascade
}

func (w *ServerWorld) cascade(result *SweepResult, pos Coord) CascadeQueue {
	queue := NewCascadeQueue(pos)
	next, moreToCheck := queue.NextToCheck()
	for moreToCheck {
		didSweep, didCascade := w.checkCascade(result, next.Pos)
		if didSweep {
			queue.AddSweep(next)
		}
		if didCascade {
			queue.AddCascade(next)
		}
		next, moreToCheck = queue.NextToCheck()
	}
	return queue
}

func (w *ServerWorld) PrintStatus(wr io.Writer) {
	exploed := w.ExplodedMines.Load()
	swept := w.SweptTiles.Load()
	fmt.Fprintf(wr, "World ID: %d\n-Total Mines: %d\n-Exploded Mines: %d\n-Total Tiles: %d\n-SweptTiles: %d\n-Mine Completion: %.2f%%\n-Tile Completion: %.2f%%\n-Time Remaining: %s\n", w.Id.Load(), w.TotalMines, exploed, WORLD_TILE_COUNT, swept, float32(exploed)*100.0/float32(w.TotalMines), float32(swept)*100.0/float32(WORLD_TILE_COUNT), time.Until(w.Expires))
}
