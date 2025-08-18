package http

import (
	"net/http"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewAccountHandler(
	logger *zap.SugaredLogger,
	authService port.AuthService,
	userService port.UserService,
	accountService port.AccountService,
) AccountHandler {
	return AccountHandler{
		logger:         logger,
		authService:    authService,
		userService:    userService,
		accountService: accountService,
	}
}

type AccountHandler struct {
	logger         *zap.SugaredLogger
	authService    port.AuthService
	userService    port.UserService
	accountService port.AccountService
}

type NewAccountRequest struct {
	Name        string `json:"name" binding:"required"`
	AccountType string `json:"accountType" binding:"required"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req NewAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	userID, err := h.authService.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	//TODO: use the userService to validate the state of the User Account: e.g. Status >= Active

	account, err := h.accountService.CreateAccount(&model.NewAccount{
		UserID: userID,
		Name:   req.Name,
		Type:   req.AccountType,
	})
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	panic("not-implemented ")
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	panic("not-implemented ")
}
