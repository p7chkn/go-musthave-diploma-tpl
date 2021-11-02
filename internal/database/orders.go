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

func (db *PostgreDataBase) GetOrders(ctx context.Context, userID string) ([]models.Order, error) {
	var resultOrders []models.Order
	sqlGetOrders := `SELECT id, user_id, number, status, uploaded_at, accrual FROM orders WHERE user_id = $1`
	rows, err := db.conn.QueryContext(ctx, sqlGetOrders, userID)
	if err != nil {
		return resultOrders, err
	}
	for rows.Next() {
		var order models.Order
		if err = rows.Scan(&order.Id, &order.UserID); err != nil {
			return resultOrders, nil
		}
		resultOrders = append(resultOrders, order)
	}
	return resultOrders, nil
}
