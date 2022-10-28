package model

type Report struct {
	ServiceID int `json:"service_id" db:"service_id"`
	Amount    int `json:"amount" db:"amount_sum"`
}
