package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

// RegisterRecommendationsRoutes sets up the recommendation-related routes
func RegisterRecommendationsRoutes(api *gin.RouterGroup) {
	recommendations := api.Group("/recommendations")
	recommendations.Use(middleware.AuthMiddleware()) // Protect all recommendation routes

	// User routes
	recommendations.GET("", controllers.GetRecommendations) // Get personalized recommendations

	// Routes for creating/updating recommendations (doctors only)
	doctorRoutes := recommendations.Group("")
	doctorRoutes.Use(middleware.DoctorMiddleware())
	{
		doctorRoutes.POST("", controllers.CreateRecommendation)
		doctorRoutes.PUT("/:id", controllers.UpdateRecommendation)
	}

	// Admin-only routes
	adminRoutes := recommendations.Group("")
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		adminRoutes.DELETE("/:id", controllers.DeleteRecommendation)
	}
}
