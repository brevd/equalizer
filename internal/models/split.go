package models

type Split struct {
	ID          int `json:"id"`
	Paid        int `json:"paid"`
	Responsible int `json:"responsible"`
	BillMateID  int `json:"bill_mate_id"`
	ExpenseID   int `json:"expense_id"`
}
