package services

import (
	"database/sql"
	"go.uber.org/zap"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

func MustSetupDatabase(db *sql.DB, log *zap.SugaredLogger) {
	log.Info("Enter a migrations start")
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	migrationsPath := basePath + "/migrations"
	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
}
