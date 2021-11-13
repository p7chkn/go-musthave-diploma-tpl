package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/p7chkn/go-musthave-diploma-tpl/cmd/gophermart/configurations"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/authentication"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/customerrors"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/models"
	"github.com/p7chkn/go-musthave-diploma-tpl/internal/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//go:generate mockery --name=RepositoryInterface --structname=MockRepositoryInterface --inpackage
type RepositoryInterface interface {
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, user models.User) (*models.User, error)
	CheckPassword(ctx context.Context, user models.User) (models.User, error)
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrders(ctx context.Context, userID string) ([]models.ResponseOrderWithAccrual, error)
	GetBalance(ctx context.Context, userID string) (models.UserBalance, error)
	CreateWithdraw(ctx context.Context, withdraw models.Withdraw, userID string) error
	GetWithdrawals(ctx context.Context, userID string) ([]models.Withdraw, error)
	ChangeOrderStatus(ctx context.Context, order string, status string, accrual int) error
}

//go:generate mockery --name=JobStoreInterface --structname=MockJobStoreInterface --inpackage
type JobStoreInterface interface {
	AddJob(ctx context.Context, job models.JobStoreRow) error
}

func New(repo RepositoryInterface, jobStore JobStoreInterface, tokenCfg *configurations.ConfigToken,
	log *zap.SugaredLogger) *Handler {
	return &Handler{
		repo:     repo,
		jobStore: jobStore,
		tokenCfg: tokenCfg,
		log:      log,
	}
}

type Handler struct {
	repo     RepositoryInterface
	jobStore JobStoreInterface
	tokenCfg *configurations.ConfigToken
	log      *zap.SugaredLogger
}

func (h *Handler) PingDB(c *gin.Context) {
	err := h.repo.Ping(c.Request.Context())
	if err != nil {
		h.log.Errorf("Error occuped on %v: %v", c.Request.RequestURI, err.Error())
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
	tokens, err := authentication.CreateToken(user.ID, h.tokenCfg)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.Header("Authorization", "Bearer "+tokens.AccessToken)
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
	c.Header("Authorization", "Bearer "+tokens.AccessToken)
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

	number, err := strconv.Atoi(string(body))

	if err != nil {
		h.handleError(c, err)
		return
	}

	if !utils.ValidLuhnNumber(number) {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}

	order := models.Order{
		UserID: c.GetString("userID"),
		Number: strconv.Itoa(number),
		Status: "NEW",
	}

	err = h.repo.CreateOrder(c.Request.Context(), order)

	if err != nil {
		var selfRegisterError *customerrors.OrderAlreadyRegisterByYou
		if errors.As(err, &selfRegisterError) {
			c.String(http.StatusOK, "")
			return
		}
		var RegisterError *customerrors.OrderAlreadyRegister
		if errors.As(err, &RegisterError) {
			c.String(http.StatusConflict, "")
		}
		h.handleError(c, err)
		return
	}
	parameters, err := json.Marshal(&models.CheckOrderStatusParameters{
		OrderNumber: strconv.Itoa(number),
	})
	if err != nil {
		h.handleError(c, err)
		return
	}
	job := models.JobStoreRow{
		Type:            "CheckOrderStatus",
		NextTimeExecute: time.Now(),
		Count:           0,
		Parameters:      string(parameters),
	}
	err = h.jobStore.AddJob(c.Request.Context(), job)
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
	withdraw := models.Withdraw{}
	err := c.BindJSON(&withdraw)
	if err != nil {
		h.handleError(c, err)
		return
	}
	number, err := strconv.Atoi(withdraw.OrderNumber)
	if err != nil {
		h.handleError(c, err)
		return
	}
	if !utils.ValidLuhnNumber(number) {
		c.String(http.StatusUnprocessableEntity, "")
		return
	}
	err = h.repo.CreateWithdraw(c.Request.Context(), withdraw, c.GetString("userID"))
	if err != nil {
		var balanceError *customerrors.NotEnoughBalanceForWithdraw
		if errors.As(err, &balanceError) {
			c.String(http.StatusPaymentRequired, "")
			return
		}
		h.handleError(c, err)
		return
	}
	c.String(http.StatusOK, "")
}

func (h *Handler) GetWithdraws(c *gin.Context) {
	withdrawals, err := h.repo.GetWithdrawals(c.Request.Context(), c.GetString("userID"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, withdrawals)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	message := make(map[string]string)
	h.log.Warnf("Wrong request occuped on %v: %v", c.Request.RequestURI, err.Error())
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
