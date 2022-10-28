package model

type Account struct {
	ID         int `json:"id" db:"id"`
	BillAmount int `json:"bill_amount" db:"bill_amount"`
	UserID     int `json:"user_id" db:"user_id"`
}
