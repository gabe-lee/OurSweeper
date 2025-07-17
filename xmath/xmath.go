package xmath

import "cmp"

func Lerp[T ~float32 | ~float64](min T, max T, percent T) T {
	delta := max - min
	add := delta * percent
	return min + add
}

func Clamp[T cmp.Ordered](minVal, val, maxVal T) T {
	return min(maxVal, max(minVal, val))
}

