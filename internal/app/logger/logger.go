package logger

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return logger.Sugar()
}
