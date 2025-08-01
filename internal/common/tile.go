package common

import (
	C "github.com/gabe-lee/OurSweeper/internal/consts"
)

type Tile uint8

func (t Tile) GetNearby() uint8 {
	return uint8(t) & C.NEARBY_MASK
}
func (t *Tile) SetNearby(near uint8) {
	*t = Tile(uint8(*t) & C.NEARBY_CLEAR)
	*t = Tile(uint8(*t) | near)
}
func (t *Tile) IncrNearbyMineCount() {
	*t = Tile(uint8(*t) + 1)
}
func (t *Tile) DecrNearbyMineCount() {
	*t = Tile(uint8(*t) - 1)
}

func (t Tile) IsMine() bool {
	return uint8(t)&C.MINE_MASK == C.MINE_MASK
}
func (t *Tile) SetMine() {
	*t = Tile(uint8(*t) | C.MINE_MASK)
}

func (t Tile) GetViz() uint8 {
	return uint8(t) & C.VIZ_MASK
}
func (t Tile) IsSwept() bool {
	// fmt.Printf("TILE:  %08b\nSWEPT: %08b\nAND:   %08b\nRESULT: %v\n", uint8(t), VIZ_EMPTY, uint8(t)&VIZ_EMPTY, uint8(t)&VIZ_EMPTY == VIZ_EMPTY) //DEBUG
	return uint8(t)&C.VIZ_EMPTY == C.VIZ_EMPTY
}
func (t *Tile) SetVizOpaque() {
	*t = Tile(uint8(*t) & C.VIZ_CLEAR)
}
func (t *Tile) SetVizFlag() {
	*t = Tile(uint8(*t) & C.VIZ_CLEAR)
	*t = Tile(uint8(*t) | C.VIZ_FLAG)
}
func (t *Tile) SetVizSweptEmpty() {
	*t = Tile(uint8(*t) & C.VIZ_CLEAR)
	*t = Tile(uint8(*t) | C.VIZ_EMPTY)
}
func (t *Tile) SetVizSweptBomb() {
	*t = Tile(uint8(*t) & C.VIZ_CLEAR)
	*t = Tile(uint8(*t) | C.VIZ_BOMB)
}

func (t Tile) GetIconForClient() uint8 {
	viz := t.GetViz()
	switch viz {
	case C.VIZ_OPAQUE:
		return C.ICON_CODE_OPAQUE
	case C.VIZ_FLAG:
		return C.ICON_CODE_FLAG
	case C.VIZ_BOMB:
		return C.ICON_CODE_BOMB
	default:
		return t.GetNearby()
	}
}

func (t Tile) GetIconRevealed() uint8 {
	if t.IsMine() {
		return C.ICON_CODE_BOMB
	}
	return t.GetNearby()
}
