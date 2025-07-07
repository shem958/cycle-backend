package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/routes"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},

		AllowCredentials: true,
	}))

	config.ConnectDB()
	routes.RegisterCycleRoutes(r)
	routes.RegisterAuthRoutes(r)
	routes.RegisterUserRoutes(r)



	r.Run(":8080")
}
