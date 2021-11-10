package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func NewCheckOrderStatusTask(accrualURL string, log *zap.SugaredLogger,
	changeStatus func(ctx context.Context, order string, status string, accrual int) error) *CheckOrderStatusTask {
	return &CheckOrderStatusTask{
		accrualURL:   accrualURL,
		log:          log,
		changeStatus: changeStatus,
	}
}

type responseFromAccrualService struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

type CheckOrderStatusTask struct {
	accrualURL   string
	log          *zap.SugaredLogger
	changeStatus func(ctx context.Context, order string, status string, accrual int) error
}

func (os *CheckOrderStatusTask) GetTitle() string {
	return "CheckOrderStatus"
}

func (os *CheckOrderStatusTask) CreateFunction(parameters map[string]string) (func(ctx context.Context) error, error) {
	orderNumber, ok := parameters["order_number"]
	if !ok {
		return nil, errors.New("wrong parameters")
	}
	return func(ctx context.Context) error {
		response, err := http.Get(os.accrualURL + orderNumber)
		if err != nil {
			os.log.Warn("Problem with access accrual service")
			return errors.New("problem with access accrual service")
		}
		if response.StatusCode == http.StatusTooManyRequests {
			os.log.Warn("Accrual service overloaded")
			return errors.New("accrual service overloaded")
		}
		if response.StatusCode == http.StatusInternalServerError {
			os.log.Warn("Accrual service is unavailable")
			return errors.New("accrual service is unavailable")
		}
		if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusNoContent {
			os.log.Warn("Order not found on accrual service")
			return errors.New("order not found on accrual service")
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		var result responseFromAccrualService
		if err := json.Unmarshal(body, &result); err != nil {
			return err
		}
		if result.Status == "REGISTERED" || result.Status == "PROCESSING" {
			os.log.Warn("checking order not finished yet")
			return errors.New("checking order not finished yet")
		}

		if err := os.changeStatus(ctx, result.Order, result.Status, result.Accrual); err != nil {
			os.log.Errorf("error on db side with update status to order: %v", err.Error())
			return err
		}
		return nil
	}, nil
}
