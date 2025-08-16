package entity

import (
	"time"
)

func NewVerificationToken(opts ...Option[*VerificationToken]) (VerificationToken, error) {
	newEntity := VerificationToken{
		expiresAt: time.Now().Add(time.Hour),
	}
	err := newEntity.Modify(opts...)
	if err != nil {
		return VerificationToken{}, err
	}
	return newEntity, nil
}

func (vt *VerificationToken) Modify(opts ...Option[*VerificationToken]) error {
	cl, err := Clone(vt)
	if err != nil {
		return err
	}
	ApplyOptions(opts, cl)

	err = validate(verificationTokenValidation{
		Token:     cl.token,
		UserID:    cl.userID,
		ExpiresAt: cl.expiresAt,
	})
	if err != nil {
		return err
	}

	*vt = *cl
	return nil
}

type VerificationToken struct {
	token     ID
	userID    ID
	expiresAt time.Time
	usedAt    *time.Time
}

type verificationTokenValidation struct {
	Token     ID        `valid:"uuid,required"`
	UserID    ID        `valid:"uuid,required"`
	ExpiresAt time.Time `valid:"required"`
}

func (vt *VerificationToken) Token() ID {
	return vt.token
}

func (vt *VerificationToken) UserID() ID {
	return vt.userID
}

func (vt *VerificationToken) ExpiresAt() time.Time {
	return vt.expiresAt
}

func WithVerificationTokenID(token ID) Option[*VerificationToken] {
	return func(vt *VerificationToken) {
		vt.token = token
	}
}

func WithVerificationTokenUserID(userID ID) Option[*VerificationToken] {
	return func(vt *VerificationToken) {
		vt.userID = userID
	}
}

func WithVerificationTokenExpiresAt(expiresAt time.Time) Option[*VerificationToken] {
	return func(vt *VerificationToken) {
		vt.expiresAt = expiresAt
	}
}

type VerificationTokenDAO struct {
	Token     ID         `db:"token"`
	UserID    ID         `db:"user_id"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
}

func (vt *VerificationToken) FromEntity() *VerificationTokenDAO {
	return &VerificationTokenDAO{
		Token:     vt.token,
		UserID:    vt.userID,
		ExpiresAt: vt.expiresAt,
		UsedAt:    vt.usedAt,
	}
}

func (vt VerificationTokenDAO) ToEntity() *VerificationToken {
	return &VerificationToken{
		token:     vt.Token,
		userID:    vt.UserID,
		expiresAt: vt.ExpiresAt,
		usedAt:    vt.UsedAt,
	}
}
