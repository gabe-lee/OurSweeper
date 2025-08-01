package vec2

type AABB[T number] struct {
	XMin T
	YMin T
	XMax T
	YMax T
}

func AABB_FromPosSize[T number](pos Vec2[T], size Vec2[T]) AABB[T] {
	return AABB[T]{
		XMin: pos.X,
		YMin: pos.Y,
		XMax: pos.X + size.X,
		YMax: pos.Y + size.Y,
	}
}

func (aabb AABB[T]) ToPosSize() (pos Vec2[T], size Vec2[T]) {
	pos = Vec2[T]{X: aabb.XMin, Y: aabb.YMin}
	size = Vec2[T]{X: aabb.XMax - aabb.XMin, Y: aabb.YMax - aabb.YMin}
	return
}

func (aabb AABB[T]) ShrinkScalar(val T) AABB[T] {
	return AABB[T]{
		XMin: aabb.XMin + val,
		XMax: aabb.XMax - val,
		YMin: aabb.YMin + val,
		YMax: aabb.YMax - val,
	}
}

func (aabb AABB[T]) GrowBotRightScalar(val T) AABB[T] {
	return AABB[T]{
		XMin: aabb.XMin,
		XMax: aabb.XMax + val,
		YMin: aabb.YMin,
		YMax: aabb.YMax + val,
	}
}

func (aabb AABB[T]) Move(delta Vec2[T]) AABB[T] {
	return AABB[T]{
		XMin: aabb.XMin + delta.X,
		XMax: aabb.XMax + delta.X,
		YMin: aabb.YMin + delta.Y,
		YMax: aabb.YMax + delta.Y,
	}
}

func (aabb AABB[T]) Combine(other AABB[T]) AABB[T] {
	return AABB[T]{
		XMin: min(aabb.XMin, other.XMin),
		XMax: max(aabb.XMax, other.XMax),
		YMin: min(aabb.YMin, other.YMin),
		YMax: max(aabb.YMax, other.YMax),
	}
}

func (aabb AABB[T]) TopLeft() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin,
		Y: aabb.YMin,
	}
}

func (aabb AABB[T]) TopCenter() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin + ((aabb.XMax - aabb.XMin) / 2),
		Y: aabb.YMin,
	}
}

func (aabb AABB[T]) TopRight() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMax,
		Y: aabb.YMin,
	}
}

func (aabb AABB[T]) MidLeft() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin,
		Y: aabb.YMin + ((aabb.YMax - aabb.YMin) / 2),
	}
}

func (aabb AABB[T]) MidCenter() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin + ((aabb.XMax - aabb.XMin) / 2),
		Y: aabb.YMin + ((aabb.YMax - aabb.YMin) / 2),
	}
}

func (aabb AABB[T]) MidRight() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMax,
		Y: aabb.YMin + ((aabb.YMax - aabb.YMin) / 2),
	}
}

func (aabb AABB[T]) BotLeft() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin,
		Y: aabb.YMax,
	}
}

func (aabb AABB[T]) BotCenter() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMin + ((aabb.XMax - aabb.XMin) / 2),
		Y: aabb.YMax,
	}
}

func (aabb AABB[T]) BotRight() Vec2[T] {
	return Vec2[T]{
		X: aabb.XMax,
		Y: aabb.YMax,
	}
}

func (aabb AABB[T]) Equals(other AABB[T], epsilon T) bool {
	return aabb.XMax == other.XMax && aabb.XMin == other.XMax && aabb.YMin == other.YMin && aabb.YMax == other.YMax
}

func (aabb AABB[T]) ApproxEqual(other AABB[T], epsilon T) bool {
	delta := aabb.XMin - other.XMin
	if delta < 0 {
		delta = other.XMin - aabb.XMin
	}
	if delta > epsilon {
		return false
	}
	delta = aabb.YMin - other.YMin
	if delta < 0 {
		delta = other.YMin - aabb.YMin
	}
	if delta > epsilon {
		return false
	}
	delta = aabb.XMax - other.XMax
	if delta < 0 {
		delta = other.XMax - aabb.XMax
	}
	if delta > epsilon {
		return false
	}
	delta = aabb.YMax - other.YMax
	if delta < 0 {
		delta = other.YMax - aabb.YMax
	}
	return delta <= epsilon
}

func (aabb AABB[T]) PointIsWithin(point Vec2[T]) bool {
	return point.X >= aabb.XMin && point.X <= aabb.XMax && point.Y >= aabb.YMin && point.Y <= aabb.YMax
}
