package entity

import (
	"time"
)

func NewUserAccount(opts ...Option[*UserAccount]) (UserAccount, error) {
	newEntity := UserAccount{
		createdAt: time.Now().UTC(),
	}
	err := newEntity.Modify(opts...)
	if err != nil {
		return UserAccount{}, err
	}
	return newEntity, nil
}

func (u *UserAccount) Modify(opts ...Option[*UserAccount]) error {
	cl, err := Clone(u)
	if err != nil {
		return err
	}
	ApplyOptions(opts, cl)

	err = validate(userAccountValidation{
		ID:            cl.id,
		UserID:        cl.userID,
		AccountNumber: cl.accountNumber,
		CreatedAt:     cl.createdAt,
	})
	if err != nil {
		return err
	}

	*u = *cl
	return nil
}

type UserAccount struct {
	id            string
	userID        string
	accountNumber string
	createdAt     time.Time
}

type userAccountValidation struct {
	ID            string    `valid:"uuid,required"`
	UserID        string    `valid:"required"`
	AccountNumber string    `valid:"required"`
	CreatedAt     time.Time `valid:"required"`
}

func (a *UserAccount) ID() string {
	return a.id
}

func (a *UserAccount) UserID() string {
	return a.userID
}

func (a *UserAccount) AccountNumber() string {
	return a.accountNumber
}

func (a *UserAccount) CreatedAt() time.Time {
	return a.createdAt
}

func WithUserAccountID(id string) Option[*UserAccount] {
	return func(a *UserAccount) {
		a.id = id
	}
}

func WithUserAccountUserID(userID string) Option[*UserAccount] {
	return func(a *UserAccount) {
		a.userID = userID
	}
}

func WithUserAccountNumber(accountNumber string) Option[*UserAccount] {
	return func(a *UserAccount) {
		a.accountNumber = accountNumber
	}
}

func WithUserAccountCreatedAt(createdAt time.Time) Option[*UserAccount] {
	return func(a *UserAccount) {
		a.createdAt = createdAt
	}
}

func (a *UserAccount) FromEntity() UserAccountDAO {
	return UserAccountDAO{
		ID:            a.id,
		UserID:        a.userID,
		AccountNumber: a.accountNumber,
		CreatedAt:     a.createdAt,
	}
}

func (a *UserAccountDAO) ToEntity() *UserAccount {
	return &UserAccount{
		id:            a.ID,
		userID:        a.UserID,
		accountNumber: a.AccountNumber,
		createdAt:     a.CreatedAt,
	}
}

type UserAccountDAO struct {
	ID            string    `db:"id"`
	UserID        string    `db:"user_id"`
	AccountNumber string    `db:"account_number"`
	CreatedAt     time.Time `db:"created_at"`
}
