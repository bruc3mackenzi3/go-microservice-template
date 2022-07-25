package model

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `gorm:"column:name"`
}

type NotFoundError struct {
	error
}

func NewNotFoundError(s string) NotFoundError {
	return NotFoundError{errors.New(s)}
}
