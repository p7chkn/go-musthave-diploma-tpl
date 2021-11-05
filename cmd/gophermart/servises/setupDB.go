package servises

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/pressly/goose"
)

func MustSetupDatabase(ctx context.Context, db *sql.DB) {
	workDirectory, _ := os.Getwd()
	migrationsPath := workDirectory + `/cmd/gophermart/servises/migrations`
	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
}
