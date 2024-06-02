package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/utils"
)

var (
	ErrOrderExists        = errors.New("order is already loaded by another user")
	ErrOrderInvalidNumber = errors.New("invalid order number")
	ErrOrdersNotAdded     = errors.New("no orders added yet")
)

type OrderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

// Загружаем номер заказа в систему, возвращет true, если пользователь уже загрузил заказ
func (s *OrderService) SaveOrderForUser(userID uint64, number string) (bool, error) {
	// Проверяем номер заказа
	validNumber := utils.CheckLuhn(number)

	if !validNumber {
		return false, ErrOrderInvalidNumber
	}

	// Проверяем загружен ли заказ с таким номером
	order, err := s.orderRepo.GetOrderByNumber(context.Background(), number)

	// Если заказ не найден, пытаемся его добавить
	if errors.Is(err, sql.ErrNoRows) {
		_, err := s.orderRepo.Create(context.Background(), model.Order{
			UserID: userID,
			Number: number,
		})

		if err != nil {
			return false, err
		}

		return false, nil
	}

	// Другие ошибки
	if err != nil {
		return false, err
	}

	// Проверяем кем был загружен заказ
	if order.UserID != userID {
		return true, ErrOrderExists
	}

	return true, nil
}

func (s *OrderService) GetOrdersByUser(userID uint64) ([]model.Order, error) {
	orders, err := s.orderRepo.GetOrdersByUser(context.Background(), userID)

	// Если не нашли заказы, отдаем ошибку
	if len(orders) == 0 {
		return nil, ErrOrdersNotAdded
	}

	// Остальные ошибки
	if err != nil {
		return nil, err
	}

	return orders, nil
}
