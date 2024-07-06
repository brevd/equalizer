package models

import "database/sql"

type BillMate struct {
	ID     int           `json:"id"`
	UserID sql.NullInt64 `json:"user_id"`
	Name   string        `json:"name"`
}
