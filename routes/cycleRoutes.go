package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

// RegisterCycleRoutes sets up cycle endpoints
func RegisterCycleRoutes(rg *gin.RouterGroup) {
	cycle := rg.Group("/cycles")
	cycle.Use(middleware.AuthMiddleware())

	cycle.GET("", controllers.GetCycles)
	cycle.GET("/", controllers.GetCycles)
	cycle.POST("", controllers.AddCycle)
	cycle.POST("/", controllers.AddCycle)
	cycle.PUT("/:id", controllers.UpdateCycle)
	cycle.DELETE("/:id", controllers.DeleteCycle)
}
