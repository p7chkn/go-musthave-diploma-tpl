package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupRouter(repo RepositoryInterface, tokenCfg *configurations.ConfigToken) *gin.Engine {
	router := gin.Default()

	handler := New(repo, tokenCfg)

	router.GET("/api/db/ping", handler.PingDB)
	router.POST("/api/user/register", handler.Register)
	router.POST("/api/user/login", handler.Login)
	router.POST("/api/user/refresh", handler.Refresh)

	router.HandleMethodNotAllowed = true

	return router
}

func TestHandler_PingDB(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name   string
		query  string
		mockDB error
		want   want
	}{
		{
			name:   "success test",
			query:  "/api/db/ping",
			mockDB: nil,
			want: want{
				code:        200,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "failed test",
			query:  "/api/db/ping",
			mockDB: errors.New("error with DB"),
			want: want{
				code:        500,
				response:    "error with DB",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := configurations.NewTokenConfig()
			repoMock := new(MockRepositoryInterface)
			repoMock.On("Ping", mock.Anything).Return(tt.mockDB)
			router := SetupRouter(repoMock, &cfg)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, tt.query, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
