package repository

import (
	"context"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (uint64, error)
	GetUserByID(ctx context.Context, ID uint64) (model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
}

type PgUserRepository struct {
	db *sqlx.DB
}

func NewPgUserRepository(db *sqlx.DB) UserRepository {
	return &PgUserRepository{
		db: db,
	}
}

// Creates new user and returns new user id
func (r *PgUserRepository) Create(ctx context.Context, user model.User) (uint64, error) {
	var newUserID uint64

	result := r.db.QueryRowContext(ctx, "INSERT INTO users (login, password, created_at) VALUES ($1, $2, $3) RETURNING id",
		user.Login,
		user.PasswordHash,
		user.CreatedAt,
	)

	err := result.Scan(&newUserID)

	if err != nil {
		return 0, err
	}

	return newUserID, nil
}

func (r *PgUserRepository) GetUserByID(ctx context.Context, ID uint64) (model.User, error) {
	var user model.User

	err := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", ID).StructScan(&user)

	return user, err
}

func (r *PgUserRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User

	err := r.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE login = $1", login).StructScan(&user)

	return user, err
}
