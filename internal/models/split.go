package models

type Split struct {
	ID          int     `json:"id"`
	Paid        float64 `json:"paid"`
	Responsible float64 `json:"responsible"`
	BillMateID  int     `json:"bill_mate_id"`
	ExpenseID   int     `json:"expense_id"`
}
