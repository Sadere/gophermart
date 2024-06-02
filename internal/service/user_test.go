package service

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestErrUserExists(t *testing.T) {
	t.Run("test error message", func(t *testing.T) {
		login := "test_login_123"
		err := ErrUserExists{Login: login}
		expectedMsg := fmt.Sprintf("user with login '%s' has already registered", login)

		assert.Equal(t, expectedMsg, err.Error())
	})

	t.Run("test comparison true", func(t *testing.T) {
		login := "test_login_123"
		err := &ErrUserExists{Login: login}

		result := errors.Is(err, err)

		assert.True(t, result)
	})

	t.Run("test comparison false", func(t *testing.T) {
		login := "test_login_123"
		err := &ErrUserExists{Login: login}

		result := errors.Is(err, &ErrUserExists{Login: "test_login_432"})

		assert.False(t, result)
	})
}

func TestRegisterUser(t *testing.T) {
	testPassword := "test_password_123"
	registeredUserPwHash, err := auth.HashPassword(testPassword)

	assert.NoError(t, err, "Failed to generate test password")

	repo := &repository.TestUserRepository{
		RegisteredUserPwHash: registeredUserPwHash,
	}
	userService := NewUserService(repo)

	type want struct {
		user model.User
		err  bool
	}

	tests := []struct {
		name     string
		login    string
		password string
		want     want
	}{
		{
			name:     "success register",
			login:    "test_user1",
			password: "test_pw",
			want: want{
				user: model.User{
					ID:    1000,
					Login: "test_user1",
				},
				err: false,
			},
		},
		{
			name:     "user exists",
			login:    "registered_user",
			password: "test_pw",
			want: want{
				user: model.User{},
				err:  true,
			},
		},
		{
			name:     "register too long password",
			login:    "test_user1",
			password: strings.Repeat("A", 100),
			want: want{
				user: model.User{},
				err:  true,
			},
		},
		{
			name:     "error user",
			login:    "error_user",
			password: "test_pw",
			want: want{
				user: model.User{},
				err:  true,
			},
		},
		{
			name:     "error creating user",
			login:    "invalid",
			password: "test_pw",
			want: want{
				user: model.User{Login: "invalid"},
				err:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resUser, err := userService.RegisterUser(tt.login, tt.password)

			assert.Equal(t, tt.want.user.ID, resUser.ID)
			assert.Equal(t, tt.want.user.Login, resUser.Login)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	testPassword := "test_password_123"
	registeredUserPwHash, err := auth.HashPassword(testPassword)

	assert.NoError(t, err, "Failed to generate test password")

	repo := &repository.TestUserRepository{
		RegisteredUserPwHash: registeredUserPwHash,
	}
	userService := NewUserService(repo)

	type want struct {
		user model.User
		err  bool
	}

	tests := []struct {
		name     string
		login    string
		password string
		want     want
	}{
		{
			name:     "success login",
			login:    "registered_user",
			password: testPassword,
			want: want{
				user: model.User{
					ID:    111,
					Login: "registered_user",
				},
				err: false,
			},
		},
		{
			name:     "wrong password",
			login:    "registered_user",
			password: "wrong_pw_12345",
			want: want{
				user: model.User{
					ID:    111,
					Login: "registered_user",
				},
				err: true,
			},
		},
		{
			name:     "unknown user",
			login:    "unknown_user_123",
			password: "test_pw",
			want: want{
				user: model.User{},
				err:  true,
			},
		},
		{
			name:     "error user",
			login:    "error_user",
			password: "test_pw",
			want: want{
				user: model.User{},
				err:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resUser, err := userService.LoginUser(tt.login, tt.password)

			assert.Equal(t, tt.want.user.ID, resUser.ID)
			assert.Equal(t, tt.want.user.Login, resUser.Login)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
