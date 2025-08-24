package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

// RegisterNotificationsRoutes sets up notification-related routes
func RegisterNotificationsRoutes(api *gin.RouterGroup) {
	notifications := api.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware()) // Protect all notification routes

	// User routes
	notifications.GET("", controllers.GetNotifications)                  // Get all notifications
	notifications.PUT("/:id/read", controllers.MarkNotificationRead)     // Mark single notification as read
	notifications.PUT("/read-all", controllers.MarkAllNotificationsRead) // Mark all notifications as read
	notifications.DELETE("/:id", controllers.DeleteNotification)         // Delete a notification

	// Admin routes
	adminRoutes := notifications.Group("")
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		adminRoutes.POST("", controllers.CreateNotification) // Create a new notification (admin only)
	}
}
