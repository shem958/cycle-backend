package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
)

// RegisterAuthRoutes sets up auth endpoints
func RegisterAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/")
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)
}
