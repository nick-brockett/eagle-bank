package service

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"golang.org/x/crypto/bcrypt"

	"github.com/nyaruka/phonenumbers"
	"github.com/pkg/errors"
)

func NewUserService(
	repo port.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserService struct {
	repo port.UserRepository
}

func (s UserService) Login(email string, password string) (*model.User, error) {
	userID, err := s.repo.Login(email, password)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(userID)
}

func ValidPassword(password string) error {
	// TODO: expand upon example basic validation for password
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}
	if !strings.ContainsAny(password, "0123456789") {
		return errors.New("Password must contain a number")
	}
	return nil
}

func (s UserService) SetPassword(user *model.User, password string) error {
	if err := ValidPassword(password); err != nil {
		return errors.Wrap(err, "invalid password received")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to encrypt password")
	}

	return s.repo.SetPassword(user, hash)
}

func (s UserService) GetUserByID(id string) (*model.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s UserService) CreateUser(newUser *model.NewUser) (*model.User, error) {

	if newUser == nil {
		return nil, errors.New("new user cannot be nil")
	}
	if err := ValidateNewUser(newUser); err != nil {
		return nil, err
	}
	if err := ValidateNewUserAddress(newUser); err != nil {
		return nil, err
	}

	return s.repo.CreateUser(newUser)
}

func (s UserService) GetUserByEmailVerificationToken(emailToken string) (*model.User, error) {
	if emailToken == "" {
		return nil, errors.New("token cannot be empty")
	}
	return s.repo.GetUserByEmailVerificationToken(emailToken)
}

func (s UserService) VerifyEmail(emailToken string) error {
	if emailToken == "" {
		return errors.New("token cannot be empty")
	}
	return s.repo.VerifyEmail(emailToken)
}

func ValidateNewUser(p *model.NewUser) error {
	if !isValidEmail(p.Email) {
		return errors.New("invalid email")
	}
	if err := validatePhone(p.PhoneNumber); err != nil {
		return err
	}
	return nil
}

func ValidateNewUserAddress(p *model.NewUser) error {
	//if p.Options == nil {
	//	return errors.New("Poll options cannot be nil")
	//}
	//if lb.HasDuplicates(p.Options) {
	//	return errors.New("Poll options cannot contain duplicates")
	//}
	//if len(p.Options) < 2 {
	//	return errors.New("Poll must have at least two options")
	//}
	//if len(p.Options) > 7 {
	//	return errors.New("Poll can not have more than 7 options")
	//}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)

func validatePhone(phone string) error {
	if !e164Regex.MatchString(phone) {
		return fmt.Errorf("phone number must be in E.164 format, e.g. +14155552671")
	}

	// Parse phone number without specifying region
	parsed, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return fmt.Errorf("failed to parse phone number: %v", err)
	}

	// Check if it's a possible number
	if !phonenumbers.IsPossibleNumber(parsed) {
		return fmt.Errorf("phone number is not possible")
	}

	// Check if it's a valid number
	if !phonenumbers.IsValidNumber(parsed) {
		return fmt.Errorf("phone number is not valid")
	}

	return nil
}
