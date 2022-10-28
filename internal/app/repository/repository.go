package repository

import (
	"github.com/ZhigarinKirill/billing-api/model"
	"github.com/jmoiron/sqlx"
)

// User - user methods interface in repository layer
type User interface {
	CreateUser(u *model.User) (int, error)
	GetAllUsers() ([]*model.User, error)
	GetUserByID(userID int) (*model.User, error)
	GetBillAmountByID(userID int) (string, error)
	DepositMoneyIntoAccount(userID int, billAmount float64, isMainAccount bool) error
	WithdrawMoneyFromAccount(userID int, amount float64, isMainAccount bool) error
	TransferMoneyBetweenUsers(fromUserID int, toUserID int, amount float64) error
}

// User - transaction methods interface in repository layer
type Transaction interface {
	CreateTransaction(t *model.Transaction) (int, error)
	ReserveMoneyFromAccount(userID, serviceID, orderID int, amount float64) error
	GetFirstReservationTransaction(userID, serviceID, orderID int, amount float64) (*model.Transaction, error)
	CompleteReservationTransaction(userID, serviceID, orderID int, amount float64) (int, error)
	CompleteTransaction(transactionID int) (int, error)
	GetReport(year int, month int) ([]*model.Report, error)
}

// Repository — responsible for obtaining data from external sources
type Repository struct {
	User
	Transaction
}

// NewRepository - repository constructor
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserPostgres(db),
		Transaction: NewTransactionPostgres(db),
	}
}
