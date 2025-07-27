package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// ReactionInput defines the structure for creating or updating a reaction
type ReactionInput struct {
	TargetID   string `json:"target_id" binding:"required"`   // UUID as string
	TargetType string `json:"target_type" binding:"required"` // "post" or "comment"
	Type       string `json:"type" binding:"required"`        // "like" or "dislike"
}

// ReactToContent handles adding or updating a reaction
func ReactToContent(c *gin.Context) {
	var input ReactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user UUID"})
		return
	}

	targetUUID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target UUID"})
		return
	}

	db := config.DB

	var existing models.Reaction
	err = db.Where("user_id = ? AND target_id = ? AND target_type = ?", userUUID, targetUUID, input.TargetType).First(&existing).Error

	if err == nil {
		// Reaction exists – update it
		existing.Type = input.Type
		db.Save(&existing)
		c.JSON(http.StatusOK, gin.H{"message": "Reaction updated"})
		return
	}

	// Create new reaction
	reaction := models.Reaction{
		UserID:     userUUID,
		TargetID:   targetUUID,
		TargetType: input.TargetType,
		Type:       input.Type,
	}

	if err := db.Create(&reaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reaction added"})
}

// RemoveReaction deletes a user's reaction to a post or comment
func RemoveReaction(c *gin.Context) {
	var input ReactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user UUID"})
		return
	}

	targetUUID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target UUID"})
		return
	}

	db := config.DB

	if err := db.Where("user_id = ? AND target_id = ? AND target_type = ?", userUUID, targetUUID, input.TargetType).Delete(&models.Reaction{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove reaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
}
