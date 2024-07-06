package models

import (
	"time"
)

// User represents a user in the system.

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	PaymentMethods []string  `json:"payment_methods"`
	Email          string    `json:"email"`
	Info           string    `json:"info"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
