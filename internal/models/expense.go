package models

import "time"

type Expense struct {
	ID            int       `json:"id"`
	Amount        int       `json:"amount"`
	Description   string    `json:"description"`
	Title         string    `json:"title"`
	Date          time.Time `json:"date"`
	PaymentMethod string    `json:"payment_method"`
	Vendor        string    `json:"vendor"`
	UserID        int       `json:"user_id"`
	BillGroupID   int       `json:"bill_group_id"`
	CategoryID    int       `json:"category_id"`
}

type ExpenseWithSplits struct {
	Expense Expense `json:"expense"`
	Splits  []Split `json:"splits"`
}
