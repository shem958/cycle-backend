package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/middleware"
)

// RegisterUserRoutes sets up user-related endpoints
func RegisterUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/users")
	user.Use(middleware.AuthMiddleware())

	// Example future routes:
	// user.GET("/me", controllers.GetProfile)
	// user.PUT("/me", controllers.UpdateProfile)
}
