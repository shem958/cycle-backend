package migrations

import (
	"log"

	"gorm.io/gorm"
)

// FixCyclesUserID handles the migration from bigint user_id to uuid user_id
func FixCyclesUserID(db *gorm.DB) error {
	log.Println("🔄 Starting migration: Fix cycles user_id column type...")

	// Check if cycles table exists
	if !db.Migrator().HasTable("cycles") {
		log.Println("✅ cycles table doesn't exist, skipping migration")
		return nil
	}

	// Check if user_id column exists and its type
	if !db.Migrator().HasColumn("cycles", "user_id") {
		log.Println("✅ user_id column doesn't exist in cycles table, skipping migration")
		return nil
	}

	// Step 1: Check if there's any data in cycles table
	var count int64
	if err := db.Table("cycles").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Printf("⚠️  Found %d records in cycles table. Manual data migration required.", count)
		log.Println("Please backup your data before proceeding with this migration.")

		// Option 1: If you want to preserve data, you'll need to:
		// 1. Create a mapping between old bigint user IDs and new UUID user IDs
		// 2. Update the cycles table accordingly
		// For now, we'll show an error and let the user decide

		log.Println("❌ Migration stopped. You have existing data that needs to be handled manually.")
		log.Println("Options:")
		log.Println("1. If data is not important: DROP TABLE cycles; and restart the app")
		log.Println("2. If data is important: Create a mapping strategy for user IDs")
		return nil
	}

	// Step 2: If no data exists, we can safely drop and recreate the column
	log.Println("🗑️  No data found in cycles table. Dropping and recreating user_id column...")

	// Drop the old user_id column
	if err := db.Migrator().DropColumn("cycles", "user_id"); err != nil {
		log.Printf("❌ Failed to drop user_id column: %v", err)
		return err
	}

	// Add the new user_id column with UUID type
	if err := db.Exec("ALTER TABLE cycles ADD COLUMN user_id UUID").Error; err != nil {
		log.Printf("❌ Failed to add new user_id column: %v", err)
		return err
	}

	log.Println("✅ Successfully migrated cycles.user_id from bigint to uuid")
	return nil
}

// FixCyclesUserIDWithDataPreservation handles migration while preserving existing data
// This assumes you have a way to map old bigint user IDs to new UUID user IDs
func FixCyclesUserIDWithDataPreservation(db *gorm.DB, userIDMapping map[int64]string) error {
	log.Println("🔄 Starting migration: Fix cycles user_id with data preservation...")

	// Check if cycles table exists
	if !db.Migrator().HasTable("cycles") {
		log.Println("✅ cycles table doesn't exist, skipping migration")
		return nil
	}

	// Step 1: Add temporary UUID column
	if err := db.Exec("ALTER TABLE cycles ADD COLUMN user_id_temp UUID").Error; err != nil {
		log.Printf("❌ Failed to add temporary user_id_temp column: %v", err)
		return err
	}

	// Step 2: Populate the temporary column with mapped UUIDs
	for oldID, newUUID := range userIDMapping {
		if err := db.Exec("UPDATE cycles SET user_id_temp = ? WHERE user_id = ?", newUUID, oldID).Error; err != nil {
			log.Printf("❌ Failed to update user_id_temp for old ID %d: %v", oldID, err)
			return err
		}
	}

	// Step 3: Drop the old user_id column
	if err := db.Migrator().DropColumn("cycles", "user_id"); err != nil {
		log.Printf("❌ Failed to drop old user_id column: %v", err)
		return err
	}

	// Step 4: Rename the temporary column to user_id
	if err := db.Exec("ALTER TABLE cycles RENAME COLUMN user_id_temp TO user_id").Error; err != nil {
		log.Printf("❌ Failed to rename user_id_temp to user_id: %v", err)
		return err
	}

	log.Println("✅ Successfully migrated cycles.user_id with data preservation")
	return nil
}
