package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// Add or update a reaction to a post or comment
func AddOrUpdateReaction(c *gin.Context) {
	var input models.Reaction
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	var existing models.Reaction
	err := db.Where("user_id = ? AND target_id = ? AND target_type = ?", input.UserID, input.TargetID, input.TargetType).
		First(&existing).Error

	if err == nil {
		// Update existing reaction
		existing.Type = input.Type
		if err := db.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reaction"})
			return
		}
		c.JSON(http.StatusOK, existing)
		return
	}

	// Create new reaction
	input.ID = uuid.New()
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add reaction"})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// Get all reactions for a target (post or comment)
func GetReactions(c *gin.Context) {
	targetID := c.Param("target_id")
	var reactions []models.Reaction

	db := config.DB
	if err := db.Where("target_id = ?", targetID).Find(&reactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reactions"})
		return
	}

	c.JSON(http.StatusOK, reactions)
}
