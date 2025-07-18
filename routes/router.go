package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/middleware"
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

	// ✅ Public API groups
	api := router.Group("/api")
	RegisterAuthRoutes(api)
	RegisterCycleRoutes(api)
	RegisterUserRoutes(api)
	RegisterCommunityRoutes(api)
	RegisterProfileRoutes(api)
	RegisterModerationRoutes(api) // if applicable

	// ✅ Admin-only routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	// Report moderation
	admin.GET("/reports", GetAllReports)
	admin.PATCH("/reports/:id/status", UpdateReportStatus)

	// Content & user management
	admin.DELETE("/posts/:id", DeletePost)
	admin.DELETE("/comments/:id", DeleteComment)
	admin.PUT("/users/:id/suspend", SuspendUser)

	return router
}
