package repository

import (
	"fmt"
	"math"
	"time"

	"github.com/ZhigarinKirill/billing-api/model"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// UserPostgres - responsible for working with users and DataBase
type UserPostgres struct {
	db *sqlx.DB
}

// NewUserPostgres - UserPostgres constructor
func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

// Creates new user. Return ID and error
func (up *UserPostgres) CreateUser(user *model.User) (int, error) {
	err := up.db.QueryRow(
		`INSERT INTO users (name) VALUES ($1) RETURNING ID`,
		user.Name,
	).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// Allows view existing users. Return slice of users and error
func (up *UserPostgres) GetAllUsers() ([]*model.User, error) {

	users := make([]*model.User, 0)
	err := up.db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Returns user and error
func (up *UserPostgres) GetUserByID(userID int) (*model.User, error) {

	user := &model.User{}
	if err := up.db.Get(user, "SELECT * FROM users WHERE id=$1", userID); err != nil {
		return nil, err
	}

	return user, nil
}

// Returns user bill amount as string and error
func (up *UserPostgres) GetBillAmountByID(userID int) (string, error) {
	account := &model.Account{}
	err := up.db.Get(account, "SELECT * FROM main_account WHERE user_id=$1", userID)
	if err != nil {
		return "", err
	}

	userBillAmount := float64(account.BillAmount) / float64(100)
	return fmt.Sprintf("%.2f", userBillAmount), nil
}

// Creates user account. Returns error
func (up *UserPostgres) CreateAccountForUser(userID int, amount int, isMainAccount bool) error {
	tx, err := up.db.Beginx()
	if err != nil {
		return err
	}

	if isMainAccount {
		_, err = tx.Exec(
			"INSERT INTO main_account (user_id, bill_amount) VALUES ($1, $2)",
			userID,
			amount,
		)
	} else {
		_, err = tx.Exec(
			"INSERT INTO reserve_account (user_id, bill_amount) VALUES ($1, $2)",
			userID,
			amount,
		)
	}
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO transactions (user_id, service_id, order_id, start_date, end_date, start_bill_amount, end_bill_amount) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, 0, 0, time.Now(), time.Now(), amount, amount,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Deposits money into user account with userID. Return error
func (up *UserPostgres) DepositMoneyIntoAccount(userID int, amount float64, isMainAccount bool) error {

	intAmount := int(math.Round(amount * 100))
	account := &model.Account{}
	if isMainAccount {
		if err := up.db.Get(account, "SELECT * FROM main_account WHERE user_id=$1", userID); err != nil {
			return up.CreateAccountForUser(userID, intAmount, isMainAccount)
		}
	} else {
		if err := up.db.Get(account, "SELECT * FROM reserve_account WHERE user_id=$1", userID); err != nil {
			return up.CreateAccountForUser(userID, intAmount, isMainAccount)
		}
	}
	return up.ChangeMoneyAmountByUserID(userID, intAmount, account.BillAmount, isMainAccount)
}

// Withdraws money from user account with userID. Return error
func (up *UserPostgres) WithdrawMoneyFromAccount(userID int, amount float64, isMainAccount bool) error {

	intAmount := int(math.Round(amount * 100))
	account := &model.Account{}
	if isMainAccount {
		log.Info().Msg(fmt.Sprintf("withdraw %.2f money from user %d", amount, userID))
		if err := up.db.Get(account, "SELECT * FROM main_account WHERE user_id=$1", userID); err != nil {
			return fmt.Errorf("user with id %d doesn't have an account", userID)
		}
	} else {
		if err := up.db.Get(account, "SELECT * FROM reserve_account WHERE user_id=$1", userID); err != nil {
			return fmt.Errorf("user with id %d doesn't have an account", userID)
		}
	}

	if intAmount > account.BillAmount {
		return fmt.Errorf("user with id %d doesn't have enough money", userID)
	}

	return up.ChangeMoneyAmountByUserID(userID, -intAmount, account.BillAmount, isMainAccount)
}

// Change money amount in user account with userID. Return error
func (up UserPostgres) ChangeMoneyAmountByUserID(userID int, amount int, startBillAmount int, isMainAccount bool) error {
	tx, err := up.db.Beginx()
	if err != nil {
		return err
	}

	endBillAmount := startBillAmount + amount
	if isMainAccount {
		_, err = tx.Exec(
			"UPDATE main_account SET bill_amount=$1 WHERE user_id=$2",
			endBillAmount,
			userID,
		)
	} else {
		_, err = tx.Exec(
			"UPDATE reserve_account SET bill_amount=$1 WHERE user_id=$2",
			endBillAmount,
			userID,
		)
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	var tID int
	err = tx.QueryRow(
		`INSERT INTO transactions (user_id, start_date, start_bill_amount, end_bill_amount)
		VALUES ($1, $2, $3, $4) RETURNING ID`,
		userID, time.Now(), startBillAmount, endBillAmount,
	).Scan(&tID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec("UPDATE transactions SET end_date=$1 WHERE id=$2",
		time.Now(),
		tID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Transfer money from user with fromUserID to user with toUserID. Returns error
func (up *UserPostgres) TransferMoneyBetweenUsers(fromUserID int, toUserID int, amount float64) error {

	fromAccount := &model.Account{}
	err := up.db.Get(fromAccount, "SELECT * FROM main_account WHERE user_id=$1", fromUserID)
	if err != nil {
		return fmt.Errorf("user with id %d doesn't have an account", fromUserID)
	}

	intAmount := int(math.Round(amount * 100))

	if fromAccount.BillAmount < intAmount {
		return fmt.Errorf("user with id %d doesn't have enough money", fromUserID)
	}

	toAccount := &model.Account{}
	err = up.db.Get(toAccount, "SELECT * FROM main_account WHERE user_id=$1", toUserID)
	// If user with toUserID does not have an account
	if err != nil {
		if err = up.WithdrawMoneyFromAccount(fromUserID, amount, true); err != nil {
			return err
		}
		if err = up.CreateAccountForUser(toUserID, intAmount, true); err != nil {
			return err
		}
	}

	if err = up.WithdrawMoneyFromAccount(fromUserID, amount, true); err != nil {
		return err
	}
	if err = up.DepositMoneyIntoAccount(toUserID, amount, true); err != nil {
		return err
	}
	return nil
}
