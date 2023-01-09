package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `gorm:"column:name"`
	Email string `gorm:"unique"`
	Phone string `gorm:"column:phone"`
}

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserEmailTaken = errors.New("email already taken")
)
