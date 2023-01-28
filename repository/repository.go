package repository

import (
	"errors"
	"fmt"

	"github.com/bruc3mackenzi3/microservice-demo/config"
	"github.com/bruc3mackenzi3/microservice-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository interface {
	InsertUser(user *model.User) error
	SelectUser(id uint) (*model.User, error)
	SelectUserByEmail(email string) (*model.User, error)
	SelectUsers() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
		fmt.Printf("Failure running query to create user in database: %T %s\n", result.Error, result.Error)
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
			return nil, model.ErrUserNotFound
		}

		fmt.Println("Failure running query to select user in database:", result.Error)
		return nil, errors.New("failed to select user")
	}
	return &user, nil
}

func (r *repository) SelectUsers() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("No users found")
			return users, nil
		}

		fmt.Println("Failure running query to select user in database:", result.Error)
		return nil, errors.New("failed to select user")
	}
	return users, nil
}

func (r *repository) SelectUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Unscoped().Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}

		return nil, errors.New("failed to select user")
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *model.User) error {
	result := r.db.Model(&user).Updates(&user)
	if result.Error != nil {
		fmt.Printf("Failure running query to create user in database: %T %s\n", result.Error, result.Error)
		return errors.New("failed to create user")
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *repository) DeleteUser(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	if result.Error != nil {
		fmt.Println("Failure running query to delete user in database:", result.Error)
		return errors.New("failed to select user")
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
