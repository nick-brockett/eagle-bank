package port

import (
	"eagle-bank.com/internal/core/domain/model"
	"github.com/gin-gonic/gin"
)

//go:generate moq -pkg mocks -out ./mocks/auth_service.go . AuthService

type AuthService interface {
	GenerateTokens(userID string, role []string) (*model.TokenPair, error)
	ValidateToken(c *gin.Context) error
	ExtractTokenID(c *gin.Context) (string, error)
	ValidateSetPasswordToken(c *gin.Context) error
}
