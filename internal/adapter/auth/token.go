package auth

import (
	"fmt"
	"strings"
	"time"

	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Config struct {
	APISecret     string `env:"API_SECRET, default=eagle-bank-secret"`
	TokenLifeSpan string `env:"TOKEN_LIFESPAN, default=60m"`
}

type Service struct {
	apiSecret string
	timeout   time.Duration
}

func (s *Service) ValidateSetPasswordToken(c *gin.Context) error {
	bearerToken := c.Request.Header.Get("Authorization")
	var bToken string
	if len(strings.Split(bearerToken, " ")) == 2 {

		bToken = strings.Split(bearerToken, " ")[1]
	}
	fmt.Println(bToken)
	return nil
}

func NewService(config Config) (port.AuthService, error) {
	timeout, err := time.ParseDuration(config.TokenLifeSpan)
	if err != nil {
		return nil, errors.New("invalid token lifespan format")
	}
	return &Service{
		apiSecret: config.APISecret,
		timeout:   timeout,
	}, nil
}

func (s *Service) GetTokenExpirationTime() time.Time {
	return time.Now().Add(s.timeout)
}

func (s *Service) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(s.timeout).Unix()
	claims["role"] = role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.apiSecret))
}

func (s *Service) ValidateToken(c *gin.Context) error {
	tokenString := s.ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.apiSecret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}

func (s *Service) ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func (s *Service) ExtractTokenID(c *gin.Context) (string, error) {
	userIDValue, exists := c.Get("user_id")
	if exists {
		return userIDValue.(string), nil
	}
	tokenString := s.ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.apiSecret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}
	return "", nil
}
