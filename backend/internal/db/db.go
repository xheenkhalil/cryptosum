package db

import (
	"fmt"
	"log"

	"github.com/cryptosum/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully opened")
	
	err = DB.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Portfolio{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	
	fmt.Println("Database migration completed")
}
