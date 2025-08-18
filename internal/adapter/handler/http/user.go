package http

import (
	"net/http"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewUserHandler(
	logger *zap.SugaredLogger,
	authService port.AuthService,
	userService port.UserService,
) UserHandler {
	return UserHandler{
		logger:      logger,
		authService: authService,
		userService: userService,
	}
}

type ErrorResponse struct {
	Message string
}

type UserHandler struct {
	logger      *zap.SugaredLogger
	authService port.AuthService
	userService port.UserService
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required,uuid"`
}

type SetPasswordRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	h.logger.Infow("Login handler started")
	var request LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
	}
	user, err := h.userService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorised"})
		return
	}
	tokens, err := h.authService.GenerateTokens(user.ID, []string{"create_account", "deposit", "withdraw"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login successful",
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expires":      tokens.AccessExpiry.Unix(),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	h.logger.Infow("GetUser handler started")
	userIDParam := c.Param("userId")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId format"})
		return
	}
	tokenUserID, err := h.authService.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if tokenUserID != userID.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	user, err := h.userService.GetUserByID(userID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.Infow("CreateUser handler started")
	var newUser model.NewUser
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}
	valid, err := newUser.Valid()
	if !valid || err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
		})
		return
	}
	user, err := h.userService.CreateUser(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	// TODO: suggested flow is that after creating a new User account at Eagle Bank
	// that eagle bank then send an email with a verification link to the email provided at CreateUser step.

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.logger.Infow("UpdateUser handler started")
	panic("not-implemented")
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	h.logger.Infow("VerifyEmail handler started")
	var req VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.GetUserByEmailVerificationToken(req.Token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email verification token"})
		return
	}

	err = h.userService.VerifyEmail(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email verification token"})
		return
	}

	tokens, err := h.authService.GenerateTokens(user.ID, []string{"set-password"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Email verified successfully",
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expires":      tokens.AccessExpiry.Unix(),
	})
}

func (h *UserHandler) SetPassword(c *gin.Context) {
	h.logger.Infow("SetPassword handler started")
	var req SetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := h.authService.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Retrieve user record by ID sent in jwt token
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//Ensure email accounts match
	if req.Email != user.Email {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	// Check correct token is being presented
	if err := h.authService.ValidateSetPasswordToken(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	err = h.userService.SetPassword(user, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "set password error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully"})
}
