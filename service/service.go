package service

import (
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUser(id uint) (*model.User, error)
	GetUsers() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

type userService struct {
	r repository.Repository
}

func NewUserService(r repository.Repository) UserService {
	return &userService{
		r: r,
	}
}

func (s *userService) CreateUser(user *model.User) error {
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

func (s *userService) GetUser(id uint) (*model.User, error) {
	user, err := s.r.SelectUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUsers() ([]model.User, error) {
	users, err := s.r.SelectUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) UpdateUser(user *model.User) error {
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

func (s *userService) DeleteUser(id uint) error {
	err := s.r.DeleteUser(id)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted user: %d\n", id)
	return nil
}
