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

	api := router.Group("/api")
	RegisterAuthRoutes(api)
	RegisterCycleRoutes(api)
	RegisterUserRoutes(api)
	RegisterCommunityRoutes(api)
	RegisterProfileRoutes(api)
	RegisterModerationRoutes(api) // ✅ Added moderation routes

	// 🛡️ Admin-only group
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	admin.GET("/reports", GetAllReports)                   // removed routes.
	admin.PATCH("/reports/:id/status", UpdateReportStatus) // removed routes.

	return router
}
