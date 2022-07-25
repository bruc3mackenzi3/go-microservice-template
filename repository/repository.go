package repository

import (
	"errors"
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	InsertUser(user *model.User) error
	SelectUser(id string) (*model.User, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new repository interface instance containing a db
// connection.  Panics on failure to estasblish connection with db.
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

	err = db.AutoMigrate(model.User{})
	if err != nil {
		fmt.Println("GORM AutoMigrate failed", err)
		panic("GORM AutoMigrate failed")
	}

	return &repository{
		db: db,
	}
}

func (r *repository) InsertUser(user *model.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		fmt.Println("Error creating user in database:", result.Error)
		return errors.New("Failed to create user")
	}
	return nil
}

func (r *repository) SelectUser(id string) (*model.User, error) {
	var user model.User
	// TODO figure out why passing empty string returns a record
	result := r.db.First(&user, id)
	// TODO implement NotFoundError returning a 404
	if result.Error != nil {
		fmt.Println("Error creating user in database:", result.Error)
		return nil, errors.New("Failed to select user")
	}
	return &user, nil
}
