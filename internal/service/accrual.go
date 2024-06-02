package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Sadere/gophermart/internal/model"
	"github.com/Sadere/gophermart/internal/repository"
	"github.com/Sadere/gophermart/internal/structs"
	"github.com/go-resty/resty/v2"
)

type AccrualService struct {
	balanceRepo  repository.BalanceRepository
	orderRepo    repository.OrderRepository
	accrualAddr  structs.NetAddress
	pullInterval time.Duration
}

func NewAccrualService(
	orderRepo repository.OrderRepository,
	balanceRepo repository.BalanceRepository,
	accrualAddr structs.NetAddress,
	pullInterval time.Duration,
) *AccrualService {
	return &AccrualService{
		orderRepo:    orderRepo,
		balanceRepo:  balanceRepo,
		accrualAddr:  accrualAddr,
		pullInterval: pullInterval,
	}
}

func (s *AccrualService) Pull() {
	statusMap := map[string]model.OrderStatus{
		"REGISTERED": model.OrderNew,
		"INVALID":    model.OrderInvalid,
		"PROCESSING": model.OrderProcessing,
		"PROCESSED":  model.OrderProcessed,
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))

		orders, err := s.orderRepo.GetPendingOrders(ctx)
		if err != nil {
			log.Printf("pull error: %v\n", err)
		}

		for _, order := range orders {
			// Указываем, что заказ попал в обработку
			if order.Status == model.OrderNew {
				order.Status = model.OrderProcessing

				err = s.orderRepo.UpdateOrder(context.Background(), order)
				if err != nil {
					log.Println("failed to update order: ", err)
				}
			}

			accOrder, err := s.pullAccrual(order.Number)

			if err != nil {
				log.Printf("failed to pull accrual, error: %v\n", err)
				continue
			}

			newStatus, ok := statusMap[accOrder.Status]
			if !ok {
				log.Printf("unknown status %s", accOrder.Status)
				continue
			}

			order.Status = newStatus

			if accOrder.Accrual != nil && *accOrder.Accrual > 0 {
				order.Accrual = accOrder.Accrual

				err = s.balanceRepo.Deposit(context.Background(), order.UserID, *order.Accrual)

				if err != nil {
					log.Println("failed to deposit user balance: ", err)
					continue
				}
			}

			err = s.orderRepo.UpdateOrder(context.Background(), order)
			if err != nil {
				log.Println("failed to update order: ", err)
			}

		}

		cancel()

		// Ждем интервал
		time.Sleep(s.pullInterval)
	}
}

func (s *AccrualService) pullAccrual(orderNumber string) (model.AccOrder, error) {
	var accOrder model.AccOrder

	baseURL := fmt.Sprintf(
		"http://%s",
		s.accrualAddr.String(),
	)

	client := resty.New()

	path := fmt.Sprintf("/api/orders/%s", orderNumber)

	result, err := client.R().
		SetResult(&accOrder).
		Get(baseURL + path)

	if err != nil {
		return accOrder, err
	}

	if result.StatusCode() != http.StatusOK {
		return accOrder, fmt.Errorf("received code = %d", result.StatusCode())
	}

	return accOrder, nil
}
