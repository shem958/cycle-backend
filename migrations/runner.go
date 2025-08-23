package migrations

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations executes all pending migrations
func RunMigrations(db *gorm.DB) error {
	log.Println("🚀 Starting database migrations...")

	// Run the cycles user_id fix migration
	if err := FixCyclesUserID(db); err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

// RunMigrationsWithUserMapping runs migrations with existing user data preservation
func RunMigrationsWithUserMapping(db *gorm.DB, userIDMapping map[int64]string) error {
	log.Println("🚀 Starting database migrations with data preservation...")

	// Run the cycles user_id fix migration with data preservation
	if err := FixCyclesUserIDWithDataPreservation(db, userIDMapping); err != nil {
		log.Printf("❌ Migration failed: %v", err)
		return err
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}
