package database

import (
	"log"
	"os"

	"github.com/vishwakarmaritu/marketplace/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Fatal: DATABASE_URL environment variable is missing. Check your .env file.")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Successfully connected to the database!")

	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Coupon{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schemas: %v", err)
	}

	log.Println("Database schemas migrated successfully!")

	DB = database
}
