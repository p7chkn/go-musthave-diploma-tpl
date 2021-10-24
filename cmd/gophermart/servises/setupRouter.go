package servises

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/app/middlewares"
)

func SetupRouter(repo handlers.RepositoryInterface) *gin.Engine {
	router := gin.Default()

	handler := handlers.New(repo)

	router.GET("/api/db/ping", middlewares.JWTMiddleware(), handler.PingDB)
	router.POST("/api/user/register", handler.Register)
	router.POST("/api/user/login", handler.Login)
	router.POST("/api/user/refresh", handler.Refresh)

	return router
}
