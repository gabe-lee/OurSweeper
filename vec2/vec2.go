package vec2

import "fmt"

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Vec2[T number] struct {
	X T
	Y T
}

func New[T number](x, y T) Vec2[T] {
	return Vec2[T]{
		X: x,
		Y: y,
	}
}

func (c Vec2[T]) Equals(other Vec2[T]) bool {
	return c.X == other.X && c.Y == other.Y
}

func (c Vec2[T]) Add(other Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: c.X + other.X,
		Y: c.Y + other.Y,
	}
}

func (c Vec2[T]) AddXY(x, y T) Vec2[T] {
	return Vec2[T]{
		X: c.X + x,
		Y: c.Y + y,
	}
}

func (c Vec2[T]) AddScalar(val T) Vec2[T] {
	return Vec2[T]{
		X: c.X + val,
		Y: c.Y + val,
	}
}

func (c Vec2[T]) Sub(other Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: c.X - other.X,
		Y: c.Y - other.Y,
	}
}

func (c Vec2[T]) SubXY(x, y T) Vec2[T] {
	return Vec2[T]{
		X: c.X - x,
		Y: c.Y - y,
	}
}

func (c Vec2[T]) SubScalar(val T) Vec2[T] {
	return Vec2[T]{
		X: c.X - val,
		Y: c.Y - val,
	}
}

func (c Vec2[T]) Mult(other Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: c.X * other.X,
		Y: c.Y * other.Y,
	}
}

func (c Vec2[T]) MultXY(x, y T) Vec2[T] {
	return Vec2[T]{
		X: c.X * x,
		Y: c.Y * y,
	}
}

func (c Vec2[T]) MultScalar(val T) Vec2[T] {
	return Vec2[T]{
		X: c.X * val,
		Y: c.Y * val,
	}
}

func (c Vec2[T]) Div(other Vec2[T]) Vec2[T] {
	return Vec2[T]{
		X: c.X / other.X,
		Y: c.Y / other.Y,
	}
}

func (c Vec2[T]) DivXY(x, y T) Vec2[T] {
	return Vec2[T]{
		X: c.X / x,
		Y: c.Y / y,
	}
}

func (c Vec2[T]) DivScalar(val T) Vec2[T] {
	return Vec2[T]{
		X: c.X / val,
		Y: c.Y / val,
	}
}

func (c Vec2[T]) Invert() Vec2[T] {
	return Vec2[T]{
		X: -c.X,
		Y: -c.Y,
	}
}
func (c Vec2[T]) InvertX() Vec2[T] {
	return Vec2[T]{
		X: -c.X,
		Y: c.Y,
	}
}
func (c Vec2[T]) InvertY() Vec2[T] {
	return Vec2[T]{
		X: c.X,
		Y: -c.Y,
	}
}

func (c Vec2[T]) Clamp(minX, maxX, minY, maxY T) Vec2[T] {
	return Vec2[T]{
		X: min(maxX, max(minX, c.X)),
		Y: min(maxY, max(minY, c.Y)),
	}
}
func (c Vec2[T]) ClampMin(minX, minY T) Vec2[T] {
	return Vec2[T]{
		X: max(minX, c.X),
		Y: max(minY, c.Y),
	}
}
func (c Vec2[T]) ClampMax(maxX, maxY T) Vec2[T] {
	return Vec2[T]{
		X: min(maxX, c.X),
		Y: min(maxY, c.Y),
	}
}

func (c Vec2[T]) IsInRange(minX, maxX, minY, maxY T) bool {
	return c.X >= minX && c.X <= maxX && c.Y >= minY && c.Y <= maxY
}

func (c Vec2[T]) String() string {
	return fmt.Sprintf("(%v, %v)", c.X, c.Y)
}

func (c Vec2[T]) ToInt() Vec2[int] {
	return Vec2[int]{
		X: int(c.X),
		Y: int(c.Y),
	}
}
func (c Vec2[T]) ToFloat64() Vec2[float64] {
	return Vec2[float64]{
		X: float64(c.X),
		Y: float64(c.Y),
	}
}
