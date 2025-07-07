package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/shem958/cycle-backend/models"
)

var DB *gorm.DB

func ConnectDB() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, continuing without it")
	}

	// Create the database connection (using SQLite)
	database, err := gorm.Open(sqlite.Open("cycles.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Auto-migrate all required models (Cycle, User)
	if err := database.AutoMigrate(&models.Cycle{}, &models.User{}); err != nil {
		log.Fatalf("❌ Failed to auto-migrate database schema: %v", err)
	}

	DB = database

	log.Println("✅ Database connected and models migrated successfully")
}
