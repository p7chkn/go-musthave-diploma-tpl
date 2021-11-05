package handlers

import (
	"context"
	"fmt"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/authentication"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/workers"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=RepositoryInterface --structname=MockRepositoryInterface --inpackage
type RepositoryInterface interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	CheckPassword(ctx context.Context, user models.User) (models.User, error)
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrders(ctx context.Context, userID string) ([]interface{}, error)
	GetBalance(ctx context.Context, userID string) (models.UserBalance, error)
}

func New(repo RepositoryInterface, tokenCfg *configurations.ConfigToken, wp *workers.WorkerPool) *Handler {
	return &Handler{
		repo:     repo,
		tokenCfg: tokenCfg,
		wp:       wp,
	}
}

type Handler struct {
	repo     RepositoryInterface
	tokenCfg *configurations.ConfigToken
	wp       *workers.WorkerPool
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
	tokens, _ := authentication.CreateToken(user.ID, h.tokenCfg)
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
	tokens, err := authentication.CreateToken(user.ID, h.tokenCfg)
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

func (h *Handler) CreateOrder(c *gin.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		h.handleError(c, err)
		return
	}

	order := models.Order{
		UserID: c.GetString("userID"),
		Number: string(body),
		Status: "NEW",
	}

	err = h.repo.CreateOrder(c.Request.Context(), order)
	// Тут еще нужно учесть вариации ошибок:
	// - регистрация заказа, который ты уже регистрировал
	// - регистрация заказа, который регистрировал кто-то другой
	// - валидация номера заказа
	// Вопрос:  нужна ли тут собственного типа ошибка?
	//
	// Так же идея еще в том, чтобы в этом проекте был воркерпулл
	// чтобы в фоне здесь отправить задачи о запросе во внешнюю систему,
	// для обработки заказа, в зависимости от ответа системы, изменять статус заказа
	// если статус заказа не окончательный, повторить запрос через N секунд, если статус окончательный,
	// освободить воркера и вероятно, еще нужно добавить поле о дате завершения заказа в таблицу заказов

	if err != nil {
		h.handleError(c, err)
		return
	}

	c.String(http.StatusAccepted, "")
}

func (h *Handler) GetOrders(c *gin.Context) {
	orders, err := h.repo.GetOrders(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, orders)
}

func (h *Handler) GetBalance(c *gin.Context) {
	balance, err := h.repo.GetBalance(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, balance)
}

func (h *Handler) MakeWithdraw(c *gin.Context) {
	// Тут в задании довольно непонятно:
	// - либо мы тут создаем новый заказ, но с отрицательныйм accrual
	// - либо мы ищем текущий заказ (среди зарегестрированных) и устанавливаем ему accrual
	// Отправлять ли тут так же в сторонний сервис запросы о состоянии заказа?
}

func (h *Handler) GetWithdraws(c *gin.Context) {
	// Тут несколько зависит от роута выше. Пока не сильно ясно.
	// тут поле называется sum а в другом месте accrual специально?
	// по сути, это же та же сущность, что и заказы, просто те из них, что с отрицательным accrual?
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
