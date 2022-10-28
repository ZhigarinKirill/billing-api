package model

import (
	"database/sql"
	"time"
)

type Transaction struct {
	ID              int          `json:"id" db:"id"`
	UserID          int          `json:"user_id" db:"user_id"`
	ServiceID       int          `json:"service_id" db:"service_id"`
	OrderID         int          `json:"order_id" db:"order_id"`
	StartDate       time.Time    `json:"start_date" db:"start_date"`
	EndDate         sql.NullTime `json:"end_date" db:"end_date"` // may be null!!!
	StartBillAmount int          `json:"start_bill_amount" db:"start_bill_amount"`
	EndBillAmount   int          `json:"end_bill_amount" db:"end_bill_amount"`
}
