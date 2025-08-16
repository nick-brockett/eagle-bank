package http_test

import (
	netHTTP "net/http"
	"testing"

	"eagle-bank.com/internal/adapter/handler/http"
	"eagle-bank.com/internal/core/domain/model"
	"eagle-bank.com/internal/core/port"
	"eagle-bank.com/internal/core/port/mocks"
	"eagle-bank.com/internal/testsupport"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser(t *testing.T) {

	tests := []struct {
		desc        string
		authService port.AuthService
		userService port.UserService
		newUser     *model.NewUser

		expectedHttpStatus int
		expectedHttpBody   string
	}{
		{
			desc:        "empty payload",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser:     nil,

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   "Invalid request",
		},
		{
			desc:        "missing new user name",
			authService: &mocks.AuthServiceMock{},
			userService: &mocks.UserServiceMock{},
			newUser: &model.NewUser{
				Name:        "",
				Email:       gofakeit.Email(),
				PhoneNumber: gofakeit.Phone(),
				Line1:       gofakeit.StreetName(),
				Town:        gofakeit.City(),
				Postcode:    gofakeit.Zip(),
			},

			expectedHttpStatus: netHTTP.StatusBadRequest,
			expectedHttpBody:   "Invalid request",
		},
	}

	for _, tt := range tests {
		tt := tt
		testHandler := http.NewUserHandler(tt.authService, tt.userService)
		c, w := testsupport.NewTestContext(tt.newUser)

		t.Run(tt.desc, func(t *testing.T) {
			testHandler.CreateUser(c)
			require.Equal(t, tt.expectedHttpStatus, w.Code)
			require.Contains(t, w.Body.String(), tt.expectedHttpBody)
		})
	}

}
