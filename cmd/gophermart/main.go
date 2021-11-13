package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/logger"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/database/postgres"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/tasks"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"

	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/services"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	log := logger.InitLogger()

	stdlog.SetOutput(os.Stdout)

	log.Info("Starting server")
	ctx, cancel := context.WithCancel(context.Background())

	log.Info("Starting parse configuration")
	cfg := configurations.New()

	log.Info("Finish parse configurations, starting connection to db")
	log.Info(cfg.DataBase.DataBaseURI)
	db, err := sql.Open("postgres", cfg.DataBase.DataBaseURI)
	log.Info("Finish db connection")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Info("Starting setup db")
	services.MustSetupDatabase(ctx, db, log)

	log.Info("Finish setup db")
	repo := postgres.NewDatabase(db)
	jobStore := postgres.NewJobStore(db)
	var listTask []tasks.TaskInterface
	listTask = append(listTask, tasks.NewCheckOrderStatusTask(cfg.AccrualSystemAdress, log, repo.ChangeOrderStatus))
	taskStore := tasks.NewTaskStore(listTask)

	wp := workers.New(jobStore, taskStore, &cfg.WorkerPool, log)

	go func() {
		wp.Run(ctx)
	}()

	handler := services.SetupRouter(repo, jobStore, &cfg.Token, log)

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
