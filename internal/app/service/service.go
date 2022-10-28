package service

import (
	"github.com/ZhigarinKirill/billing-api/internal/app/repository"
	"github.com/ZhigarinKirill/billing-api/model"
)

// User - user methods interface in service layer
type User interface {
	CreateUser(u *model.User) (int, error)
	GetAllUsers() ([]*model.User, error)
	GetUserByID(userID int) (*model.User, error)
	GetBillAmountByID(userID int) (string, error)
	DepositMoneyIntoAccount(userID int, billAmount float64, isMainAccount bool) error
	WithdrawMoneyFromAccount(userID int, amount float64, isMainAccount bool) error
	TransferMoneyBetweenUsers(fromUserID int, toUserID int, amount float64) error
}

// Transaction - transaction methods interface in service layer
type Transaction interface {
	CreateTransaction(t *model.Transaction) (int, error)
	ReserveMoneyFromAccount(userID, serviceID, orderID int, amount float64) error
	GetFirstReservationTransaction(userID, serviceID, orderID int, amount float64) (*model.Transaction, error)
	CompleteReservationTransaction(userID, serviceID, orderID int, amount float64) (int, error)
	CompleteTransaction(transactionID int) (int, error)
	GetReport(year int, month int) ([]*model.Report, error)
}

// Service — responsible for business logic
type Service struct {
	User
	Transaction
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User:        NewUserService(repo.User),
		Transaction: NewTransactionService(repo.Transaction),
	}
}
