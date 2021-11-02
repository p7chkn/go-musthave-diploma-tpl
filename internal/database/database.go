package database

import (
	"context"
	"database/sql"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
	"log"
)

type PostgreDataBase struct {
	conn *sql.DB
}

func NewDatabaseRepository(db *sql.DB) handlers.RepositoryInterface {
	return handlers.RepositoryInterface(NewDatabase(db))
}

func NewDatabase(db *sql.DB) *PostgreDataBase {
	result := &PostgreDataBase{
		conn: db,
	}
	return result
}

func (db *PostgreDataBase) Ping(ctx context.Context) error {

	err := db.conn.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
