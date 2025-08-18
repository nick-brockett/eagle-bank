package http_test

import (
	"encoding/json"
	netHTTP "net/http"
	"testing"

	"eagle-bank.com/internal/adapter/handler/http"
	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port/mocks"
	"eagle-bank.com/internal/testsupport"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser(t *testing.T) {

	testUser := model.User{
		ID:          uuid.NewString(),
		Name:        gofakeit.Name(),
		Email:       gofakeit.Email(),
		PhoneNumber: gofakeit.Phone(),
		Line1:       gofakeit.StreetName(),
		Town:        gofakeit.City(),
		Postcode:    gofakeit.Zip(),
	}

	validTestUserBytes, err := json.Marshal(testUser)
	require.NoError(t, err)

	tests := []struct {
		desc        string
		authService *mocks.AuthServiceMock
		userService *mocks.UserServiceMock
		newUser     *model.NewUser

		expectedHttpStatus                 int
		expectedHttpBody                   string
		expectedCreateUserServiceCallCount int
	}{
		{
			desc:        "empty payload",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser:     nil,

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing user name",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        "",
				Email:       testUser.Email,
				PhoneNumber: testUser.PhoneNumber,
				Line1:       testUser.Line1,
				Town:        testUser.Town,
				Postcode:    testUser.Postcode,
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing email",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        testUser.Name,
				Email:       "",
				PhoneNumber: testUser.PhoneNumber,
				Line1:       testUser.Line1,
				Town:        testUser.Town,
				Postcode:    testUser.Postcode,
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing phone number",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        testUser.Name,
				Email:       testUser.Email,
				PhoneNumber: "",
				Line1:       testUser.Line1,
				Town:        testUser.Town,
				Postcode:    testUser.Postcode,
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing address line1",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        testUser.Name,
				Email:       testUser.Email,
				PhoneNumber: testUser.PhoneNumber,
				Line1:       "",
				Town:        testUser.Town,
				Postcode:    testUser.Postcode,
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing town",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        testUser.Name,
				Email:       testUser.Email,
				PhoneNumber: testUser.PhoneNumber,
				Line1:       testUser.Line1,
				Town:        "",
				Postcode:    testUser.Postcode,
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "missing postcode",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        testUser.Name,
				Email:       testUser.Email,
				PhoneNumber: testUser.PhoneNumber,
				Line1:       testUser.Line1,
				Town:        testUser.Town,
				Postcode:    "",
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   `{"message":"Invalid request"}`,
		},
		{
			desc:        "internal service error",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{
				CreateUserFunc: func(user *model.NewUser) (*model.User, error) {
					return nil, errors.New("test internal service error")
				},
			},
			newUser: &model.NewUser{
				Name:        gofakeit.Name(),
				Email:       gofakeit.Email(),
				PhoneNumber: gofakeit.Phone(),
				Line1:       gofakeit.StreetName(),
				Town:        gofakeit.City(),
				Postcode:    gofakeit.Zip(),
			},

			expectedHttpStatus:                 netHTTP.StatusInternalServerError,
			expectedHttpBody:                   `{"message":"test internal service error"}`,
			expectedCreateUserServiceCallCount: 1,
		},
		{
			desc:        "success",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{
				CreateUserFunc: func(user *model.NewUser) (*model.User, error) {
					return &testUser, nil
				},
			},
			newUser: &model.NewUser{
				Name:        gofakeit.Name(),
				Email:       gofakeit.Email(),
				PhoneNumber: gofakeit.Phone(),
				Line1:       gofakeit.StreetName(),
				Town:        gofakeit.City(),
				Postcode:    gofakeit.Zip(),
			},

			expectedHttpStatus:                 netHTTP.StatusCreated,
			expectedHttpBody:                   string(validTestUserBytes),
			expectedCreateUserServiceCallCount: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		testHandler := http.NewUserHandler(tt.authService, tt.userService)
		c, w := testsupport.NewTestContext(tt.newUser)

		t.Run(tt.desc, func(t *testing.T) {
			testHandler.CreateUser(c)
			assert.Equal(t, tt.expectedHttpStatus, w.Code)
			assert.JSONEq(t, w.Body.String(), tt.expectedHttpBody)

			require.Equal(t, tt.expectedCreateUserServiceCallCount, len(tt.userService.CreateUserCalls()))

		})
	}

}
