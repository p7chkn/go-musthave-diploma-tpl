package tasks

import (
	"context"
	"encoding/json"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type ResponseFromAccrualService struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func CheckOrderStatus(accrualURL string, log *zap.SugaredLogger,
	changeStatus func(ctx context.Context, order string, status string, accrual int) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		response, err := http.Get(accrualURL)
		if err != nil {
			log.Warn("Problem with access accrual service")
			return customerrors.NewRepeatError()
		}
		if response.StatusCode == http.StatusTooManyRequests {
			log.Warn("Accrual service overloaded")
			return customerrors.NewRepeatError()
		}
		if response.StatusCode == http.StatusInternalServerError {
			log.Warn("Accrual service is unavailable")
			return customerrors.NewRepeatError()
		}
		if response.StatusCode == http.StatusNotFound {
			log.Warn("Order not found on accrual service")
			return nil
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return customerrors.NewRepeatError()
		}
		var result ResponseFromAccrualService
		if err := json.Unmarshal(body, &result); err != nil {
			return customerrors.NewRepeatError()
		}
		if result.Status == "REGISTERED" || result.Status == "PROCESSING" {
			log.Warn("checking order not finished yet")
			return customerrors.NewRepeatError()
		}

		if err := changeStatus(ctx, result.Order, result.Status, result.Accrual); err != nil {
			return err
		}
		return nil
	}
}
