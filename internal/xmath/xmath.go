package xmath

type Float interface {
	~float32 | ~float64
}

func Lerp[T ~float32 | ~float64](min T, max T, percent T) T {
	delta := max - min
	add := delta * percent
	return min + add
}
