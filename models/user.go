package models

import (
	"errors"
	"finalproject/helpers"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"not null; unique"`
	Email     string `gorm:"not null; unique" validate:"email"`
	Password  string `gorm:"not null;" validate:"min=6"`
	Age       uint   `gorm:"not null" validate:"min=9"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var validate *validator.Validate

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	validate = validator.New()
	err = validate.Struct(u)

	if err != nil {
		return errors.New("can't save invalid data")
	}
	u.Password, err = helpers.HashPassword(u.Password)

	if err != nil {
		return errors.New("can't save invalid data")
	}

	return
}
