package models

import "time"

type Budget struct {
	ID         int       `json:"id"`
	CategoryID int       `json:"category_id"`
	TimePeriod int       `json:"time_period"`
	UserID     int       `json:"user_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
