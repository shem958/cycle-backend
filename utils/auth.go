package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromContextOrAbort extracts the user_id from context, parses it as UUID, or aborts with 400 if invalid.
func GetUserIDFromContextOrAbort(c *gin.Context) uuid.UUID {
	userIDStr, ok := c.MustGet("user_id").(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context or not a string"})
		c.Abort()
		return uuid.Nil
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		c.Abort()
		return uuid.Nil
	}
	return userID
}

// ParseUUIDParamOrAbort extracts a URL param, parses it as UUID, or aborts with 400 if invalid.
func ParseUUIDParamOrAbort(c *gin.Context, param string) uuid.UUID {
	idStr := c.Param(param)
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for param: " + param})
		c.Abort()
		return uuid.Nil
	}
	return id
}
