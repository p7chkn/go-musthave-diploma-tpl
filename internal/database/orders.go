package database

import (
	"context"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgreDataBase) CreateOrder(ctx context.Context, order models.Order) error {
	sqlCreateOrder := `INSERT INTO orders (user_id, number, status, accrual) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.ExecContext(ctx, sqlCreateOrder, order.UserID, order.Number, order.Status, order.Accrual)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == pgerrcode.UniqueViolation {
			existingOrder, err := db.getOrder(ctx, order.Number)
			if err != nil {
				return err
			}
			if existingOrder.UserID == order.UserID {
				return customerrors.NewOrderAlreadyRegisterByYouError()
			}
			return customerrors.NewOrderAlreadyRegisterError()
		}
	}

	return err
}

func (db *PostgreDataBase) GetOrders(ctx context.Context, userID string) ([]models.ResponseOrderWithAccrual, error) {

	var result []models.ResponseOrderWithAccrual
	sqlGetOrders := `SELECT number, status, uploaded_at, accrual FROM orders
					 WHERE user_id = $1 ORDER BY uploaded_at`
	rows, err := db.conn.QueryContext(ctx, sqlGetOrders, userID)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var order models.ResponseOrderWithAccrual
		if err = rows.Scan(&order.Number, &order.Status, &order.UploadedAt, &order.Accrual); err != nil {
			return result, nil
		}
		result = append(result, order)

	}

	return result, nil
}

func (db *PostgreDataBase) getOrder(ctx context.Context, number string) (*models.Order, error) {
	resultOrder := &models.Order{}
	sqlGetUser := `SELECT id, user_id, number, status, accrual, uploaded_at FROM orders WHERE number = $1`
	query := db.conn.QueryRowContext(ctx, sqlGetUser, number)
	err := query.Scan(&resultOrder.ID, &resultOrder.UserID, &resultOrder.Number, &resultOrder.Status,
		&resultOrder.Accrual, &resultOrder.UploadedAt)
	if err != nil {
		return resultOrder, err
	}
	return resultOrder, nil
}

func (db *PostgreDataBase) ChangeOrderStatus(ctx context.Context, order string, status string, accrual int) error {
	return nil
}
