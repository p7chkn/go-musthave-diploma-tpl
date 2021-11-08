package utils

import (
	"context"
	"strconv"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

type FunctionStack []func(ctx context.Context) error

func NewFunctionStack() FunctionStack {
	return FunctionStack{}
}

func (fs *FunctionStack) IsEmpty() bool {
	return len(*fs) == 0
}

func (fs *FunctionStack) Push(f func(ctx context.Context) error) {
	*fs = append(*fs, f)
}

func (fs *FunctionStack) Pop() (func(ctx context.Context) error, bool) {
	if fs.IsEmpty() {
		return nil, false
	} else {
		index := len(*fs) - 1
		element := (*fs)[index]
		*fs = (*fs)[:index]
		return element, true
	}
}
