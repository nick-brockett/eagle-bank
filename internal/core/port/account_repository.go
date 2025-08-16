package port

import (
	"eagle-bank.com/internal/core/domain/model"
)

type AccountRepository interface {
	CreateAccount(newAccount *model.NewAccount) (*model.UserAccount, error)
}
