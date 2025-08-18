package repository

import (
	"fmt"
	"time"

	"eagle-bank.com/internal/adapter/storage/postgres"
	"eagle-bank.com/internal/adapter/storage/postgres/repository/entity"
	"eagle-bank.com/internal/core/domain/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

/**
 * AccountRepository implements port.AccountRepository interface
 * and provides access to the postgres database
 */

type AccountRepository struct {
	pg *postgres.DBContext
}

// NewAccountRepository creates a new account repository instance
func NewAccountRepository(db *postgres.DBContext) *AccountRepository {
	return &AccountRepository{
		db,
	}
}

// TODO: inject a clock into this method for ease of testing

func (ar *AccountRepository) CreateAccount(newAccount *model.NewAccount) (*model.UserAccount, error) {

	account, err := entity.NewAccount(
		entity.WithAccountUserID(newAccount.UserID),
		entity.WithAccountNumber(newAccount.AccountNumber),
		entity.WithAccountBalance(decimal.Zero),
		entity.WithAccountName(newAccount.Name),
		entity.WithAccountType(newAccount.Type),
		entity.WithAccountCurrency("GBP"),
		entity.WithAccountCreatedAt(time.Now()),
		entity.WithAccountUpdatedAt(time.Now()),
	)
	if err != nil {
		return nil, err
	}

	userAccount, err := entity.NewUserAccount(
		entity.WithUserAccountID(uuid.NewString()),
		entity.WithUserAccountUserID(newAccount.UserID),
		entity.WithUserAccountNumber(newAccount.AccountNumber),
	)
	if err != nil {
		return nil, err
	}

	tx, err := ar.pg.DB.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure rollback if commit is not successful
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	accountQuery := `	INSERT INTO eagle.accounts (account_number, sort_code, name, account_type, balance, currency, created_at) 
				VALUES (:account_number, :sort_code, :name, :account_type, :balance, :currency, :created_at)`

	_, err = tx.NamedExec(accountQuery, account.FromEntity())
	if err != nil {
		if pgErr := new(pq.Error); errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Unique violation
				return nil, errors.New("Account with this number already exists")
			}
		}
		return nil, errors.New("error encountered creating account ")
	}

	userAccountQuery := `	INSERT INTO eagle.user_accounts (id, user_id, account_number, created_at) 
				VALUES (:id, :user_id, :account_number, :created_at)`

	_, err = tx.NamedExec(userAccountQuery, userAccount.FromEntity())
	if err != nil {
		if pgErr := new(pq.Error); errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Unique violation
				return nil, errors.New(" User Account with this number already exists")
			}
		}
		return nil, errors.New("error encountered creating account ")
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return ar.GetAccountByNumber(account.AccountNumber())
}

func (ar *AccountRepository) GetAccountByNumber(accountNumber string) (*model.UserAccount, error) {
	query := `SELECT id, user_id, account_number FROM eagle.user_accounts WHERE account_number = :account_number`
	var userAccount entity.UserAccountDAO
	namedStmt, err := ar.pg.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"account_number": accountNumber,
	}
	err = namedStmt.Get(&userAccount, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return &model.UserAccount{
		UserID:        userAccount.UserID,
		AccountNumber: userAccount.AccountNumber,
	}, nil
}
