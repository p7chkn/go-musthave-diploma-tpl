package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/authentication"
	"go.uber.org/zap"
	"net/http"
)

func JWTMiddleware(accessSecret string, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authentication.TokenValid(c.Request, accessSecret)
		if err != nil {
			message := make(map[string]string)
			log.Warnf("Wrong request occuped on %v: %v", c.Request.RequestURI, err.Error())
			message["detail"] = err.Error()
			c.IndentedJSON(http.StatusBadRequest, message)
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
