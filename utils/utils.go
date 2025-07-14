package utils

import "errors"

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

type ErrorCollector struct {
	Err error
}

func (e *ErrorCollector) Do(err error) {
	e.Err = errors.Join(e.Err, err)
}
