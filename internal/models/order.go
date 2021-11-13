package models

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
	Accrual    float64   `json:"accrual"`
}

type ResponseOrder struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type ResponseOrderWithAccrual struct {
	ResponseOrder
	Accrual float64 `json:"accrual"`
}

func (ro ResponseOrderWithAccrual) MarshalJSON() ([]byte, error) {
	if ro.Accrual != 0 {
		type ResponseAlias ResponseOrderWithAccrual
		aliasValue := struct {
			ResponseAlias
			UploadedAt string `json:"uploaded_at"`
		}{
			ResponseAlias: ResponseAlias(ro),
			UploadedAt:    ro.UploadedAt.Format(time.RFC3339),
		}
		return json.Marshal(aliasValue)
	}
	type ResponseAlias ResponseOrder
	aliasValue := struct {
		ResponseAlias
		UploadedAt string `json:"uploaded_at"`
	}{
		ResponseAlias: ResponseAlias(ro.ResponseOrder),
		UploadedAt:    ro.UploadedAt.Format(time.RFC3339),
	}
	return json.Marshal(aliasValue)
}
