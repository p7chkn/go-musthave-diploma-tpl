package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHandler_Login(t *testing.T) {
	type want struct {
		code        int
		InResponse  []string
		contentType string
	}
	tests := []struct {
		name      string
		query     string
		body      string
		mockUser  models.User
		mockError error
		want      want
	}{
		{
			name:  "success test",
			query: "/api/user/login",
			mockUser: models.User{
				Login:    "test",
				Password: "test",
			},
			mockError: nil,
			body:      `{"login": "test", "password": "test"}`,
			want: want{
				code:        200,
				InResponse:  []string{"AccessToken", "RefreshToken"},
				contentType: "application/json; charset=utf-8",
			},
		},
		{
			name:  "failed test",
			query: "/api/user/login",
			mockUser: models.User{
				Login:    "test",
				Password: "test",
			},
			mockError: errors.New("invalid user"),
			body:      `{"login": "test", "password": "test"}`,
			want: want{
				code:        400,
				contentType: "application/json; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := configurations.NewTokenConfig()
			repoMock := new(MockRepositoryInterface)
			repoMock.On("CheckPassword", mock.Anything, tt.mockUser).Return(tt.mockUser, tt.mockError)
			router := SetupRouter(repoMock, &cfg)
			w := httptest.NewRecorder()
			body := strings.NewReader(tt.body)
			req, _ := http.NewRequest(http.MethodPost, tt.query, body)
			router.ServeHTTP(w, req)
			resBody, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(resBody)
			for _, field := range tt.want.InResponse {
				assert.Contains(t, string(resBody), field)
			}
			assert.Equal(t, tt.want.code, w.Code)
			assert.Equal(t, tt.want.contentType, w.Header()["Content-Type"][0])
		})
	}
}
