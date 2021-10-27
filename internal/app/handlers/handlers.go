package handlers

import (
	"context"
	"fmt"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/authentication"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=RepositoryInterface --structname=MockRepositoryInterface --inpackage
type RepositoryInterface interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	CheckPassword(ctx context.Context, user models.User) (models.User, error)
}

func New(repo RepositoryInterface, tokenCfg *configurations.ConfigToken) *Handler {
	return &Handler{
		repo:     repo,
		tokenCfg: tokenCfg,
	}
}

type Handler struct {
	repo     RepositoryInterface
	tokenCfg *configurations.ConfigToken
}

func (h *Handler) PingDB(c *gin.Context) {
	err := h.repo.Ping(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.String(http.StatusOK, "")
}

func (h *Handler) Register(c *gin.Context) {
	data := models.User{}
	err := c.BindJSON(&data)
	if err != nil {
		h.handleError(c, err)
		return
	}
	errSlice := data.Validate()
	if len(errSlice) > 0 {
		h.handleErrors(c, errSlice)
		return
	}
	user, err := h.repo.CreateUser(c.Request.Context(), data)
	if err != nil {
		h.handleError(c, err)
		return
	}
	tokens, _ := authentication.CreateToken(user.Id, h.tokenCfg)
	c.IndentedJSON(http.StatusOK, tokens)
}

func (h *Handler) Login(c *gin.Context) {
	data := models.User{}
	err := c.BindJSON(&data)
	if err != nil {
		h.handleError(c, err)
		return
	}
	user, err := h.repo.CheckPassword(c.Request.Context(), data)
	if err != nil {
		h.handleError(c, err)
		return
	}
	tokens, err := authentication.CreateToken(user.Id, h.tokenCfg)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {
	data := authentication.RefreshTokenData{}
	err := c.BindJSON(&data)
	if err != nil {
		h.handleError(c, err)
		return
	}
	tokens, err := authentication.RefreshToken(data.RefreshToken, h.tokenCfg)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, tokens)
}
func (h *Handler) handleError(c *gin.Context, err error) {
	message := make(map[string]string)
	message["detail"] = err.Error()
	c.IndentedJSON(http.StatusBadRequest, message)
}

func (h *Handler) handleErrors(c *gin.Context, errorSlice []error) {
	message := make(map[string]string)
	for index, err := range errorSlice {
		message[fmt.Sprintf("Error #1: %v", index)] = err.Error()
	}
	c.IndentedJSON(http.StatusBadRequest, message)
}
