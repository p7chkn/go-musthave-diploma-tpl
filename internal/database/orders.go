package database

import (
	"context"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
)

func (db *PostgreDataBase) CreateOrder(ctx context.Context, order models.Order) error {
	sqlCreateOrder := `INSERT INTO orders (user_id, number, status, accrual) VALUES ($1, $2, $3, $4)`
	_, err := db.conn.ExecContext(ctx, sqlCreateOrder, order.UserID, order.Number, order.Status, order.Accrual)
	return err
}

func (db *PostgreDataBase) GetOrders(ctx context.Context, userID string) ([]interface{}, error) {

	var result []interface{}
	sqlGetOrders := `SELECT number, status, uploaded_at, accrual FROM orders
					 WHERE user_id = $1 ORDER BY uploaded_at ASC`
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
