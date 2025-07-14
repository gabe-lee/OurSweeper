package internal

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gabe-lee/OurSweeper/serializer"
	"github.com/gabe-lee/OurSweeper/utils"
)

type Coord struct {
	X int
	Y int
}

func (c Coord) Equals(other Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Coord) Sub(other Coord) Coord {
	return Coord{
		X: c.X - other.X,
		Y: c.Y - other.Y,
	}
}

func (c Coord) Mult(other Coord) Coord {
	return Coord{
		X: c.X * other.X,
		Y: c.Y * other.Y,
	}
}

func (c Coord) Div(other Coord) Coord {
	return Coord{
		X: c.X / other.X,
		Y: c.Y / other.Y,
	}
}

func (c Coord) MultScalar(val int) Coord {
	return Coord{
		X: c.X * val,
		Y: c.Y * val,
	}
}

func (c Coord) DivScalar(val int) Coord {
	return Coord{
		X: c.X / val,
		Y: c.Y / val,
	}
}

func (c Coord) ShiftDownScalar(val int) Coord {
	return Coord{
		X: c.X >> val,
		Y: c.Y >> val,
	}
}

func (c Coord) ShiftUpScalar(val int) Coord {
	return Coord{
		X: c.X << val,
		Y: c.Y << val,
	}
}

func (c Coord) Invert() Coord {
	return Coord{
		X: -c.X,
		Y: -c.Y,
	}
}
func (c Coord) InvertX() Coord {
	return Coord{
		X: -c.X,
		Y: c.Y,
	}
}
func (c Coord) InvertY() Coord {
	return Coord{
		X: c.X,
		Y: -c.Y,
	}
}

func (c Coord) Clamp(minX, maxX, minY, maxY int) Coord {
	return Coord{
		X: min(maxX, max(minX, c.X)),
		Y: min(maxY, max(minY, c.Y)),
	}
}
func (c Coord) ClampMin(minX, minY int) Coord {
	return Coord{
		X: max(minX, c.X),
		Y: max(minY, c.Y),
	}
}
func (c Coord) ClampMax(maxX, maxY int) Coord {
	return Coord{
		X: min(maxX, c.X),
		Y: min(maxY, c.Y),
	}
}

func (c Coord) IsInRange(minX, maxX, minY, maxY int) bool {
	return c.X >= minX && c.X <= maxX && c.Y >= minY && c.Y <= maxY
}

func (c Coord) GetBounds2(growBy Coord, minX, maxX, minY, maxY int) (bounds Bounds2) {
	bounds.BotRight = c.Add(growBy).ClampMax(maxX, maxY)
	bounds.TopLeft = c.Add(growBy.Invert()).ClampMin(minX, minY)
	return
}

func (c Coord) GetBounds4(growBy Coord, minX, maxX, minY, maxY int) (bounds Bounds4) {
	bounds.BotRight = c.Add(growBy).ClampMax(maxX, maxY)
	bounds.TopLeft = c.Add(growBy.Invert()).ClampMin(minX, minY)
	bounds.TopRight = Coord{
		X: bounds.BotRight.X,
		Y: bounds.TopLeft.Y,
	}
	bounds.BotLeft = Coord{
		X: bounds.TopLeft.X,
		Y: bounds.BotRight.Y,
	}
	return
}

func (c Coord) IsInRangeExcludeEdges(minX, maxX, minY, maxY int) bool {
	return c.X > minX && c.X < maxX && c.Y > minY && c.Y < maxY
}

func (c Coord) ToIndex(yShift int) int {
	return c.Y<<yShift | c.X
}

func CoordFromIndex(index int, yShift int, xMask int) Coord {
	return Coord{
		X: index & xMask,
		Y: index >> yShift,
	}
}

func (c Coord) ToByteCoord() ByteCoord {
	return ByteCoord{
		X: byte(c.X),
		Y: byte(c.Y),
	}
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

type ByteCoord struct {
	X byte
	Y byte
}

// ReadWire implements wire_serializer.WireSerializer.
func (b *ByteCoord) Deserialize(r io.Reader, order binary.ByteOrder) error {
	var e utils.ErrorCollector
	e.Do(binary.Read(r, order, &b.X))
	e.Do(binary.Read(r, order, &b.Y))
	return e.Err
}

// WriteWire implements wire_serializer.WireSerializer.
func (b *ByteCoord) Serialize(w io.Writer, order binary.ByteOrder) error {
	var e utils.ErrorCollector
	e.Do(binary.Write(w, order, b.X))
	e.Do(binary.Write(w, order, b.Y))
	return e.Err
}

func (b ByteCoord) ToCoord() Coord {
	return Coord{
		X: int(b.X),
		Y: int(b.Y),
	}
}

var _ serializer.Serializer = (*ByteCoord)(nil)
var _ serializer.Deserializer = (*ByteCoord)(nil)

type NearbyCoords struct {
	Coords [8]Coord
	Bits   [8]uint64
	Len    int
}

func (c Coord) GetNearbyCoords(minX, maxX, minY, maxY int) NearbyCoords {
	near := NearbyCoords{}
	for i, offset := range NearCoordTable8 {
		pos := c.Add(offset)
		if pos.IsInRange(minX, maxX, minY, maxY) {
			near.Coords[near.Len] = pos
			near.Bits[near.Len] = NearBitTable8[i]
			near.Len += 1
		}
	}
	return near
}

type Bounds2 struct {
	TopLeft  Coord
	BotRight Coord
}

type Bounds4 struct {
	TopLeft  Coord
	TopRight Coord
	BotLeft  Coord
	BotRight Coord
}

func (b Bounds4) DivScalar(scale int) Bounds4 {
	return Bounds4{
		TopLeft:  b.TopLeft.DivScalar(scale),
		TopRight: b.TopRight.DivScalar(scale),
		BotLeft:  b.BotLeft.DivScalar(scale),
		BotRight: b.BotRight.DivScalar(scale),
	}
}

func (b Bounds4) MultScalar(scale int) Bounds4 {
	return Bounds4{
		TopLeft:  b.TopLeft.MultScalar(scale),
		TopRight: b.TopRight.MultScalar(scale),
		BotLeft:  b.BotLeft.MultScalar(scale),
		BotRight: b.BotRight.MultScalar(scale),
	}
}

func (b Bounds4) ShiftDownScalar(shift int) Bounds4 {
	return Bounds4{
		TopLeft:  b.TopLeft.ShiftDownScalar(shift),
		TopRight: b.TopRight.ShiftDownScalar(shift),
		BotLeft:  b.BotLeft.ShiftDownScalar(shift),
		BotRight: b.BotRight.ShiftDownScalar(shift),
	}
}

func (b Bounds4) ShiftUpScalar(shift int) Bounds4 {
	return Bounds4{
		TopLeft:  b.TopLeft.ShiftUpScalar(shift),
		TopRight: b.TopRight.ShiftUpScalar(shift),
		BotLeft:  b.BotLeft.ShiftUpScalar(shift),
		BotRight: b.BotRight.ShiftUpScalar(shift),
	}
}
