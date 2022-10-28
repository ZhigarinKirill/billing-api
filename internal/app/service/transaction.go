package service

import (
	"errors"

	"github.com/ZhigarinKirill/billing-api/internal/app/repository"
	"github.com/ZhigarinKirill/billing-api/model"
)

// TransactionService - realizes methods for transaction logic
type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}
func (ts *TransactionService) CreateTransaction(t *model.Transaction) (int, error) {
	return ts.repo.CreateTransaction(t)
}

func (ts *TransactionService) ReserveMoneyFromAccount(userID, serviceID, orderID int, amount float64) error {
	if amount <= 0 {
		return errors.New("the amount must not be negative")
	}
	return ts.repo.ReserveMoneyFromAccount(userID, serviceID, orderID, amount)
}

func (ts *TransactionService) GetFirstReservationTransaction(userID, serviceID, orderID int, amount float64) (*model.Transaction, error) {
	return ts.repo.GetFirstReservationTransaction(userID, serviceID, orderID, amount)
}

func (ts *TransactionService) CompleteReservationTransaction(userID, serviceID, orderID int, amount float64) (int, error) {
	if amount <= 0 {
		return 0, errors.New("the amount must not be negative")
	}
	return ts.repo.CompleteReservationTransaction(userID, serviceID, orderID, amount)
}

func (ts *TransactionService) CompleteTransaction(transactionID int) (int, error) {
	return ts.repo.CompleteTransaction(transactionID)
}

func (ts *TransactionService) GetReport(year int, month int) ([]*model.Report, error) {
	return ts.repo.GetReport(year, month)
}
