package color

import "math"

type colorparam interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

const (
	INVISIBLE = iota
	WHITE
	RED
	GREEN
	BLUE
	CYAN
	MAGENTA
	YELLOW
	BLACK
	_colorCount
)

const (
	max8   uint8   = math.MaxUint8
	max16  uint16  = math.MaxUint16
	max32  uint32  = math.MaxUint32
	max64  uint64  = math.MaxUint64
	maxf32 float32 = 1.0
	maxf64 float64 = 1.0
)

type Color[T colorparam] struct {
	Red   T
	Green T
	Blue  T
	Alpha T
}

func New[T colorparam](r, g, b, a T) Color[T] {
	return Color[T]{
		Red:   r,
		Green: g,
		Blue:  b,
		Alpha: a,
	}
}

var COLOR_U8 = [_colorCount]Color[uint8]{
	INVISIBLE: New[uint8](0, 0, 0, 0),
	WHITE:     New(max8, max8, max8, max8),
	RED:       New(max8, 0, 0, max8),
	GREEN:     New(0, max8, 0, max8),
	BLUE:      New(0, 0, max8, max8),
	CYAN:      New(0, max8, max8, max8),
	YELLOW:    New(max8, max8, 0, max8),
	MAGENTA:   New(max8, 0, max8, max8),
	BLACK:     New(0, 0, 0, max8),
}

var COLOR_U16 = [_colorCount]Color[uint16]{
	INVISIBLE: New[uint16](0, 0, 0, 0),
	WHITE:     New(max16, max16, max16, max16),
	RED:       New(max16, 0, 0, max16),
	GREEN:     New(0, max16, 0, max16),
	BLUE:      New(0, 0, max16, max16),
	CYAN:      New(0, max16, max16, max16),
	YELLOW:    New(max16, max16, 0, max16),
	MAGENTA:   New(max16, 0, max16, max16),
	BLACK:     New(0, 0, 0, max16),
}

var COLOR_U32 = [_colorCount]Color[uint32]{
	INVISIBLE: New[uint32](0, 0, 0, 0),
	WHITE:     New(max32, max32, max32, max32),
	RED:       New(max32, 0, 0, max32),
	GREEN:     New(0, max32, 0, max32),
	BLUE:      New(0, 0, max32, max32),
	CYAN:      New(0, max32, max32, max32),
	YELLOW:    New(max32, max32, 0, max32),
	MAGENTA:   New(max32, 0, max32, max32),
	BLACK:     New(0, 0, 0, max32),
}

var COLOR_U64 = [_colorCount]Color[uint64]{
	INVISIBLE: New[uint64](0, 0, 0, 0),
	WHITE:     New(max64, max64, max64, max64),
	RED:       New(max64, 0, 0, max64),
	GREEN:     New(0, max64, 0, max64),
	BLUE:      New(0, 0, max64, max64),
	CYAN:      New(0, max64, max64, max64),
	YELLOW:    New(max64, max64, 0, max64),
	MAGENTA:   New(max64, 0, max64, max64),
	BLACK:     New(0, 0, 0, max64),
}

var COLOR_F32 = [_colorCount]Color[float32]{
	INVISIBLE: New[float32](0, 0, 0, 0),
	WHITE:     New(maxf32, maxf32, maxf32, maxf32),
	RED:       New(maxf32, 0, 0, maxf32),
	GREEN:     New(0, maxf32, 0, maxf32),
	BLUE:      New(0, 0, maxf32, maxf32),
	CYAN:      New(0, maxf32, maxf32, maxf32),
	YELLOW:    New(maxf32, maxf32, 0, maxf32),
	MAGENTA:   New(maxf32, 0, maxf32, maxf32),
	BLACK:     New(0, 0, 0, maxf32),
}

var COLOR_F64 = [_colorCount]Color[float64]{
	INVISIBLE: New[float64](0, 0, 0, 0),
	WHITE:     New(maxf64, maxf64, maxf64, maxf64),
	RED:       New(maxf64, 0, 0, maxf64),
	GREEN:     New(0, maxf64, 0, maxf64),
	BLUE:      New(0, 0, maxf64, maxf64),
	CYAN:      New(0, maxf64, maxf64, maxf64),
	YELLOW:    New(maxf64, maxf64, 0, maxf64),
	MAGENTA:   New(maxf64, 0, maxf64, maxf64),
	BLACK:     New(0, 0, 0, maxf64),
}
