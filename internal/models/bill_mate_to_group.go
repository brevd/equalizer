package models

type BillMateToGroup struct {
	ID          int `json:"id"`
	BillMateID  int `json:"bill_mate_id"`
	BillGroupID int `json:"bill_group_id"`
}
