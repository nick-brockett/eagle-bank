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
)

func main() {
	s := "Eagle Bank"
	fmt.Printf("Hello and welcome to, %s!\n", s)
	ctx := context.Background()

	// wire up the database service
	dbCfg := postgres.Config{}
	if err := envconfig.Process(ctx, &dbCfg); err != nil {
		log.Fatal(err)
	}

	dbContext, err := postgres.NewDBContext(ctx, dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbContext.Close(); err != nil {
			log.Printf("Error closing DB: %v\n", err)
		}
	}()

	// wire up the auth service
	authCfg := auth.Config{}
	if err := envconfig.Process(ctx, &authCfg); err != nil {
		log.Fatal(err)
	}

	authService, err := auth.NewService(authCfg)
	if err != nil {
		log.Fatal("Error initializing auth service", err)
	}

	userRepo := repository.NewUserRepository(dbContext)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(authService, userService)

	accountRepo := repository.NewAccountRepository(dbContext)
	accountService := service.NewAccountService(accountRepo)
	accountHandler := http.NewAccountHandler(authService, userService, accountService)

	router, err := http.NewRouter(authService, userHandler, accountHandler)
	if err != nil {
		log.Fatal("Error initializing router", "error", err)
	}

	err = router.Serve(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
