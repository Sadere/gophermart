package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrOrderNotFound     = errors.New("order not found")
)

type BalanceService struct {
	orderRepo   repository.OrderRepository
	balanceRepo repository.BalanceRepository
}

func NewBalanceService(orderRepo repository.OrderRepository, balanceRepo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		orderRepo:   orderRepo,
		balanceRepo: balanceRepo,
	}
}

func (s *BalanceService) RegisterWithdraw(userID uint64, orderNumber string, sum float64) error {
	order, err := s.orderRepo.GetOrderByNumber(context.Background(), orderNumber)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrOrderNotFound
	}

	if err != nil {
		return err
	}

	// Проверяем принадлежность заказа пользователю
	if order.UserID != userID {
		return ErrOrderNotFound
	}

	// Проверяем хватает ли баллов
	if order.Accrual == nil || sum > *order.Accrual {
		return ErrInsufficientFunds
	}

	withdrawRequest := model.Withdrawal{
		UserID:  userID,
		OrderID: order.ID,
		Amount:  sum,
	}
	err = s.balanceRepo.Withdraw(context.Background(), withdrawRequest)

	if errors.Is(err, repository.ErrInsufficientFunds) {
		return ErrInsufficientFunds
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *BalanceService) ListUserWithdrawals(userID uint64) ([]model.Withdrawal, error) {
	withdrawals, err := s.balanceRepo.GetUserWithdrawals(context.Background(), userID)

	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (s *BalanceService) GetUserBalance(userID uint64) (*model.UserBalance, error) {
	balance, err := s.balanceRepo.GetUserBalance(context.Background(), userID)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
