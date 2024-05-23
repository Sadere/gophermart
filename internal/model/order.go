package model

import "github.com/Sadere/gophermart/internal/structs"

type OrderStatus string

const (
	OrderNew        OrderStatus = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	OrderProcessing OrderStatus = "PROCESSING" // — вознаграждение за заказ рассчитывается;
	OrderInvalid    OrderStatus = "INVALID"    // — система расчёта вознаграждений отказала в расчёте;
	OrderProcessed  OrderStatus = "PROCESSED"  // — данные по заказу проверены и информация о расчёте успешно получена.
)

type Order struct {
	ID        uint64          `json:"-" db:"id"`
	UserID    uint64          `json:"-" db:"user_id"`
	CreatedAt structs.RFCTime `json:"uploaded_at" db:"created_at"`
	Number    string          `json:"number" db:"number"`
	Status    OrderStatus     `json:"status" db:"status"`
	Accrual   *float64        `json:"accrual,omitempty" db:"accrual"`
}
