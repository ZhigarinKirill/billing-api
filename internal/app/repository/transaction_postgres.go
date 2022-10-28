package repository

import (
	"database/sql"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/ZhigarinKirill/billing-api/model"
	"github.com/jmoiron/sqlx"
)

// TransactionPostgres - responsible for working with transactions and DataBase
type TransactionPostgres struct {
	db *sqlx.DB
}

// NewTransactionPostgres - TransactionPostgres constructor
func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

// Creates new transaction. Returns ID and error
func (tp *TransactionPostgres) CreateTransaction(t *model.Transaction) (int, error) {
	err := tp.db.QueryRow(
		`INSERT INTO transactions (user_id, service_id, order_id, start_date, start_bill_amount, end_bill_amount) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID`,
		t.UserID, t.ServiceID, t.OrderID, time.Now(), t.StartBillAmount, t.EndBillAmount,
	).Scan(&t.ID)

	if err != nil {
		return 0, err
	}

	return t.ID, nil
}

// Reserves money from main account into reserv. Returns error
func (tp *TransactionPostgres) ReserveMoneyFromAccount(userID, serviceID, orderID int, amount float64) error {

	up := &UserPostgres{
		db: tp.db,
	}
	if err := up.WithdrawMoneyFromAccount(userID, amount, true); err != nil {
		return err
	}

	if err := up.DepositMoneyIntoAccount(userID, amount, false); err != nil {
		return err
	}

	stringStartBillAmount, err := up.GetBillAmountByID(userID)
	if err != nil {
		return err
	}
	float64StartBillAmount, err := strconv.ParseFloat(stringStartBillAmount, 64)
	if err != nil {
		return err
	}

	startBillAmount := int(math.Round(float64StartBillAmount * 100))
	if err != nil {
		return err
	}

	transaction := &model.Transaction{
		UserID:          userID,
		ServiceID:       serviceID,
		OrderID:         orderID,
		StartBillAmount: startBillAmount,
		EndBillAmount:   startBillAmount - int(math.Round(amount*100)),
	}

	_, err = tp.CreateTransaction(transaction)
	return err
}

// Returns first reservation transaction with specifying parametres by start date
func (tp *TransactionPostgres) GetFirstReservationTransaction(userID, serviceID, orderID int, amount float64) (*model.Transaction, error) {
	transaction := &model.Transaction{}
	intAmount := int(math.Round(amount * 100))
	err := tp.db.Get(transaction, `SELECT * FROM transactions WHERE user_id = $1 AND service_id = $2 
							AND order_id = $3 AND end_bill_amount - start_bill_amount = $4 AND 
							end_date IS NULL ORDER BY start_date ASC LIMIT 1`, userID, serviceID, orderID, -intAmount)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return nil, errors.New("no pending transactions with given parameters")
	}
	return transaction, nil
}

// Completes reservation transaction. Returns 1 if successful and error
func (tp *TransactionPostgres) CompleteReservationTransaction(userID, serviceID, orderID int, amount float64) (int, error) {
	transaction, err := tp.GetFirstReservationTransaction(userID, serviceID, orderID, amount)
	if err != nil {
		return 0, err
	}

	up := &UserPostgres{
		db: tp.db,
	}
	if err := up.WithdrawMoneyFromAccount(userID, amount, false); err != nil {
		return 0, err
	}

	return tp.CompleteTransaction(transaction.ID)
}

// Completes transaction. Returns 1 if successful and error
func (tp *TransactionPostgres) CompleteTransaction(transactionID int) (int, error) {
	res, err := tp.db.Exec(
		"UPDATE transactions SET end_date=$1 WHERE id=$2",
		time.Now(),
		transactionID,
	)
	if err != nil {
		return 0, err
	}

	updatedNum, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(updatedNum), nil
}

// Report generation
func (tp *TransactionPostgres) GetReport(year int, month int) ([]*model.Report, error) {
	reports := make([]*model.Report, 0)
	err := tp.db.Select(&reports, `SELECT service_id, SUM(start_bill_amount - end_bill_amount) amount_sum FROM transactions WHERE service_id IS NOT NULL AND end_date IS NOT NULL
										AND EXTRACT(YEAR FROM end_date) = $1 
										AND EXTRACT(MONTH FROM end_date) = $2
										GROUP BY service_id`, year, month)
	if err != nil {
		return nil, err
	}
	return reports, nil
}
