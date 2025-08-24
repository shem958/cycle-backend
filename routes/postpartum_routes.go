package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

// RegisterPostpartumRoutes registers all postpartum-related routes
func RegisterPostpartumRoutes(api *gin.RouterGroup) {
	postpartum := api.Group("/postpartum")
	postpartum.Use(middleware.AuthMiddleware()) // Protect all postpartum routes

	// Dashboard endpoint - combines logs and checkups
	postpartum.GET("/dashboard/:id", controllers.GetPostpartumDashboard)

	// Postpartum logs
	postpartum.POST("/", controllers.CreatePostpartumLog)
	postpartum.GET("/logs/:id", controllers.GetPostpartumLogs)

	// Postpartum checkups
	checkups := postpartum.Group("/checkups")
	{
		checkups.POST("/", controllers.CreatePostpartumCheckup)
		checkups.GET("/user/:user_id", controllers.GetPostpartumCheckupsByUser)
		checkups.GET("/:id", controllers.GetPostpartumCheckupByID)
		checkups.PUT("/:id", controllers.UpdatePostpartumCheckup)
		checkups.DELETE("/:id", controllers.DeletePostpartumCheckup)
	}
}
