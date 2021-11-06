package services

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose"
)

func MustSetupDatabase(ctx context.Context, db *sql.DB) {
	log.Println("Enter a migrations start")
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
}
