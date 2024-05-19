package model

import "time"

type OrderStatus string

const (
	New        OrderStatus = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	Processing OrderStatus = "PROCESSING" // — вознаграждение за заказ рассчитывается;
	Invalid    OrderStatus = "INVALID"    // — система расчёта вознаграждений отказала в расчёте;
	Processed  OrderStatus = "PROCESSED"  // — данные по заказу проверены и информация о расчёте успешно получена.
)

type Order struct {
	ID        uint64      `json:"id" db:"id"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	Number    string      `json:"number" db:"number"`
	Status    OrderStatus `json:"status" db:"status"`
	Accrual   *float64    `json:"accrual" db:"accrual"`
}
