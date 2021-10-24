package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/authentication"
	"net/http"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authentication.TokenValid(c.Request)
		if err != nil {
			message := make(map[string]string)
			message["detail"] = err.Error()
			c.IndentedJSON(http.StatusBadRequest, message)
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
