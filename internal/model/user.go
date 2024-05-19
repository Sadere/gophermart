package model

import "time"

type User struct {
	ID           uint64    `json:"id" db:"id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-" db:"password"`
}
