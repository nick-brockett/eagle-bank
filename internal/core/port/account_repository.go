package port

import (
	"eagle-bank.com/internal/core/domain/model"
)

//go:generate moq -pkg mocks -out ./mocks/account_repository.go . AccountRepository

type AccountRepository interface {
	CreateAccount(newAccount *model.NewAccount) (*model.UserAccount, error)
}
