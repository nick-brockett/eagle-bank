package port

import (
	"eagle-bank.com/internal/core/domain/model"
)

type UserService interface {
	CreateUser(user *model.NewUser) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByEmailVerificationToken(emailToken string) (*model.User, error)
	VerifyEmail(emailToken string) error
	SetPassword(user *model.User, password string) error
	Login(email string, password string) (*model.User, error)
}
