package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// Migrate DB models
	models.Migrate()

	// Init Gin app
	r := gin.Default()

	// Routes would go here...

	r.Run()
}
