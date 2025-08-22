package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
	"github.com/shem958/cycle-backend/utils"
)

// GetCycles retrieves all cycles for the authenticated user
func GetCycles(c *gin.Context) {
	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}
	var cycles []models.Cycle

	if err := config.DB.Where("user_id = ?", userID).Find(&cycles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cycles"})
		return
	}

	c.JSON(http.StatusOK, cycles)
}

// AddCycle creates a new cycle for the authenticated user
func AddCycle(c *gin.Context) {
	var input models.Cycle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}
	input.UserID = userID

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cycle"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateCycle updates a cycle owned by the authenticated user
func UpdateCycle(c *gin.Context) {
	cycleID := utils.ParseUUIDParamOrAbort(c, "id")
	if cycleID == uuid.Nil {
		return
	}
	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}

	var cycle models.Cycle
	if err := config.DB.Where("id = ? AND user_id = ?", cycleID, userID).First(&cycle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cycle not found or unauthorized"})
		return
	}

	var input models.Cycle
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data", "details": err.Error()})
		return
	}

	// Update allowed fields
	cycle.StartDate = input.StartDate
	cycle.Length = input.Length
	cycle.Mood = input.Mood
	cycle.Symptoms = input.Symptoms

	if err := config.DB.Save(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cycle"})
		return
	}

	c.JSON(http.StatusOK, cycle)
}

// DeleteCycle removes a cycle owned by the authenticated user
func DeleteCycle(c *gin.Context) {
	cycleID := utils.ParseUUIDParamOrAbort(c, "id")
	if cycleID == uuid.Nil {
		return
	}
	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}

	var cycle models.Cycle
	if err := config.DB.Where("id = ? AND user_id = ?", cycleID, userID).First(&cycle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cycle not found or unauthorized"})
		return
	}

	if err := config.DB.Delete(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cycle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cycle deleted successfully"})
}
