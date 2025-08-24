package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// CreatePostpartumLog handles adding a new postpartum entry
func CreatePostpartumLog(c *gin.Context) {
	var input struct {
		UserID            string  `json:"user_id"`
		Date              string  `json:"date"`
		Mood              string  `json:"mood"`
		PainLevel         int     `json:"pain_level"`
		Notes             string  `json:"notes"`
		Breastfeeding     bool    `json:"breastfeeding"`
		SleepHours        float64 `json:"sleep_hours"`
		AppetiteLevel     string  `json:"appetite_level"`
		FollowUpScheduled bool    `json:"follow_up_scheduled"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userID, err := uuid.Parse(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (use YYYY-MM-DD)"})
		return
	}

	log := models.PostpartumLog{
		ID:                uuid.New(),
		UserID:            userID,
		Date:              date,
		Mood:              input.Mood,
		PainLevel:         input.PainLevel,
		Notes:             input.Notes,
		Breastfeeding:     input.Breastfeeding,
		SleepHours:        input.SleepHours,
		AppetiteLevel:     input.AppetiteLevel,
		FollowUpScheduled: input.FollowUpScheduled,
	}

	if err := config.DB.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save postpartum log"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Postpartum log created"})
}

// GetPostpartumLogs retrieves logs for a user
func GetPostpartumLogs(c *gin.Context) {
	userID := c.Param("id")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var logs []models.PostpartumLog
	if err := config.DB.Where("user_id = ?", parsedID).Order("date desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetPostpartumDashboard retrieves combined postpartum data for the dashboard
func GetPostpartumDashboard(c *gin.Context) {
	userID := c.Param("id")
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get logs
	var logs []models.PostpartumLog
	if err := config.DB.Where("user_id = ?", parsedID).Order("date desc").Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}

	// Get checkups
	var checkups []models.PostpartumCheckup
	if err := config.DB.Where("user_id = ?", parsedID).Order("visit_date desc").Find(&checkups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve checkups"})
		return
	}

	// Construct dashboard response
	dashboard := gin.H{
		"logs":     logs,
		"checkups": checkups,
	}

	// Add latest metrics if available
	if len(logs) > 0 {
		latestLog := logs[0]
		dashboard["latestMetrics"] = gin.H{
			"mood":                latestLog.Mood,
			"pain_level":          latestLog.PainLevel,
			"sleep_hours":         latestLog.SleepHours,
			"breastfeeding":       latestLog.Breastfeeding,
			"appetite_level":      latestLog.AppetiteLevel,
			"follow_up_scheduled": latestLog.FollowUpScheduled,
			"date":                latestLog.Date,
			"notes":               latestLog.Notes,
		}
	}

	c.JSON(http.StatusOK, dashboard)
}
