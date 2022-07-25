package repository

import "github.com/bruc3mackenzi3/microservice-demo/src/model"

type Repository interface {
	InsertUser(name string) (string, error)
	SelectUser(id string) (model.User, error)
}

type repository struct {
}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) InsertUser(name string) (string, error) {
	return "", nil
}

func (r *repository) SelectUser(id string) (model.User, error) {
	return model.User{}, nil
}
