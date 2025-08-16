package service

import (
	"fmt"
	"math/rand/v2"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
)

func NewAccountService(
	repo port.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

type AccountService struct {
	repo port.AccountRepository
}

func (s AccountService) CreateAccount(newAccount *model.NewAccount) (*model.UserAccount, error) {
	newAccount.AccountNumber = GenerateAccountNumber()
	return s.repo.CreateAccount(newAccount)
}

// GenerateAccountNumber returns an 8-digit account number as a string
func GenerateAccountNumber() string {
	return fmt.Sprintf("%08d", rand.IntN(100_000_000))
}
