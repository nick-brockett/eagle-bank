package port

import (
	"eagle-bank.com/internal/core/domain/model"
)

//go:generate moq -pkg mocks -out ./mocks/account_service.go . AccountService

type AccountService interface {
	CreateAccount(newAccount *model.NewAccount) (*model.UserAccount, error)
}
