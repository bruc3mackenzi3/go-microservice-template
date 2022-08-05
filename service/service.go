package service

import (
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
)

type Service interface {
	CreateUser(user *model.User) error
	GetUser(id uint) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

type service struct {
	r repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{
		r: r,
	}
}

func (s *service) CreateUser(user *model.User) error {
	// Check if email is in use by an existing user
	_, err := s.r.SelectUserByEmail(user.Email)
	if err == nil {
		return model.ErrUserEmailTaken
	} else if err != model.ErrUserNotFound {
		return err
	}

	err = s.r.InsertUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("Created new user: %+v\n", user)
	return nil
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

func (s *service) UpdateUser(user *model.User) error {
	// Check if email is in use by another user
	existingUser, err := s.r.SelectUserByEmail(user.Email)
	if existingUser != nil && user.ID != existingUser.ID && err == nil {
		return model.ErrUserEmailTaken
	} else if err != nil && err != model.ErrUserNotFound {
		return err
	}

	err = s.r.UpdateUser(user)
	if err != nil {
		return err
	}

	fmt.Printf("Updated user: %+v\n", user)
	return nil
}

func (s *service) DeleteUser(id uint) error {
	err := s.r.DeleteUser(id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted user: %d\n", id)
	return nil
}
