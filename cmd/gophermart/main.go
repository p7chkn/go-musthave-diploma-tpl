package main

import (
	"context"
	"database/sql"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/database/postgres"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"

	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/services"
)

func main() {

	log := logger.InitLogger()

	log.Info("Starting server")
	ctx, cancel := context.WithCancel(context.Background())

	log.Info("Starting parse configuration")
	cfg := configurations.New()

	log.Info("Finish parse configurations, starting connection to db")
	db, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	log.Info("Finish db connection")
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Starting setup db")
	services.MustSetupDatabase(db, log)

	log.Info("Finish setup db")
	repo := postgres.NewDatabase(db)
	wp := workers.New(repo, repo, &cfg.WorkerPool, log, cfg.AccrualSystemAdress)

	go func() {
		wp.Run(ctx)
	}()

	handler := services.SetupRouter(repo, repo, &cfg.Token, wp, log, cfg.AccrualSystemAdress)

	server := &http.Server{
		Addr:    cfg.ServerAdress,
		Handler: handler,
	}
	go func() {
		log.Info("Starting server")
		log.Info(server.ListenAndServe())
		cancel()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}
	server.Shutdown(context.Background())
}
