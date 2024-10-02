package model

import "time"

type User struct {
	ID           uint64    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password"`
}

type UserBalance struct {
	Balance   float64 `json:"current" db:"balance"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}
