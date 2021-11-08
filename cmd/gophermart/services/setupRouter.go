package services

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/middlewares"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
	"go.uber.org/zap"
)

func SetupRouter(repo handlers.RepositoryInterface, tokenCfg *configurations.ConfigToken,
	wp *workers.WorkerPool, log *zap.SugaredLogger, accrualURL string) *gin.Engine {
	router := gin.Default()

	handler := handlers.New(repo, tokenCfg, wp, log, accrualURL)

	router.GET("/api/db/ping", handler.PingDB)
	router.POST("/api/user/register", handler.Register)
	router.POST("/api/user/login", handler.Login)
	router.POST("/api/user/refresh", handler.Refresh)
	router.POST("/api/user/orders", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret, log), handler.CreateOrder)
	router.GET("/api/user/orders", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret, log), handler.GetOrders)
	router.GET("/api/user/balance", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret, log), handler.GetBalance)
	router.POST("/api/user/balance/withdraw", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret, log), handler.MakeWithdraw)
	router.GET("/api/user/balance/withdrawals", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret, log), handler.GetWithdraws)

	router.HandleMethodNotAllowed = true

	return router
}
