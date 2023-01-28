package integrationtests

import (
	"fmt"
	"os"

	"github.com/bruc3mackenzi3/microservice-demo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	if os.Getenv("INTEGRATION_TESTS") == "" {
		return
	}

	connectionString := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database")
	}
}

func deleteUser(email string) {
	// Unscoped() is used to override GORM's soft delete feature
	db.Unscoped().Where("email = ?", email).Delete(&model.User{})
}
