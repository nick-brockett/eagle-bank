package http

import (
	"eagle-bank.com/internal/core/port"
	"github.com/gin-gonic/gin"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

func NewRouter(
	authService port.AuthService,
	userHandler UserHandler,
	accountHandler AccountHandler,
) (*Router, error) {

	router := gin.Default()

	v1 := router.Group("/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.CreateUser)
			user.POST("/verify-email", userHandler.VerifyEmail)
			user.POST("/login", userHandler.Login)

			authUser := user.Group("/").Use(AuthMiddleware(authService))
			{
				authUser.GET("/:userId", userHandler.GetUser)
				authUser.POST("/set-password", userHandler.SetPassword)
				authUser.PATCH("/:userId", userHandler.UpdateUser)
			}
		}
		account := v1.Group("/accounts")
		{
			authAccount := account.Group("/").Use(AuthMiddleware(authService))
			{
				authAccount.POST("/", accountHandler.CreateAccount)
				authAccount.GET("/", accountHandler.ListAccounts)
				authAccount.GET("/:accountNumber", accountHandler.GetAccount)
			}

		}
	}
	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
