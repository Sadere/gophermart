package service

import (
	"context"
	"errors"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/utils"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrOrderNotFound     = errors.New("order not found")
)

type BalanceService struct {
	balanceRepo repository.BalanceRepository
}

func NewBalanceService(balanceRepo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		balanceRepo: balanceRepo,
	}
}

func (s *BalanceService) RegisterWithdraw(userID uint64, orderNumber string, sum float64) error {
	// Проверяем валидность номера
	if !utils.CheckLuhn(orderNumber) {
		return ErrOrderInvalidNumber
	}

	// Получаем баланс пользователя
	userBalance, err := s.balanceRepo.GetUserBalance(context.Background(), userID)
	if err != nil {
		return err
	}

	// Проверяем достаточно ли на балансе
	if sum > userBalance.Balance {
		return ErrInsufficientFunds
	}

	withdrawRequest := model.Withdrawal{
		UserID: userID,
		Number: orderNumber,
		Amount: sum,
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
