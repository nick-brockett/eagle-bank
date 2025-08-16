package http

import (
	"net/http"

	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(s port.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := s.ValidateToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		userID, err := s.ExtractTokenID(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
