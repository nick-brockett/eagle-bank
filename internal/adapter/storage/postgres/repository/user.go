package repository

import (
	"database/sql"
	"fmt"
	"time"

	"eagle-bank.com/internal/adapter/storage/postgres"
	"eagle-bank.com/internal/adapter/storage/postgres/repository/dao"
	"eagle-bank.com/internal/adapter/storage/postgres/repository/entity"
	"eagle-bank.com/internal/core/domain/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

/**
 * UserRepository implements port.UserRepository interface
 * and provides access to the postgres database
 */

type UserRepository struct {
	pg *postgres.DBContext
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *postgres.DBContext) *UserRepository {
	return &UserRepository{
		db,
	}
}
func (ur *UserRepository) CreateUser(newUser *model.NewUser) (*model.User, error) {
	if newUser == nil {
		return nil, errors.New("new user cannot be nil")
	}

	newUserID := entity.ID(uuid.NewString())
	user, err := entity.NewUser(
		entity.WithUserID(newUserID),
		entity.WithUserName(newUser.Name),
		entity.WithUserEmail(newUser.Email),
		entity.WithUserPhoneNumber(newUser.PhoneNumber),
	)

	if err != nil {
		return nil, err
	}

	newUserAddressID := entity.ID(uuid.NewString())
	userAddress, err := entity.NewAddress(
		entity.WithUserAddressID(newUserAddressID),
		entity.WithUserAddressUserID(newUserID),
		entity.WithUserAddressLine1(newUser.Line1),
		entity.WithUserAddressLine2(newUser.Line2),
		entity.WithUserAddressLine3(newUser.Line3),
		entity.WithUserAddressTown(newUser.Town),
		entity.WithUserAddressCounty(newUser.County),
		entity.WithUserAddressPostcode(newUser.Postcode),
	)

	if err != nil {
		return nil, err
	}

	newVerificationToken := entity.ID(uuid.NewString())
	token, err := entity.NewVerificationToken(
		entity.WithVerificationTokenID(newVerificationToken),
		entity.WithVerificationTokenUserID(newUserID),
	)

	if err != nil {
		return nil, err
	}

	tx, err := ur.pg.DB.Beginx()
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

	userQuery := `	INSERT INTO eagle.users (id, name, email, phone_number, status, created_at) 
				VALUES (:id, :name, :email, :phone_number, :status, :created_at)`

	_, err = tx.NamedExec(userQuery, user.FromEntity())
	if err != nil {
		if pgErr := new(pq.Error); errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Unique violation
				return nil, errors.New("User with this email address already exists")
			}
		}
		return nil, errors.New("error encountered creating user ")
	}

	userAddressQuery := `	INSERT INTO eagle.addresses (id, user_id, line1, line2, line3, town, county, postcode, created_at) 
				VALUES (:id, :user_id, :line1, :line2, :line3, :town, :county, :postcode, :created_at)`

	_, err = tx.NamedExec(userAddressQuery, userAddress.FromEntity())
	if err != nil {
		return nil, err
	}

	tokenQuery := `	INSERT INTO eagle.user_verification_tokens (token, user_id, expires_at) 
				VALUES (:token, :user_id, :expires_at)`
	_, err = tx.NamedExec(tokenQuery, token.FromEntity())
	if err != nil {
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	return ur.GetUserByID(newUserID.String())
}

func (ur *UserRepository) GetUserByEmailVerificationToken(emailToken string) (*model.User, error) {
	if emailToken == "" {
		return nil, errors.New("emailToken cannot be empty")
	}
	query := `	SELECT token, user_id, expires_at, used_at 
				FROM eagle.user_verification_tokens
				WHERE token = :token
				`

	var verificationToken entity.VerificationTokenDAO
	namedStmt, err := ur.pg.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"token": emailToken,
	}
	err = namedStmt.Get(&verificationToken, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return ur.GetUserByID(verificationToken.UserID.String())
}

