package port

import (
	"eagle-bank.com/internal/adapter/storage/postgres/repository/entity"
	"eagle-bank.com/internal/core/domain/model"
)

//go:generate moq -pkg mocks -out ./mocks/user_repository.go . UserRepository

type UserRepository interface {
	CreateUser(newUser *model.NewUser) (*model.User, error)
	UpdateUser(user *model.User) (*model.User, error)
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*entity.UserDAO, error)
	GetUserByEmailVerificationToken(emailToken string) (*model.User, error)
	VerifyEmail(emailToken string) error
	SetPassword(user *model.User, hash []byte) error
	Login(email string, password string) (string, error)
	//Update(user *model.User) error
	//Delete(id string) error
}
