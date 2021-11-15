package customerrors

import "errors"

func NewOrderAlreadyRegisterByYouError() *OrderAlreadyRegisterByYou {
	return &OrderAlreadyRegisterByYou{
		BaseError{
			Err: errors.New("order already register by you"),
		},
	}
}

type OrderAlreadyRegisterByYou struct {
	BaseError
}

func NewOrderAlreadyRegisterError() *OrderAlreadyRegister {
	return &OrderAlreadyRegister{
		BaseError{
			Err: errors.New("order already register"),
		},
	}
}

type OrderAlreadyRegister struct {
	BaseError
}
