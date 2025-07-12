package tile

import (
	C "github.com/gabe-lee/OurSweeper/internal/common"
)

type Tile uint8

const (
	NEARBY_MASK  uint8 = 0b0_00_0_1111
	NEARBY_CLEAR uint8 = 0b1_11_1_0000
	MINE_MASK    uint8 = 0b0_00_1_0000
	VIZ_MASK     uint8 = 0b0_11_0_0000
	VIZ_CLEAR    uint8 = 0b1_00_1_1111
	VIZ_OPAQUE   uint8 = 0b0_00_0_0000
	VIZ_FLAG     uint8 = 0b0_01_0_0000
	VIZ_EMPTY    uint8 = 0b0_10_0_0000
	VIZ_BOMB     uint8 = 0b0_11_0_0000
)

func (t Tile) GetNearby() uint8 {
	return uint8(t) & NEARBY_MASK
}
func (t *Tile) SetNearby(near uint8) {
	*t = Tile(uint8(*t) & NEARBY_CLEAR)
	*t = Tile(uint8(*t) | near)
}
func (t *Tile) IncrNearbyMineCount() {
	*t = Tile(uint8(*t) + 1)
}
func (t *Tile) DecrNearbyMineCount() {
	*t = Tile(uint8(*t) - 1)
}

func (t Tile) IsMine() bool {
	return uint8(t)&MINE_MASK == MINE_MASK
}
func (t *Tile) SetMine() {
	*t = Tile(uint8(*t) | MINE_MASK)
}

func (t Tile) GetViz() uint8 {
	return uint8(t) & VIZ_MASK
}
func (t Tile) IsSwept() bool {
	// fmt.Printf("TILE:  %08b\nSWEPT: %08b\nAND:   %08b\nRESULT: %v\n", uint8(t), VIZ_EMPTY, uint8(t)&VIZ_EMPTY, uint8(t)&VIZ_EMPTY == VIZ_EMPTY) //DEBUG
	return uint8(t)&VIZ_EMPTY == VIZ_EMPTY
}
func (t *Tile) SetVizOpaque() {
	*t = Tile(uint8(*t) & VIZ_CLEAR)
}
func (t *Tile) SetVizFlag() {
	*t = Tile(uint8(*t) & VIZ_CLEAR)
	*t = Tile(uint8(*t) | VIZ_FLAG)
}
func (t *Tile) SetVizSweptEmpty() {
	*t = Tile(uint8(*t) & VIZ_CLEAR)
	*t = Tile(uint8(*t) | VIZ_EMPTY)
}
func (t *Tile) SetVizSweptBomb() {
	*t = Tile(uint8(*t) & VIZ_CLEAR)
	*t = Tile(uint8(*t) | VIZ_BOMB)
}

func (t Tile) GetIconServer() uint8 {
	viz := t.GetViz()
	switch viz {
	case VIZ_OPAQUE:
		return C.ICON_CODE_OPAQUE
	case VIZ_FLAG:
		return C.ICON_CODE_FLAG
	case VIZ_BOMB:
		return C.ICON_CODE_BOMB
	default:
		return t.GetNearby()
	}
}

func (t Tile) GetIconClient() uint8 {
	return t.GetNearby()
}
