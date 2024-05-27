package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Sadere/gophermart/internal/database"
	"github.com/Sadere/gophermart/internal/model"
	"github.com/jmoiron/sqlx"
)

var ErrInsufficientFunds = errors.New("requested sum is greater than available accrual")

type BalanceRepository interface {
	Withdraw(ctx context.Context, withdraw model.Withdrawal) error
	GetUserWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error)
	GetUserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error)
	Deposit(ctx context.Context, userID uint64, sum float64) error
}

type PgBalanceRepository struct {
	db *sqlx.DB
}

func NewPgBalanceRepository(db *sqlx.DB) BalanceRepository {
	return &PgBalanceRepository{
		db: db,
	}
}

func (r *PgBalanceRepository) Withdraw(ctx context.Context, withdraw model.Withdrawal) error {
	err := database.WrapTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		// Блокируем баланс пользователя
		var balance float64
		userQuery := "SELECT balance FROM users WHERE id = $1 FOR UPDATE"
		err := r.db.QueryRowContext(ctx, userQuery, withdraw.UserID).Scan(&balance)

		if err != nil {
			return err
		}

		if withdraw.Amount > balance {
			return ErrInsufficientFunds
		}

		// Снимаем баланс пользователя
		updateBalanceQuery := "UPDATE users SET balance = balance - $1 WHERE id = $2"
		_, err = r.db.ExecContext(ctx, updateBalanceQuery, withdraw.Amount, withdraw.UserID)

		if err != nil {
			return err
		}

		// Добавляем запись о выводе средств
		insertWithdrawalQuery := `INSERT INTO withdrawals
			(user_id, number, created_at, amount)
				VALUES
			($1, $2, $3, $4)`
		_, err = r.db.ExecContext(
			ctx,
			insertWithdrawalQuery,
			withdraw.UserID,
			withdraw.Number,
			time.Now(),
			withdraw.Amount,
		)
		if err != nil {
			return err
		}

		// Увеличиваем сумму выведенных средств пользователя
		updateWithdrawnQuery := "UPDATE users SET withdrawn = withdrawn + $1 WHERE id = $2"
		_, err = r.db.ExecContext(ctx, updateWithdrawnQuery, withdraw.Amount, withdraw.UserID)

		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *PgBalanceRepository) GetUserWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal

	selectWithdrawalsQuery := `
		SELECT
			id,
			user_id,
			number,
			created_at,
			amount
		FROM withdrawals
		WHERE user_id = $1
	`
	err := r.db.SelectContext(
		ctx,
		&withdrawals,
		selectWithdrawalsQuery,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (r *PgBalanceRepository) GetUserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error) {
	var balance model.UserBalance

	err := r.db.QueryRowxContext(ctx, "SELECT balance, withdrawn FROM users WHERE id = $1", userID).
		StructScan(&balance)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (r *PgBalanceRepository) Deposit(ctx context.Context, userID uint64, sum float64) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE users SET balance = balance + $1 WHERE id = $2",
		sum,
		userID,
	)

	return err
}
