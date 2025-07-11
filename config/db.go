package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the PostgreSQL database connection
func ConnectDB() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found. Using system environment variables.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL not set in environment")
	}

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	// Check DB connection with a simple ping query
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get raw DB connection: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	// Assign DB connection to global variable
	DB = db

	fmt.Println("✅ PostgreSQL database connected successfully")
}
