package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all route groups
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// ✅ CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	RegisterAuthRoutes(api)
	RegisterCycleRoutes(api)
	RegisterUserRoutes(api)
	RegisterCommunityRoutes(api)
	RegisterProfileRoutes(api)

	return router
}
