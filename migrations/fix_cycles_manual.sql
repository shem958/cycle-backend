-- Manual SQL fix for cycles table user_id column type conversion
-- Run this script directly in your PostgreSQL database if you don't have important data to preserve

-- Option 1: Drop and recreate the cycles table (if no important data)
DROP TABLE IF EXISTS cycles CASCADE;

-- The table will be recreated automatically when you restart your Go application
-- with the correct UUID type for user_id

-- Option 2: If you have data to preserve, use this approach instead:
-- (Comment out Option 1 above and uncomment the lines below)

-- Step 1: Create backup of existing data
-- CREATE TABLE cycles_backup AS SELECT * FROM cycles;

-- Step 2: Drop the cycles table
-- DROP TABLE IF EXISTS cycles CASCADE;

-- Step 3: Restart your Go application to recreate the table with correct schema

-- Step 4: If you need to restore data, you'll need to map old bigint user_ids to new UUIDs
-- This requires manual mapping based on your users table
-- Example:
-- INSERT INTO cycles (user_id, start_date, length, mood, symptoms, created_at, updated_at)
-- SELECT u.id, cb.start_date, cb.length, cb.mood, cb.symptoms, cb.created_at, cb.updated_at
-- FROM cycles_backup cb
-- JOIN users u ON u.some_old_reference = cb.user_id;

-- Step 5: Clean up backup table
-- DROP TABLE cycles_backup;

-- Note: Make sure to backup your database before running any of these commands!
