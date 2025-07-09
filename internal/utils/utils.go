package utils

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
