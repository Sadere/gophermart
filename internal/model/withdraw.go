package model

import "github.com/Sadere/gophermart/internal/structs"

type Withdrawal struct {
	ID        uint64          `json:"-" db:"id"`
	UserID    uint64          `json:"-" db:"user_id"`
	OrderID   uint64          `json:"-" db:"order_id"`
	CreatedAt structs.RFCTime `json:"processed_at" db:"created_at"`
	Amount    float64         `json:"sum" db:"amount"`
	Order     Order           `json:"-" db:"order"`
}
