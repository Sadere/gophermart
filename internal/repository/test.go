package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sadere/gophermart/internal/model"
)

type TestUserRepository struct {
	RegisteredUserPwHash string
}

func (tu *TestUserRepository) Create(ctx context.Context, user model.User) (uint64, error) {
	if user.Login == "invalid" {
		return 0, errors.New("test error")
	}

	return 1000, nil
}

func (tu *TestUserRepository) GetUserByID(ctx context.Context, ID uint64) (model.User, error) {
	var user model.User

	if ID == 0 {
		return user, sql.ErrNoRows
	}

	return model.User{
		ID:    111,
		Login: "test_user",
	}, nil
}

func (tu *TestUserRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	if login == "registered_user" {
		return model.User{
			ID:           111,
			Login:        "registered_user",
			PasswordHash: tu.RegisteredUserPwHash,
		}, nil
	}

	if login == "error_user" {
		return user, errors.New("error")
	}

	return user, sql.ErrNoRows
}
