package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/routes"
)

func main() {
	// Load environment variables from .env (if present)
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ .env file not found, using system environment variables")
	}

	// Connect to the database
	config.ConnectDB()

	// Initialize and setup router
	router := routes.SetupRouter()

	// Determine port to run server on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
