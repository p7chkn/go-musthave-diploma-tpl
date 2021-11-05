package main

import (
	"context"
	"database/sql"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"

	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/services"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := configurations.New()

	db, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}
	services.MustSetupDatabase(ctx, db)

	wp := workers.New(ctx, cfg.WorkerPool.NumOfWorkers, cfg.WorkerPool.PoolBuffer)

	go func() {
		wp.Run(ctx)
	}()

	repo := database.NewDatabaseRepository(db)
	handler := services.SetupRouter(repo, &cfg.Token, wp)

	server := &http.Server{
		Addr:    cfg.ServerAdress,
		Handler: handler,
	}
	go func() {
		log.Println(server.ListenAndServe())
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
