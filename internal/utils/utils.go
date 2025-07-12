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



func Clamp[T cmp.Ordered](minVal, val, maxVal T) T {
	return min(maxVal, max(minVal, val))
}
