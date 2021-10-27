package main

import (
	"context"
	"database/sql"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"

	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/servises"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := configurations.New()

	db, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}
	servises.MustSetupDatabase(ctx, db)

	repo := database.NewDatabaseRepository(db)
	handler := servises.SetupRouter(repo, &cfg.Token)

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
