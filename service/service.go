package service

import (
	"fmt"
	"strconv"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
)

type Service interface {
	CreateUser(user model.User) (string, error)
	GetUser(id uint) (*model.User, error)
}

type service struct {
	r repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) CreateUser(user model.User) (string, error) {
	// Check if email is in use by an existing user
	_, err := s.r.SelectUserByEmail(user.Email)
	if err == nil {
		return "", model.ErrUserEmailTaken
	} else if err != model.ErrUserNotFound {
		return "", err
	}

	err = s.r.InsertUser(&user)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created new user: %+v\n", user)
	return strconv.Itoa(int(user.ID)), err
}

func (s *service) GetUser(id uint) (*model.User, error) {
	user, err := s.r.SelectUser(id)
	if err != nil {
		fmt.Printf("Error getting user with id=%d: %v\n", id, err)
		return nil, err
	}
	fmt.Printf("Got user: %+v\n", user)

	return user, nil
}
