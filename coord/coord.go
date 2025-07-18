package coord

import (
	"fmt"

	"github.com/gabe-lee/OurSweeper/wire"
)

type (
	IncomingWire = wire.IncomingWire
	OutgoingWire = wire.OutgoingWire
)

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Coord[T number] struct {
	X T
	Y T
}

func NewCoord[T number](x, y T) Coord[T] {
	return Coord[T]{
		X: x,
		Y: y,
	}
}

func (c Coord[T]) Equals(other Coord[T]) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c Coord[T]) Add(other Coord[T]) Coord[T] {
	return Coord[T]{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Coord[T]) AddXY(x, y T) Coord[T] {
	return Coord[T]{
		X: c.X + x,
		Y: c.Y + y,
	}
}

func (c Coord[T]) AddScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X + val,
		Y: c.Y + val,
	}
}

func (c Coord[T]) Sub(other Coord[T]) Coord[T] {
	return Coord[T]{
		X: c.X - other.X,
		Y: c.Y - other.Y,
	}
}

func (c Coord[T]) SubXY(x, y T) Coord[T] {
	return Coord[T]{
		X: c.X - x,
		Y: c.Y - y,
	}
}

func (c Coord[T]) SubScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X - val,
		Y: c.Y - val,
	}
}

func (c Coord[T]) Mult(other Coord[T]) Coord[T] {
	return Coord[T]{
		X: c.X * other.X,
		Y: c.Y * other.Y,
	}
}

func (c Coord[T]) MultXY(x, y T) Coord[T] {
	return Coord[T]{
		X: c.X * x,
		Y: c.Y * y,
	}
}

func (c Coord[T]) MultScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X * val,
		Y: c.Y * val,
	}
}

func (c Coord[T]) Div(other Coord[T]) Coord[T] {
	return Coord[T]{
		X: c.X / other.X,
		Y: c.Y / other.Y,
	}
}

func (c Coord[T]) DivXY(x, y T) Coord[T] {
	return Coord[T]{
		X: c.X / x,
		Y: c.Y / y,
	}
}

func (c Coord[T]) DivScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X / val,
		Y: c.Y / val,
	}
}

func (c Coord[T]) ShiftDownScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X >> val,
		Y: c.Y >> val,
	}
}

func (c Coord[T]) ShiftUpScalar(val T) Coord[T] {
	return Coord[T]{
		X: c.X << val,
		Y: c.Y << val,
	}
}

func (c Coord[T]) Invert() Coord[T] {
	return Coord[T]{
		X: -c.X,
		Y: -c.Y,
	}
}
func (c Coord[T]) InvertX() Coord[T] {
	return Coord[T]{
		X: -c.X,
		Y: c.Y,
	}
}
func (c Coord[T]) InvertY() Coord[T] {
	return Coord[T]{
		X: c.X,
		Y: -c.Y,
	}
}

func (c Coord[T]) Clamp(minX, maxX, minY, maxY T) Coord[T] {
	return Coord[T]{
		X: min(maxX, max(minX, c.X)),
		Y: min(maxY, max(minY, c.Y)),
	}
}
func (c Coord[T]) ClampMin(minX, minY T) Coord[T] {
	return Coord[T]{
		X: max(minX, c.X),
		Y: max(minY, c.Y),
	}
}
func (c Coord[T]) ClampMax(maxX, maxY T) Coord[T] {
	return Coord[T]{
		X: min(maxX, c.X),
		Y: min(maxY, c.Y),
	}
}

func (c Coord[T]) IsInRange(minX, maxX, minY, maxY T) bool {
	return c.X >= minX && c.X <= maxX && c.Y >= minY && c.Y <= maxY
}

func (c Coord[T]) GetBounds2(growBy Coord[T], minX, maxX, minY, maxY T) (bounds Bounds2[T]) {
	bounds.BotRight = c.Add(growBy).ClampMax(maxX, maxY)
	bounds.TopLeft = c.Add(growBy.Invert()).ClampMin(minX, minY)
	return
}

