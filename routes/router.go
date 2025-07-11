package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all route groups
func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")

	RegisterAuthRoutes(api)
	RegisterCycleRoutes(api)
	RegisterUserRoutes(api)

	return router
}
