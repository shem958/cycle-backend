package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterInsightsRoutes(rg *gin.RouterGroup) {
	insight := rg.Group("/insights")
	insight.Use(middleware.AuthMiddleware())

	insight.GET("/cycle", controllers.GetCycleInsights)
}
