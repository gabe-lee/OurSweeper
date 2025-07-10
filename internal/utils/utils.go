package utils

import "cmp"

type ErrorChecker struct {
	Err error
}

func (e *ErrorChecker) IsErr(err error) bool {
	e.Err = err
	return err != nil
}

type FirstError struct {
	Err error
}

func (e *FirstError) Add(err error) {
	if e.Err == nil {
		e.Err = err
	}
}

func AddOffset(x, y int, off [2]int) (xx, yy int) {
	return x + off[0], y + off[1]
}

var NEAR_X = [8]int{
	-1, 0, 1,
	-1, 1,
	-1, 0, 1,
}
var NEAR_Y = [8]int{
	-1, -1, -1,
	0, 0,
	1, 1, 1,
}

type NearbyCoords struct {
	X   [8]int
	Y   [8]int
	Len int
}

func GetNearbyCoords(minX, minY, x, y, maxX, maxY int) NearbyCoords {
	near := NearbyCoords{}
	for i := range NEAR_X {
		xx := x + NEAR_X[i]
		yy := y + NEAR_Y[i]
		if xx >= minX && xx <= maxX && yy >= minY && yy <= maxY {
			near.X[near.Len] = xx
			near.Y[near.Len] = yy
			near.Len += 1
		}
	}
	return near
}

func Clamp[T cmp.Ordered](minVal, val, maxVal T) T {
	return min(maxVal, max(minVal, val))
}
