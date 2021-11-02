package servises

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/middlewares"
)

func SetupRouter(repo handlers.RepositoryInterface, tokenCfg *configurations.ConfigToken) *gin.Engine {
	router := gin.Default()

	handler := handlers.New(repo, tokenCfg)

	router.GET("/api/db/ping", handler.PingDB)
	router.POST("/api/user/register", handler.Register)
	router.POST("/api/user/login", handler.Login)
	router.POST("/api/user/refresh", handler.Refresh)
	router.POST("/api/user/orders", middlewares.JWTMiddleware(tokenCfg.AccessTokenSecret), handler.CreateOrder)

	router.HandleMethodNotAllowed = true

	return router
}
