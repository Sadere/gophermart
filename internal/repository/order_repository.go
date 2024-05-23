package repository

import (
	"context"
	"time"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (uint64, error)
	GetOrderByNumber(ctx context.Context, number string) (model.Order, error)
	GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error)
}

type PgOrderRepository struct {
	db *sqlx.DB
}

func NewPgOrderRepository(db *sqlx.DB) OrderRepository {
	return &PgOrderRepository{
		db: db,
	}
}

func (r *PgOrderRepository) Create(ctx context.Context, order model.Order) (uint64, error) {
	var newOrderID uint64

	result := r.db.QueryRowContext(ctx, "INSERT INTO orders (number, user_id, created_at) VALUES ($1, $2, $3) RETURNING id",
		order.Number,
		order.UserID,
		time.Now(),
	)

	err := result.Scan(&newOrderID)

	if err != nil {
		return 0, err
	}

	return newOrderID, nil
}

func (r *PgOrderRepository) GetOrderByNumber(ctx context.Context, number string) (model.Order, error) {
	var order model.Order

	err := r.db.QueryRowxContext(ctx, "SELECT * FROM orders WHERE number = $1", number).StructScan(&order)

	return order, err
}

func (r *PgOrderRepository) GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error) {
	var result []model.Order

	sql := "SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at ASC"
	err := r.db.SelectContext(ctx, &result, sql, userID)

	if err != nil {
		return nil, err
	}

	return result, nil
}
