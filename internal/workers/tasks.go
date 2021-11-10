package workers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type responseFromAccrualService struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func checkOrderStatus(accrualURL string, log *zap.SugaredLogger, orderNumber string,
	changeStatus func(ctx context.Context, order string, status string, accrual int) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		response, err := http.Get(accrualURL + orderNumber)
		if err != nil {
			log.Warn("Problem with access accrual service")
			return errors.New("problem with access accrual service")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			log.Warn("Accrual service overloaded")
			return errors.New("accrual service overloaded")
		}
		if response.StatusCode == http.StatusInternalServerError {
			log.Warn("Accrual service is unavailable")
			return errors.New("accrual service is unavailable")
		}
		if response.StatusCode == http.StatusNotFound {
			log.Warn("Order not found on accrual service")
			return errors.New("order not found on accrual service")
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return customerrors.NewRepeatError()
		}
		var result responseFromAccrualService
		if err := json.Unmarshal(body, &result); err != nil {
			return customerrors.NewRepeatError()
		}
		if result.Status == "REGISTERED" || result.Status == "PROCESSING" {
			log.Warn("checking order not finished yet")
			return errors.New("checking order not finished yet")
		}

		if err := changeStatus(ctx, result.Order, result.Status, result.Accrual); err != nil {
			log.Errorf("error on db side with update status to order: %v", err.Error())
			return err
		}
		return nil
	}
}
