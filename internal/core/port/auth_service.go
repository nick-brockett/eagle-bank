package port

import (
	"time"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(c *gin.Context) error
	ExtractTokenID(c *gin.Context) (string, error)
	GetTokenExpirationTime() time.Time
	ValidateSetPasswordToken(c *gin.Context) error
}
