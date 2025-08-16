package entity_test

import (
	"testing"
	"time"

	"eagle-bank.com/internal/adapter/storage/postgres/repository/entity"
	"eagle-bank.com/internal/testsupport"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	id := entity.ID(uuid.NewString())
	name := gofakeit.Name()
	email := gofakeit.Email()
	phoneNumber := gofakeit.Phone()
	createdAt := testsupport.TimeNowRoundedMicroseconds()

	defaultOpts := entity.Options[*entity.User]{
		entity.WithUserID(id),
		entity.WithUserName(name),
		entity.WithUserEmail(email),
		entity.WithUserPhoneNumber(phoneNumber),
		entity.WithUserCreatedAt(createdAt),
	}

	tests := []struct {
		desc                string
		opts                entity.Options[*entity.User]
		expectedErrorString string
	}{
		{
			desc: "succeeds when valid",
			opts: defaultOpts,

			expectedErrorString: "",
		},
		{
			desc: "fails when id is empty ",
			opts: defaultOpts.Merge(entity.WithUserID("")),

			expectedErrorString: "ID: non zero value required",
		},
		{
			desc: "fails when id is not a uuid",
			opts: defaultOpts.Merge(entity.WithUserID("123")),

			expectedErrorString: "ID: 123 does not validate as uuid",
		},
		{
			desc: "fails when name is empty ",
			opts: defaultOpts.Merge(entity.WithUserName("")),

			expectedErrorString: "Name: non zero value required",
		},
		{
			desc: "fails when email is empty ",
			opts: defaultOpts.Merge(entity.WithUserEmail("")),

			expectedErrorString: "Email: non zero value required",
		},
		{
			desc: "fails when phoneNumber is empty ",
			opts: defaultOpts.Merge(entity.WithUserPhoneNumber("")),

			expectedErrorString: "PhoneNumber: non zero value required",
		},
		{
			desc: "fails when status is empty ",
			opts: defaultOpts.Merge(entity.WithUserStatus("")),

			expectedErrorString: "Status: non zero value required",
		},
		{
			desc: "fails when createdAt is empty ",
			opts: defaultOpts.Merge(entity.WithUserCreatedAt(time.Time{})),

			expectedErrorString: "CreatedAt: non zero value required",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {

			userEntity, err := entity.NewUser(tt.opts...)
			if tt.expectedErrorString == "" {
				require.NoError(t, err)
				assert.Equal(t, id, userEntity.ID())
				assert.Equal(t, name, userEntity.Name())
				assert.Equal(t, email, userEntity.Email())
				assert.Equal(t, phoneNumber, userEntity.PhoneNumber())
				assert.Equal(t, entity.UserAwaitingVerificationStatus, userEntity.Status())
				assert.Equal(t, createdAt, userEntity.CreatedAt())
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorString)
				require.Equal(t, entity.User{}, userEntity)
			}

		})
	}
}
