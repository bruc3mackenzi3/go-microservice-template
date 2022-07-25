package repository

import (
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	InsertUser(name string) (string, error)
	SelectUser(id string) (model.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository() Repository {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s",
		config.Env.PostgresHost,
		config.Env.PostgresDB,
		config.Env.PostgresUser,
		config.Env.PostgresPassword,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return &repository{
		db: db,
	}
}

func (r *repository) InsertUser(name string) (string, error) {
	return "", nil
}

func (r *repository) SelectUser(id string) (model.User, error) {
	return model.User{}, nil
}
