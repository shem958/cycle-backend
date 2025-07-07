package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// GetCycles - returns all cycles for the authenticated user
func GetCycles(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var cycles []models.Cycle

	if err := config.DB.Where("user_id = ?", userID).Find(&cycles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cycles"})
		return
	}

	c.JSON(http.StatusOK, cycles)
}

// AddCycle - creates a new cycle for the authenticated user
func AddCycle(c *gin.Context) {
	var cycle models.Cycle

	if err := c.ShouldBindJSON(&cycle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)
	cycle.UserID = userID // 👈 Assign ownership

	if err := config.DB.Create(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cycle entry"})
		return
	}

	c.JSON(http.StatusCreated, cycle)
}

// UpdateCycle - updates a cycle only if it belongs to the authenticated user
func UpdateCycle(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")
	var cycle models.Cycle

	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&cycle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cycle not found"})
		return
	}

	var updated models.Cycle
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data", "details": err.Error()})
		return
	}

	cycle.StartDate = updated.StartDate
	cycle.Length = updated.Length
	cycle.Mood = updated.Mood
	cycle.Symptoms = updated.Symptoms

	if err := config.DB.Save(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cycle"})
		return
	}

	c.JSON(http.StatusOK, cycle)
}

// DeleteCycle - deletes a cycle only if it belongs to the authenticated user
func DeleteCycle(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	id := c.Param("id")

	var cycle models.Cycle
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&cycle).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cycle not found"})
		return
	}

	if err := config.DB.Delete(&cycle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cycle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cycle deleted successfully"})
}