func (ur *UserRepository) VerifyEmail(emailToken string) error {

	query := `SELECT user_id 
				FROM eagle.user_verification_tokens  
				WHERE token = :token
				AND used_at IS NULL 
				AND expires_at > NOW()`

	var userID string
	namedStmt, err := ur.pg.DB.PrepareNamed(query)
	if err != nil {
		return err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"token": emailToken,
	}
	err = namedStmt.Get(&userID, args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invalid or expired token")
		} else {
			return errors.New("invalid or expired token")
		}
	}

	// Mark token as used
	_, err = ur.pg.DB.Exec(`
		UPDATE eagle.user_verification_tokens
		SET used_at = $1
		WHERE token = $2`, time.Now().UTC(), emailToken)
	if err != nil {
		return err
	}

	// Update user record to set status
	_, err = ur.pg.DB.Exec(`
		UPDATE eagle.users
		SET status = $1
		WHERE id = $2`, entity.UserEmailVerifiedStatus, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) SetPassword(user *model.User, hash []byte) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	userEntity, err := ur.GetEntityByID(user.ID)
	if err != nil {
		return err
	}
	if userEntity == nil {
		return errors.New("failed to get user entity")
	}
	passwordHash := string(hash)
	err = userEntity.Modify(
		entity.WithUserStatus(entity.UserActiveStatus),
		entity.WithUserPasswordHash(&passwordHash),
	)
	if err != nil {
		return err
	}

	_, err = ur.pg.DB.NamedExec(`
	      UPDATE eagle.users
	      SET password_hash = :password_hash, status = :status
	      WHERE id = :id
	  `, map[string]interface{}{
		"password_hash": string(hash),
		"status":        entity.UserActiveStatus,
		"id":            user.ID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Login(email string, password string) (string, error) {
	user, err := ur.GetUserByEmail(email)
	if err != nil || user == nil || user.PasswordHash == nil {
		return "", errors.New("invalid email or password")
	}
	// Compare the stored bcrypt hash with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	return user.ID, nil
}

func (ur *UserRepository) GetUserByID(id string) (*model.User, error) {
	query := `SELECT u.id, 
       				u.name, 
       				u.email, 
       				u.phone_number, 
       				u.status,
       				a.line1,
       				a.line2,
       				a.line3,
       				a.town,
       				a.county,
       				a.postcode
				FROM eagle.users u 
				JOIN eagle.addresses a ON a.user_id = u.id
				WHERE u.id = :user_id`

	var user dao.UserViewDAO
	namedStmt, err := ur.pg.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"user_id": id,
	}
	err = namedStmt.Get(&user, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return user.ConvertToModel(), nil
}

func (ur *UserRepository) GetUserByEmail(email string) (*entity.UserDAO, error) {
	query := `SELECT u.id, 
       				u.name, 
       				u.email, 
       				u.phone_number, 
       				u.password_hash,
       				u.status
				FROM eagle.users u 
				WHERE u.email = :email`

	var user entity.UserDAO
	namedStmt, err := ur.pg.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"email": email,
	}
	err = namedStmt.Get(&user, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}

	return &user, nil
}

func (ur *UserRepository) GetEntityByID(id string) (*entity.User, error) {
	query := `SELECT id, 
       				name, 
       				email, 
       				phone_number, 
       				password_hash,
       				status,
       				created_at
				FROM eagle.users u 
				WHERE u.id = :user_id`

	var user entity.UserDAO
	namedStmt, err := ur.pg.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	defer func(namedStmt *sqlx.NamedStmt) {
		_ = namedStmt.Close()
	}(namedStmt)

	args := map[string]interface{}{
		"user_id": id,
	}
	err = namedStmt.Get(&user, args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	return user.ToEntity(), nil
}

func (ur *UserRepository) UpdateUser(user *model.User) (*model.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	userEntity, err := ur.GetEntityByID(user.ID)
	if err != nil {
		return nil, err
	}
	if userEntity == nil {
		return nil, errors.New("failed to get user entity")
	}
	err = userEntity.Modify(
		entity.WithUserStatus(entity.UserActiveStatus),
		entity.WithUserPasswordHash(&user.Password),
	)
	if err != nil {
		return nil, err
	}

	userUpdateQuery := `	UPDATE eagle.users SET name = :name, 
	                       phone_number = :phone_number, 
	                       status = :status, 
	                       password_hash = :password, 
	                       updated_at = :updated_at 
				WHERE id = :user_id`

	namedStmt, err := ur.pg.DB.PrepareNamed(userUpdateQuery)
	if err != nil {
		return nil, err
	}

	defer namedStmt.Close()
	args := map[string]interface{}{
		"user_id":      userEntity.ID(),
		"name":         userEntity.Name(),
		"phone_number": userEntity.PhoneNumber(),
		"status":       userEntity.Status(),
		"password":     userEntity.PasswordHash(),
		"updated_at":   time.Now().UTC(),
	}

	result, err := namedStmt.Exec(args)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected < 1 {
		return nil, errors.New("failed to set password")
	}

	return ur.GetUserByID(string(userEntity.ID()))
}
