package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	UserAwaitingVerificationStatus = "awaiting_verification"
	UserEmailVerifiedStatus        = "email_verified"
	UserActiveStatus               = "active"
)

func NewUser(opts ...Option[*User]) (User, error) {
	newEntity := User{
		id:        ID(uuid.NewString()),
		status:    UserAwaitingVerificationStatus,
		createdAt: time.Now().UTC(),
	}
	err := newEntity.Modify(opts...)
	if err != nil {
		return User{}, err
	}
	return newEntity, nil
}

func (u *User) Modify(opts ...Option[*User]) error {
	cl, err := Clone(u)
	if err != nil {
		return err
	}
	ApplyOptions(opts, cl)

	err = validate(userValidation{
		ID:          cl.id,
		Name:        cl.name,
		Email:       cl.email,
		PhoneNumber: cl.phoneNumber,
		Status:      cl.status,
		CreatedAt:   cl.createdAt,
	})
	if err != nil {
		return err
	}

	*u = *cl
	return nil
}

type User struct {
	id           ID
	name         string
	email        string
	phoneNumber  string
	passwordHash *string
	status       string
	createdAt    time.Time
}

type userValidation struct {
	ID          ID        `valid:"uuid,required"`
	Name        string    `valid:"required"`
	Email       string    `valid:"required"`
	PhoneNumber string    `valid:"required"`
	Status      string    `valid:"required"`
	CreatedAt   time.Time `valid:"required"`
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) PhoneNumber() string {
	return u.phoneNumber
}

func (u *User) PasswordHash() *string {
	return u.passwordHash
}

func (u *User) Status() string {
	return u.status
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func WithUserID(id ID) Option[*User] {
	return func(u *User) {
		u.id = id
	}
}

func WithUserName(name string) Option[*User] {
	return func(u *User) {
		u.name = name
	}
}

func WithUserEmail(email string) Option[*User] {
	return func(u *User) {
		u.email = email
	}
}

func WithUserPhoneNumber(phoneNumber string) Option[*User] {
	return func(u *User) {
		u.phoneNumber = phoneNumber
	}
}

func WithUserPasswordHash(passwordHash *string) Option[*User] {
	return func(u *User) {
		u.passwordHash = passwordHash
	}
}

func WithUserStatus(status string) Option[*User] {
	return func(u *User) {
		u.status = status
	}
}

func WithUserCreatedAt(createdAt time.Time) Option[*User] {
	return func(u *User) {
		u.createdAt = createdAt
	}
}

func (u *User) FromEntity() UserDAO {
	return UserDAO{
		ID:           u.id.String(),
		Name:         u.name,
		Email:        u.email,
		PhoneNumber:  u.phoneNumber,
		PasswordHash: u.passwordHash,
		Status:       u.status,
		CreatedAt:    u.createdAt,
	}
}

func (u *UserDAO) ToEntity() *User {
	return &User{
		id:           ID(u.ID),
		name:         u.Name,
		email:        u.Email,
		phoneNumber:  u.PhoneNumber,
		passwordHash: u.PasswordHash,
		status:       u.Status,
		createdAt:    u.CreatedAt,
	}
}

type UserDAO struct {
	ID           string    `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PhoneNumber  string    `db:"phone_number"`
	PasswordHash *string   `db:"password_hash"`
	Status       string    `db:"status"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
