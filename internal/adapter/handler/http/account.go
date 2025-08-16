package http

import (
	"net/http"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
)

func NewAccountHandler(
	authService port.AuthService,
	userService port.UserService,
	accountService port.AccountService,
) AccountHandler {
	return AccountHandler{
		authService:    authService,
		userService:    userService,
		accountService: accountService,
	}
}

type AccountHandler struct {
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
