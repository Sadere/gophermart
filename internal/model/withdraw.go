package model

import "github.com/Sadere/gophermart/internal/structs"

type Withdrawal struct {
	ID        uint64          `json:"-" db:"id"`
	UserID    uint64          `json:"-" db:"user_id"`
	Number    string          `json:"order" db:"number"`
	CreatedAt structs.RFCTime `json:"processed_at" db:"created_at"`
	Amount    float64         `json:"sum" db:"amount"`
}
