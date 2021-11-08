package database

import (
	"context"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"time"
)

func (db *PostgreDataBase) CreateWithdraw(ctx context.Context, withdraw models.Withdraw, userID string) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var enough bool
	row := tx.QueryRowContext(ctx, "SELECT (balance >= $1) FROM users WHERE id = $2 FOR UPDATE", withdraw.Sum, userID)
	if row.Err() != nil {
		return row.Err()
	}
	err = row.Scan(&enough)
	if err != nil {
		return err
	}
	if !enough {
		return customerrors.NewNotEnoughBalanceForWithdraw()
	}
	_, err = tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", withdraw.Sum, userID)
	if err != nil {
		return err
	}
	sqlCreateWithdraw := `INSERT INTO withdrawals (order_number, user_id, sum, status, processed_at) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, sqlCreateWithdraw, withdraw.OrderNumber, userID, withdraw.Sum, "PROCESSED", time.Now())
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *PostgreDataBase) GetWithdrawals(ctx context.Context, userID string) ([]models.Withdraw, error) {
	var result []models.Withdraw
	sqlGetOrders := `SELECT order_number, sum, status, processed_at FROM withdrawals
					 WHERE user_id = $1 ORDER BY processed_at ASC`
	rows, err := db.conn.QueryContext(ctx, sqlGetOrders, userID)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var order models.Withdraw
		if err = rows.Scan(&order.OrderNumber, &order.Sum, &order.Status, &order.ProcessedAt); err != nil {
			return result, nil
		}
		result = append(result, order)

	}
	return result, nil
}
