package main

import (
	"context"
	"database/sql"
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

	connURL := "postgresql://postgres:1234@localhost:5432?sslmode=disable"
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatal(err)
	}
	servises.MustSetupDatabase(ctx, db)

	repo := database.NewDatabaseRepository(db)
	handler := servises.SetupRouter(repo)

	server := &http.Server{
		Addr:    "localhost:8080",
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
