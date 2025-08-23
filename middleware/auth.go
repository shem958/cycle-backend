package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Helper function to set CORS headers and abort with error
		abortWithCORSError := func(status int, message string) {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.JSON(status, gin.H{"error": message})
			c.Abort()
		}

		auth := c.GetHeader("Authorization")
		if auth == "" {
			abortWithCORSError(http.StatusUnauthorized, "Missing token")
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			abortWithCORSError(http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			abortWithCORSError(http.StatusUnauthorized, "Invalid claims")
			return
		}

		// ✅ Ensure user_id is a string (UUID)
		userID, ok := claims["user_id"].(string)
		if !ok {
			abortWithCORSError(http.StatusUnauthorized, "Invalid user ID format")
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", claims["role"])
		c.Next()
	}
}
