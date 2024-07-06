package models

import "time"

type Expense struct {
	ID            int       `json:"id"`
	Amount        float64   `json:"amount"`
	Description   string    `json:"description"`
	Title         string    `json:"title"`
	Date          time.Time `json:"date"`
	PaymentMethod string    `json:"payment_method"`
	Vendor        string    `json:"vendor"`
	UserID        int       `json:"user_id"`
	BillGroupID   int       `json:"bill_group_id"`
	CategoryID    int       `json:"category_id"`
}
