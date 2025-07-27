package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shem958/cycle-backend/models" // ✅ This is allowed here
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env file not found. Using system environment variables.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL not set in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get raw DB connection: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	// ✅ Centralize AutoMigration here
	err = db.AutoMigrate(
		&models.User{},
		&models.Cycle{},
		&models.Post{},
		&models.Comment{},
		&models.Report{},
		&models.Reaction{}, // ✅ include your new model
	)
	if err != nil {
		log.Fatalf("❌ AutoMigration failed: %v", err)
	}

	DB = db
	fmt.Println("✅ PostgreSQL database connected successfully")
}
