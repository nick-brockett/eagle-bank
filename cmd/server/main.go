package main

import (
	"context"
	"fmt"
	"log"

	"eagle-bank.com/internal/adapter/auth"
	"eagle-bank.com/internal/adapter/handler/http"
	"eagle-bank.com/internal/adapter/storage/postgres"
	"eagle-bank.com/internal/adapter/storage/postgres/repository"
	"eagle-bank.com/internal/core/service"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	s := "Eagle Bank"
	fmt.Printf("Hello and welcome to, %s!\n", s)

	ctx := context.Background()

	// TODO: trial new logging framework
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize zap logger: %v", err)
	}
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	// wire up the database service
	dbCfg := postgres.Config{}
	if err := envconfig.Process(ctx, &dbCfg); err != nil {
		logger.Fatalw("failed to derive database credentials", "error", err)
	}

	dbContext, err := postgres.NewDBContext(ctx, dbCfg)
	if err != nil {
		logger.Fatalw("failed to connect to database", "error", err)
	}

	defer func() {
		if err := dbContext.Close(); err != nil {
			logger.Errorw("error closing DB", "error", err)
		}
	}()

	// wire up the auth service
	authCfg := auth.Config{}
	if err := envconfig.Process(ctx, &authCfg); err != nil {
		logger.Fatalw("failed to load auth config", "error", err)
	}

	authService, err := auth.NewService(authCfg)
	if err != nil {
		logger.Fatalw("failed to initialise auth service", "error", err)
	}

	userRepo := repository.NewUserRepository(dbContext)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(logger, authService, userService)

	accountRepo := repository.NewAccountRepository(dbContext)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := http.NewAccountHandler(logger, authService, userService, accountService)

	router, err := http.NewRouter(authService, userHandler, accountHandler)
	if err != nil {
		logger.Fatalw("error initializing router", "error", err)
	}

	logger.Infow("server starting", "port", 8080)
	err = router.Serve(":8080")
	if err != nil {
		logger.Fatalw("server error", "error", err)
	}

}