func (c Coord[T]) GetBounds4(growBy Coord[T], minX, maxX, minY, maxY T) (bounds Bounds4[T]) {
	bounds.BotRight = c.Add(growBy).ClampMax(maxX, maxY)
	bounds.TopLeft = c.Add(growBy.Invert()).ClampMin(minX, minY)
	bounds.TopRight = Coord[T]{
		X: bounds.BotRight.X,
		Y: bounds.TopLeft.Y,
	}
	bounds.BotLeft = Coord[T]{
		X: bounds.TopLeft.X,
		Y: bounds.BotRight.Y,
	}
	return
}

func (c Coord[T]) IsInRangeExcludeEdges(minX, maxX, minY, maxY T) bool {
	return c.X > minX && c.X < maxX && c.Y > minY && c.Y < maxY
}

func (c Coord[T]) ToIndex(yShift T) T {
	return c.Y<<yShift | c.X
}

func CoordFromIndex[T number](index T, yShift T, xMask T) Coord[T] {
	return Coord[T]{
		X: index & xMask,
		Y: index >> yShift,
	}
}

func (c Coord[T]) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

func (c *Coord[T]) WireRead(w *IncomingWire) {
	w.TryRead_Auto(&c.X)
	w.TryRead_Auto(&c.Y)
}

func (c *Coord[T]) WireWrite(w *wire.OutgoingWire) {
	w.TryWrite_Auto(c.X)
	w.TryWrite_Auto(c.Y)
}

func (c Coord[T]) ToCoordByte() Coord[byte] {
	return Coord[byte]{
		X: byte(c.X),
		Y: byte(c.Y),
	}
}
func (c Coord[T]) ToCoordInt() Coord[int] {
	return Coord[int]{
		X: int(c.X),
		Y: int(c.Y),
	}
}

var _ wire.WireWriter = (*Coord[byte])(nil)
var _ wire.WireReader = (*Coord[byte])(nil)
var _ wire.WireWriter = (*Coord[int])(nil)
var _ wire.WireReader = (*Coord[int])(nil)

// func (c Coord) GetNearbyCoords(minX, maxX, minY, maxY int) NearbyCoords {
// 	near := NearbyCoords{}
// 	for i, offset := range NearCoordTable8 {
// 		pos := c.Add(offset)
// 		if pos.IsInRange(minX, maxX, minY, maxY) {
// 			near.Coords[near.Len] = pos
// 			near.Bits[near.Len] = NearBitTable8[i]
// 			near.Len += 1
// 		}
// 	}
// 	return near
// }

type Bounds2[T number] struct {
	TopLeft  Coord[T]
	BotRight Coord[T]
}

type Bounds4[T number] struct {
	TopLeft  Coord[T]
	TopRight Coord[T]
	BotLeft  Coord[T]
	BotRight Coord[T]
}

func (b Bounds4[T]) DivScalar(scale T) Bounds4[T] {
	return Bounds4[T]{
		TopLeft:  b.TopLeft.DivScalar(scale),
		TopRight: b.TopRight.DivScalar(scale),
		BotLeft:  b.BotLeft.DivScalar(scale),
		BotRight: b.BotRight.DivScalar(scale),
	}
}

func (b Bounds4[T]) MultScalar(scale T) Bounds4[T] {
	return Bounds4[T]{
		TopLeft:  b.TopLeft.MultScalar(scale),
		TopRight: b.TopRight.MultScalar(scale),
		BotLeft:  b.BotLeft.MultScalar(scale),
		BotRight: b.BotRight.MultScalar(scale),
	}
}

func (b Bounds4[T]) ShiftDownScalar(shift T) Bounds4[T] {
	return Bounds4[T]{
		TopLeft:  b.TopLeft.ShiftDownScalar(shift),
		TopRight: b.TopRight.ShiftDownScalar(shift),
		BotLeft:  b.BotLeft.ShiftDownScalar(shift),
		BotRight: b.BotRight.ShiftDownScalar(shift),
	}
}

func (b Bounds4[T]) ShiftUpScalar(shift T) Bounds4[T] {
	return Bounds4[T]{
		TopLeft:  b.TopLeft.ShiftUpScalar(shift),
		TopRight: b.TopRight.ShiftUpScalar(shift),
		BotLeft:  b.BotLeft.ShiftUpScalar(shift),
		BotRight: b.BotRight.ShiftUpScalar(shift),
	}
}
