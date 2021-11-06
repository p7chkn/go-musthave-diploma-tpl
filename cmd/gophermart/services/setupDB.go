package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose"
)

func MustSetupDatabase(ctx context.Context, db *sql.DB) {
	fmt.Println("Enter a migrations start")
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
}
