package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterCycleRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Apply authentication to all /api/cycles routes
	protected := api.Group("/cycles")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", controllers.GetCycles)
		protected.POST("", controllers.AddCycle)
		protected.PUT("/:id", controllers.UpdateCycle)
		protected.DELETE("/:id", controllers.DeleteCycle)
	}
}
