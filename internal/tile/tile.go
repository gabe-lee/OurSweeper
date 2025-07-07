package tile

type Tile uint8

const (
	NEARBY_MASK  uint8 = 15
	NEARBY_CLEAR uint8 = 240
	MINE_MASK    uint8 = 16
	VIZ_MASK     uint8 = 96
	VIZ_CLEAR    uint8 = 159
	VIZ_OPAQUE   uint8 = 0
	VIZ_FLAG     uint8 = 32
	VIZ_SWEPT    uint8 = 64
	VIZ_BOMB     uint8 = 96
	LOCK_MASK    uint8 = 128
	LOCK_CLEAR   uint8 = 127

	ICON_0 uint8 = 0
	ICON_1 uint8 = 1
	ICON_2 uint8 = 2
	ICON_3 uint8 = 3
	ICON_4 uint8 = 4
	ICON_5 uint8 = 5
	ICON_6 uint8 = 6
	ICON_7 uint8 = 7
	ICON_8 uint8 = 8
	// 9-15 reserved

	ICON_OPAQUE uint8 = 16
	ICON_FLAG   uint8 = 17
	ICON_BOMB   uint8 = 18
)

func (t Tile) GetNearby() uint8 {
	return uint8(t) & NEARBY_MASK
}
func (t *Tile) SetNearby(near uint8) {
	*t = Tile(uint8(*t) & NEARBY_CLEAR)
	*t = Tile(uint8(*t) | near)
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
func (t *Tile) SetViz(viz uint8) {
	*t = Tile(uint8(*t) & VIZ_CLEAR)
	*t = Tile(uint8(*t) | viz)
}

func (t Tile) IsLocked() bool {
	return uint8(t)&LOCK_MASK == LOCK_MASK
}
func (t *Tile) Lock() {
	*t = Tile(uint8(*t) | LOCK_MASK)
}
func (t *Tile) Unlock() {
	*t = Tile(uint8(*t) & LOCK_CLEAR)
}

func (t Tile) GetIcon() uint8 {
	viz := t.GetViz()
	switch viz {
	case VIZ_OPAQUE:
		return ICON_OPAQUE
	case VIZ_FLAG:
		return ICON_FLAG
	case VIZ_BOMB:
		return ICON_BOMB
	default:
		return t.GetNearby()
	}
}
