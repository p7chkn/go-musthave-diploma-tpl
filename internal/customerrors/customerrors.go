package customerrors

import (
	"errors"
	"fmt"
)

type ErrorWithAccrualSystem struct {
	Err   error
	Title string
}

func (err *ErrorWithAccrualSystem) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *ErrorWithAccrualSystem) Unwrap() error {
	return err.Err
}

func NewRepeatError() *ErrorWithAccrualSystem {
	return &ErrorWithAccrualSystem{
		Err:   errors.New("need repeat"),
		Title: "Need repeat",
	}
}
