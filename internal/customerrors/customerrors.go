package customerrors

import (
	"errors"
	"fmt"
)

type BaseError struct {
	Err error
}

func (err *BaseError) Error() string {
	return fmt.Sprintf("%v", err.Err)
}

func (err *BaseError) Unwrap() error {
	return err.Err
}

type RepeatError struct {
	BaseError
}

func NewRepeatError() *RepeatError {
	return &RepeatError{
		BaseError{
			Err: errors.New("need repeat"),
		},
	}
}
