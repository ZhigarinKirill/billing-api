package service

import (
	"errors"

	"github.com/ZhigarinKirill/billing-api/internal/app/repository"
	"github.com/ZhigarinKirill/billing-api/model"
)

// UserService - realizes methods for user logic
type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) CreateUser(user *model.User) (int, error) {
	if err := user.Validate(); err != nil {
		return 0, err
	}
	return us.repo.CreateUser(user)
}

func (us *UserService) GetAllUsers() ([]*model.User, error) {
	return us.repo.GetAllUsers()
}

func (us *UserService) GetUserByID(userID int) (*model.User, error) {
	return us.repo.GetUserByID(userID)
}

func (us *UserService) GetBillAmountByID(userID int) (string, error) {
	return us.repo.GetBillAmountByID(userID)
}

func (us *UserService) DepositMoneyIntoAccount(userID int, amount float64, isMainAccount bool) error {
	if amount <= 0 {
		return errors.New("the amount must not be negative")
	}
	return us.repo.DepositMoneyIntoAccount(userID, amount, isMainAccount)
}

func (us *UserService) WithdrawMoneyFromAccount(userID int, amount float64, isMainAccount bool) error {
	if amount <= 0 {
		return errors.New("the amount must not be negative")
	}
	return us.repo.WithdrawMoneyFromAccount(userID, amount, isMainAccount)
}

func (us *UserService) TransferMoneyBetweenUsers(fromUserID int, toUserID int, amount float64) error {
	if amount <= 0 {
		return errors.New("the amount must not be negative")
	}
	return us.repo.TransferMoneyBetweenUsers(fromUserID, toUserID, amount)
}
