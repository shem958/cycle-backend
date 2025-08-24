package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Helper function to set CORS headers and abort with error
func abortWithCORSError(c *gin.Context, status int, message string) {
	c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	c.JSON(status, gin.H{"error": message})
	c.Abort()
}

// verifyUserRole checks if the user has one of the allowed roles
func verifyUserRole(c *gin.Context, allowedRoles []string) bool {
	userRole, exists := c.Get("user_role")
	if !exists {
		abortWithCORSError(c, http.StatusUnauthorized, "Role not found")
		return false
	}

	roleStr, ok := userRole.(string)
	if !ok {
		abortWithCORSError(c, http.StatusInternalServerError, "Invalid role type")
		return false
	}

	for _, role := range allowedRoles {
		if roleStr == role {
			return true
		}
	}

	abortWithCORSError(c, http.StatusForbidden, "Access denied")
	return false
}

// AdminMiddleware ensures the user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !verifyUserRole(c, []string{"admin"}) {
			return
		}
		c.Next()
	}
}

// DoctorMiddleware ensures the user has doctor role
func DoctorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !verifyUserRole(c, []string{"doctor"}) {
			return
		}
		c.Next()
	}
}

// AdminOrDoctorMiddleware ensures the user has either admin or doctor role
func AdminOrDoctorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !verifyUserRole(c, []string{"admin", "doctor"}) {
			return
		}
		c.Next()
	}
}

// UserMiddleware ensures the user has a valid role
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !verifyUserRole(c, []string{"user", "admin", "doctor"}) {
			return
		}
		c.Next()
	}
}
