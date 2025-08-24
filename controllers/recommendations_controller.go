package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// GetRecommendations retrieves personalized health recommendations for a user
func GetRecommendations(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	parsedID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get current time for filtering active recommendations
	now := time.Now()

	var recommendations []models.Recommendation
	if err := config.DB.Where("user_id = ? AND active = ? AND valid_from <= ? AND (valid_until IS NULL OR valid_until >= ?)",
		parsedID, true, now, now).
		Order("priority DESC, created_at DESC").
		Find(&recommendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations"})
		return
	}

	if len(recommendations) == 0 {
		// If no personalized recommendations, return default ones
		defaultRecommendations := []gin.H{
			{
				"category": "Exercise",
				"advice":   "Consider starting with gentle exercises like walking or prenatal yoga.",
				"priority": 3,
			},
			{
				"category": "Diet",
				"advice":   "Ensure you're getting adequate nutrition with a balanced diet rich in fruits and vegetables.",
				"priority": 4,
			},
			{
				"category": "Mental Health",
				"advice":   "Take time for self-care and relaxation. Don't hesitate to reach out for support when needed.",
				"priority": 5,
			},
			{
				"category": "Rest",
				"advice":   "Aim for 7-9 hours of sleep per night and take short naps during the day if needed.",
				"priority": 3,
			},
		}
		c.JSON(http.StatusOK, defaultRecommendations)
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

// CreateRecommendation handles creating a new recommendation (admin/doctor only)
func CreateRecommendation(c *gin.Context) {
	var input models.Recommendation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values
	input.ID = uuid.New()
	input.Active = true
	input.ValidFrom = time.Now()
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recommendation"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateRecommendation handles updating an existing recommendation
func UpdateRecommendation(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recommendation ID"})
		return
	}

	var recommendation models.Recommendation
	if err := config.DB.First(&recommendation, parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recommendation not found"})
		return
	}

	var input models.Recommendation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = parsedID
	input.UpdatedAt = time.Now()

	if err := config.DB.Save(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recommendation"})
		return
	}

	c.JSON(http.StatusOK, input)
}

// DeleteRecommendation handles deleting a recommendation
func DeleteRecommendation(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recommendation ID"})
		return
	}

	if err := config.DB.Delete(&models.Recommendation{}, parsedID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recommendation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recommendation deleted successfully"})
}
