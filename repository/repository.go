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
	SelectUser(id uint) (*model.User, error)
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
		// AutoMigrate errors trying to create a unique constraint on User.Name
		fmt.Println("GORM AutoMigrate returned error:", err)
	}

	return &repository{
		db: db,
	}
}

func (r *repository) InsertUser(user *model.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		fmt.Println("Error creating user in database:", result.Error)
		return errors.New("failed to create user")
	}
	return nil
}

func (r *repository) SelectUser(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("User not found")
			return nil, model.NewNotFoundError(fmt.Sprintf("User %d not found", id))
		}

		fmt.Println("Error selecting user in database:", result.Error)
		return nil, errors.New("failed to select user")
	}
	return &user, nil
}
