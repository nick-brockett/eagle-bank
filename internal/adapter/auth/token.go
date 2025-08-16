package auth

import (
	"fmt"
	"strings"
	"time"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Config struct {
	APISecret          string `env:"API_SECRET, default=eagle-bank-secret"`
	AccessTokenExpiry  string `env:"ACCESS_TOKEN_EXPIRY, default=15m"`
	RefreshTokenExpiry string `env:"REFRESH_TOKEN_EXPIRY, default=60m"`
}

type Service struct {
	apiSecret          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewService(config Config) (port.AuthService, error) {

	accessTokenExpiry, err := time.ParseDuration(config.AccessTokenExpiry)
	if err != nil {
		return nil, errors.New("invalid access token expiry format")
	}

	refreshTokenExpiry, err := time.ParseDuration(config.RefreshTokenExpiry)
	if err != nil {
		return nil, errors.New("invalid refresh token expiry format")
	}

	return &Service{
		apiSecret:          config.APISecret,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}, nil
}

func (s *Service) ValidateSetPasswordToken(c *gin.Context) error {
	bearerToken := c.Request.Header.Get("Authorization")
	var bToken string
	if len(strings.Split(bearerToken, " ")) == 2 {

		bToken = strings.Split(bearerToken, " ")[1]
	}
	//TODO inspect claims
	_ = bToken
	return nil
}

func (s *Service) getAccessTokenExpirationTime() time.Time {
	return time.Now().Add(s.accessTokenExpiry)
}

func (s *Service) getRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(s.accessTokenExpiry)
}

func (s *Service) GenerateTokens(userID string, roles []string) (*model.TokenPair, error) {

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     s.getAccessTokenExpirationTime(),
		"iat":     time.Now().Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.apiSecret))

	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     s.getRefreshTokenExpirationTime(),
		"iat":     time.Now().Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.apiSecret))
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  s.getAccessTokenExpirationTime(),
		RefreshExpiry: s.getRefreshTokenExpirationTime(),
	}, nil
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

	//TODO: consider the whole blacklist and revoking mechanism for each token IDs (REDIS)
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
