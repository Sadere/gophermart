package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/structs"
)

// Test user repo

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

// Test Order repo

type TestOrderRepository struct{}

func NewTestOrderRepository() OrderRepository {
	return &TestOrderRepository{}
}

func (r *TestOrderRepository) Create(ctx context.Context, order model.Order) (uint64, error) {
	if order.Number == "27078" {
		return 0, errors.New("error create order")
	}

	return 444, nil
}

func (r *TestOrderRepository) GetOrderByNumber(ctx context.Context, number string) (model.Order, error) {
	var order model.Order

	if number == "43513" {
		return order, errors.New("error get order by number")
	}

	if number == "56317" {
		order.UserID = 111
		return order, nil
	}

	if number == "24844" {
		order.UserID = 222
		return order, nil
	}

	return order, sql.ErrNoRows
}

func (r *TestOrderRepository) GetOrdersByUser(ctx context.Context, userID uint64) ([]model.Order, error) {
	var result []model.Order

	if userID == 111 {
		result = append(result, model.Order{
			ID:     1,
			Number: "111",
			Status: model.OrderNew,
		})
		return result, nil
	}

	// Error
	if userID == 222 {
		result = append(result, model.Order{
			ID:     2,
			Number: "222",
		})
		return result, errors.New("GetOrdersByUser() test error")
	}

	return result, nil
}

func (r *TestOrderRepository) GetPendingOrders(ctx context.Context) ([]model.Order, error) {
	var pendingOrders []model.Order

	return pendingOrders, nil
}

func (r *TestOrderRepository) UpdateOrder(ctx context.Context, order model.Order) error {

	return nil
}

// Test Balance repo

type TestBalanceRepository struct{}

func NewTestBalanceRepository() BalanceRepository {
	return &TestBalanceRepository{}
}

func (r *TestBalanceRepository) Withdraw(ctx context.Context, withdraw model.Withdrawal) error {
	if withdraw.UserID == 444 {
		return ErrInsufficientFunds
	}

	if withdraw.UserID == 555 {
		return errors.New("Withdraw() test error")
	}

	return nil
}

func (r *TestBalanceRepository) GetUserWithdrawals(ctx context.Context, userID uint64) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal

	if userID == 111 {
		withdrawals = append(withdrawals, model.Withdrawal{
			ID:     1,
			UserID: 111,
			Number: "78477",
			CreatedAt: structs.RFCTime{
				Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			Amount: 200,
		})
	}

	if userID == 222 {
		return nil, errors.New("GetUserWithdrawals() test error")
	}

	return withdrawals, nil
}

func (r *TestBalanceRepository) GetUserBalance(ctx context.Context, userID uint64) (*model.UserBalance, error) {
	var balance model.UserBalance

	balance.Balance = 200
	balance.Withdrawn = 200

	if userID == 222 {
		balance.Balance = 0
		balance.Withdrawn = 0
	}

	if userID == 333 {
		return nil, errors.New("GetUserBalance() test error")
	}

	return &balance, nil
}

func (r *TestBalanceRepository) Deposit(ctx context.Context, userID uint64, sum float64) error {

	return nil
}
