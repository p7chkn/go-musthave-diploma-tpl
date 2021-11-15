package customerrors

import "errors"

type NotEnoughBalanceForWithdraw struct {
	BaseError
}

func NewNotEnoughBalanceForWithdraw() *NotEnoughBalanceForWithdraw {
	return &NotEnoughBalanceForWithdraw{
		BaseError{
			Err: errors.New("not enough balance"),
		},
	}
}
