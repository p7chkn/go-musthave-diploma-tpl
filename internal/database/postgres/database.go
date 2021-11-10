package postgres

import (
	"context"
	"database/sql"
)

type PostgreDataBase struct {
	conn *sql.DB
}

func NewDatabase(db *sql.DB) *PostgreDataBase {
	result := &PostgreDataBase{
		conn: db,
	}
	return result
}

func (db *PostgreDataBase) Ping(ctx context.Context) error {
	return db.conn.PingContext(ctx)
}
