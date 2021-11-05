package servises

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/middlewares"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
)

func SetupRouter(repo handlers.RepositoryInterface, tokenCfg *configurations.ConfigToken, wp *workers.WorkerPool) *gin.Engine {
	router := gin.Default()

	handler := handlers.New(repo, tokenCfg, wp)

	router.GET("/api/db/ping", handler.PingDB)
	router.POST("/api/user/register", handler.Register)
	router.POST("/api/user/login", handler.Login)
	router.POST("/api/user/refresh", handler.Refresh)
	router.POST("/api/user/orders", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret), handler.CreateOrder)
	router.GET("/api/user/orders", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret), handler.GetOrders)
	router.GET("/api/user/balance", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret), handler.GetBalance)

	router.HandleMethodNotAllowed = true

	return router
}
