package models

import (
	"log"

	"github.com/shem958/cycle-backend/config"
)

// Migrate runs AutoMigrate on all models
func Migrate() {
	err := config.DB.AutoMigrate(
		&User{},
		&Cycle{},
		// Add more models here as your app grows
		&Post{},
		&Comment{},
		&Report{},
	)
	if err != nil {
		log.Fatalf("❌ Failed to migrate database models: %v", err)
	}

	log.Println("✅ Database models migrated successfully")
}
