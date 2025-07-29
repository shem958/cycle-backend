package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
	"github.com/shem958/cycle-backend/utils"
)

func CreateMonitoringRecord(c *gin.Context) {
	var input struct {
		Type      string `json:"type" binding:"required"` // "pregnancy" or "postpartum"
		Data      string `json:"data"`
		Notes     string `json:"notes"`
		StartDate string `json:"start_date" binding:"required"`
		EndDate   string `json:"end_date"`
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	start, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	var end *time.Time
	if input.EndDate != "" {
		e, err := time.Parse(time.RFC3339, input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
		end = &e
	}

	// 🔐 Encrypt data and notes
	encryptedData, err := utils.Encrypt(input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt data"})
		return
	}

	encryptedNotes, err := utils.Encrypt(input.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt notes"})
		return
	}

	record := models.MonitoringRecord{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      input.Type,
		Data:      encryptedData,
		Notes:     encryptedNotes,
		StartDate: start,
		EndDate:   end,
	}

	if err := config.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record created"})
}

func GetUserMonitoringRecords(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	var records []models.MonitoringRecord
	if err := config.DB.Where("user_id = ?", userID).Order("start_date desc").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	// 🔓 Decrypt each record
	for i, record := range records {
		decryptedData, err := utils.Decrypt(record.Data)
		if err == nil {
			records[i].Data = decryptedData
		}

		decryptedNotes, err := utils.Decrypt(record.Notes)
		if err == nil {
			records[i].Notes = decryptedNotes
		}
	}

	c.JSON(http.StatusOK, records)
}
