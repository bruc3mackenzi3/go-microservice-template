package service

import (
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/src/model"
	"github.com/bruc3mackenzi3/microservice-demo/src/repository"
)

type Service interface {
	CreateUser(name string) (string, error)
	GetUser(id string) (model.User, error)
}

type service struct {
	r repository.Repository
}

func (s *service) CreateUser(name string) (string, error) {
	fmt.Println("DEBUG: Service CreateUser called")
	return "", nil
}

func (s *service) GetUser(id string) (model.User, error) {
	fmt.Println("DEBUG: Service GetUser called")
	return model.User{}, nil
}

func NewService() Service {
	return &service{
		r: repository.NewRepository(),
	}
}
