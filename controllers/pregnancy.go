package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// CreatePregnancy starts a new pregnancy record
func CreatePregnancy(c *gin.Context) {
	var payload struct {
		UserID    string    `json:"user_id" binding:"required"`
		StartDate time.Time `json:"start_date" binding:"required"`
		Notes     string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pregnancy := models.Pregnancy{
		ID:          uuid.New(),
		UserID:      userUUID,
		StartDate:   payload.StartDate,
		CurrentWeek: int(time.Since(payload.StartDate).Hours() / (24 * 7)),
		Notes:       payload.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := config.DB.Create(&pregnancy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pregnancy"})
		return
	}

	c.JSON(http.StatusCreated, pregnancy)
}

// GetPregnanciesByUser retrieves pregnancy records for a user
func GetPregnanciesByUser(c *gin.Context) {
	userID := c.Param("user_id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var pregnancies []models.Pregnancy
	if err := config.DB.Where("user_id = ?", userUUID).Order("created_at desc").Find(&pregnancies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pregnancies"})
		return
	}

	c.JSON(http.StatusOK, pregnancies)
}

// LogSymptom allows a user to log a symptom during pregnancy
func LogSymptom(c *gin.Context) {
	var payload struct {
		UserID      string    `json:"user_id" binding:"required"`
		PregnancyID string    `json:"pregnancy_id" binding:"required"`
		Date        time.Time `json:"date" binding:"required"`
		Symptoms    string    `json:"symptoms" binding:"required"`
		Notes       string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userUUID, err := uuid.Parse(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	pregnancyUUID, err := uuid.Parse(payload.PregnancyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pregnancy ID"})
		return
	}

	symptom := models.SymptomLog{
		ID:          uuid.New(),
		UserID:      userUUID,
		PregnancyID: pregnancyUUID,
		Date:        payload.Date,
		Symptoms:    payload.Symptoms,
		Notes:       payload.Notes,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&symptom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log symptom"})
		return
	}

	c.JSON(http.StatusCreated, symptom)
}

// GetSymptoms retrieves all symptom logs for a pregnancy
func GetSymptoms(c *gin.Context) {
	pregnancyID := c.Param("pregnancy_id")
	pregnancyUUID, err := uuid.Parse(pregnancyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pregnancy ID"})
		return
	}

	var symptoms []models.SymptomLog
	if err := config.DB.Where("pregnancy_id = ?", pregnancyUUID).Order("date desc").Find(&symptoms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve symptoms"})
		return
	}

	c.JSON(http.StatusOK, symptoms)
}
