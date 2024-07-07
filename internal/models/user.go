package models

import (
	"time"
)

// User represents a user in the system.

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	PaymentMethods []string  `json:"payment_methods"`
	Email          string    `json:"email" binding:"required"`
	Info           string    `json:"info"`
	Password       string    `json:"-" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
