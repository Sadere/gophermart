package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Sadere/gophermart/internal/auth"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
)

type ErrUserExists struct {
	Login string
}

func (e *ErrUserExists) Error() string {
	return fmt.Sprintf("user with login '%s' has already registered", e.Login)
}

func (e *ErrUserExists) Is(tgt error) bool {
	target, ok := tgt.(*ErrUserExists)
	if !ok {
		return false
	}
	return e.Login == target.Login
}

var (
	ErrBadCredentials = errors.New("bad credentials")
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) RegisterUser(login string, password string) (model.User, error) {
	var newUser model.User
	_, err := s.userRepo.GetUserByLogin(context.Background(), login)

	// Проверяем существует ли пользователь с таким логином
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		return newUser, &ErrUserExists{Login: login}
	}

	// Хешируем пароль
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return newUser, errors.New("failed to generate password hash")
	}

	// Сохраняем юзера
	newUser = model.User{
		Login:        login,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	var newUserID uint64
	newUserID, err = s.userRepo.Create(context.Background(), newUser)

	if err != nil {
		return newUser, errors.New("failed to create user")
	}

	newUser.ID = newUserID

	return newUser, nil
}

func (s *UserService) LoginUser(login string, password string) (model.User, error) {
	user, err := s.userRepo.GetUserByLogin(context.Background(), login)

	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrBadCredentials
	}

	if err != nil {
		return user, errors.New("failed to authenticate user")
	}

	if !auth.CheckPassword(user.PasswordHash, password) {
		return user, ErrBadCredentials
	}

	return user, nil
}
