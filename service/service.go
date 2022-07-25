package service

import (
	"fmt"
	"strconv"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"github.com/bruc3mackenzi3/microservice-demo/repository"
)

type Service interface {
	CreateUser(name string) (string, error)
	GetUser(id string) (*model.User, error)
}

type service struct {
	r repository.Repository
}

func NewService() Service {
	return &service{
		r: repository.NewRepository(),
	}
}

func (s *service) CreateUser(name string) (string, error) {
	user := model.User{Name: name}
	err := s.r.InsertUser(&user)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created new user: %+v\n", user)

	return strconv.Itoa(int(user.ID)), err
}

func (s *service) GetUser(id string) (*model.User, error) {
	fmt.Println("DEBUG: Service GetUser called")

	user, err := s.r.SelectUser(id)
	if err != nil {
		fmt.Println("Error getting user:", err)
	}

	return user, nil
}
