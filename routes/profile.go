package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterProfileRoutes(r *gin.RouterGroup) {
	profile := r.Group("/api/profile")
	profile.Use(middleware.AuthMiddleware())
	{
		profile.GET("/", controllers.GetProfile)
		profile.PUT("/", controllers.UpdateProfile)
	}
}
