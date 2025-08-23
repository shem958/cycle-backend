package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Helper function to set CORS headers and abort with error
		abortWithCORSError := func(status int, message string) {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.JSON(status, gin.H{"error": message})
			c.Abort()
		}

		uidRaw, exists := c.Get("user_id")
		if !exists {
			abortWithCORSError(http.StatusUnauthorized, "Unauthorized")
			return
		}

		userIDStr, ok := uidRaw.(string)
		if !ok {
			abortWithCORSError(http.StatusInternalServerError, "Invalid user ID type")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			abortWithCORSError(http.StatusBadRequest, "Invalid user ID format")
			return
		}

		var user models.User
		if err := config.DB.First(&user, "id = ?", userID).Error; err != nil || user.Role != "admin" {
			abortWithCORSError(http.StatusForbidden, "Access denied")
			return
		}

		c.Next()
	}
}
