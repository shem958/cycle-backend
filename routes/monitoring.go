package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterMonitoringRoutes(rg *gin.RouterGroup) {
	monitoring := rg.Group("/monitoring")
	monitoring.Use(middleware.AuthMiddleware())

	monitoring.POST("/", controllers.CreateMonitoringRecord)
	monitoring.GET("/", controllers.GetUserMonitoringRecords)
}
