package model

import "github.com/go-playground/validator/v10"

type User struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" validate:"required" db:"name"`
}

var validate *validator.Validate

// Validate проверяет валидность полей пользователя
func (u *User) Validate() error {
	validate = validator.New()
	return validate.Struct(u)
}
