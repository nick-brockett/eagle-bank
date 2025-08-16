package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	eagleSortCode       = "10-10-10"
	accountPrivateType  = "private"
	accountBusinessType = "business"
)

func NewAccount(opts ...Option[*Account]) (Account, error) {
	newEntity := Account{
		sortCode:  eagleSortCode,
		createdAt: time.Now().UTC(),
	}
	err := newEntity.Modify(opts...)
	if err != nil {
		return Account{}, err
	}
	return newEntity, nil
}

func (u *Account) Modify(opts ...Option[*Account]) error {
	cl, err := Clone(u)
	if err != nil {
		return err
	}
	ApplyOptions(opts, cl)

	err = validate(accountValidation{
		AccountNumber: cl.accountNumber,
		UserID:        cl.userID,
		SortCode:      cl.sortCode,
		Name:          cl.name,
		AccountType:   cl.accountType,
		Balance:       cl.balance.String(),
		Currency:      cl.currency,
		CreatedAt:     cl.createdAt,
		UpdatedAt:     cl.updatedAt,
	})
	if err != nil {
		return err
	}

	*u = *cl
	return nil
}

type Account struct {
	accountNumber string
	userID        string
	sortCode      string
	name          string
	accountType   string
	balance       decimal.Decimal
	currency      string
	createdAt     time.Time
	updatedAt     time.Time
}

type accountValidation struct {
	AccountNumber string    `valid:"required"`
	UserID        string    `valid:"uuid,required"`
	SortCode      string    `valid:"required"`
	Name          string    `valid:"required"`
	AccountType   string    `valid:"required"`
	Balance       string    `valid:"required"`
	Currency      string    `valid:"required"`
	CreatedAt     time.Time `valid:"required"`
	UpdatedAt     time.Time `valid:"required"`
}

func (a *Account) AccountNumber() string {
	return a.accountNumber
}

func (a *Account) UserID() string {
	return a.userID
}

func (a *Account) SortCode() string {
	return a.sortCode
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) AccountType() string {
	return a.accountType
}

func (a *Account) Balance() decimal.Decimal {
	return a.balance
}

func (a *Account) Currency() string {
	return a.currency
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func WithAccountNumber(accountNumber string) Option[*Account] {
	return func(a *Account) {
		a.accountNumber = accountNumber
	}
}

func WithAccountUserID(userID string) Option[*Account] {
	return func(a *Account) {
		a.userID = userID
	}
}

func WithAccountSortCode(sortCode string) Option[*Account] {
	return func(a *Account) {
		a.sortCode = sortCode
	}
}

func WithAccountName(name string) Option[*Account] {
	return func(a *Account) {
		a.name = name
	}
}

func WithAccountType(accountType string) Option[*Account] {
	return func(a *Account) {
		a.accountType = accountType
	}
}

func WithAccountBalance(balance decimal.Decimal) Option[*Account] {
	return func(a *Account) {
		a.balance = balance
	}
}

func WithAccountCurrency(currency string) Option[*Account] {
	return func(a *Account) {
		a.currency = currency
	}
}
func WithAccountCreatedAt(createdAt time.Time) Option[*Account] {
	return func(a *Account) {
		a.createdAt = createdAt
	}
}

func WithAccountUpdatedAt(updatedAt time.Time) Option[*Account] {
	return func(a *Account) {
		a.updatedAt = updatedAt
	}
}

func (a *Account) FromEntity() AccountDAO {
	return AccountDAO{
		AccountNumber: a.accountNumber,
		UserID:        a.userID,
		SortCode:      a.sortCode,
		Name:          a.name,
		AccountType:   a.accountType,
		Balance:       a.balance,
		Currency:      a.currency,
		CreatedAt:     a.createdAt,
		UpdatedAt:     a.updatedAt,
	}
}

func (a *AccountDAO) ToEntity() *Account {
	return &Account{
		accountNumber: a.AccountNumber,
		userID:        a.UserID,
		sortCode:      a.SortCode,
		name:          a.Name,
		accountType:   a.AccountType,
		balance:       a.Balance,
		currency:      a.Currency,
		createdAt:     a.CreatedAt,
		updatedAt:     a.UpdatedAt,
	}
}

type AccountDAO struct {
	AccountNumber string          `db:"account_number"`
	UserID        string          `db:"user_id"`
	SortCode      string          `db:"sort_code"`
	Name          string          `db:"name"`
	AccountType   string          `db:"account_type"`
	Balance       decimal.Decimal `db:"balance"`
	Currency      string          `db:"currency"`
	CreatedAt     time.Time       `db:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"`
}
