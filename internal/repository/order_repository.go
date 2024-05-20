package repository

import (
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
}

type PgOrderRepository struct {
	db *sqlx.DB
}

func NewPgOrderRepository(db *sqlx.DB) OrderRepository {
	return &PgOrderRepository{
		db: db,
	}
}
