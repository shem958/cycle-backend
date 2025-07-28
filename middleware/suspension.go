package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// BlockSuspendedMiddleware prevents suspended users from accessing protected routes
func BlockSuspendedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user models.User
		if err := config.DB.First(&user, "id = ?", userID.(uuid.UUID)).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User lookup failed"})
			return
		}

		if user.Suspended {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Your account is suspended"})
			return
		}

		c.Next()
	}
}
