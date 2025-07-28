package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

func CreateAppointment(c *gin.Context) {
	var input struct {
		UserID      string `json:"user_id" binding:"required"`
		DoctorID    string `json:"doctor_id" binding:"required"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Location    string `json:"location"`
		ScheduledAt string `json:"scheduled_at"` // "2025-07-01T15:04:00"
		IsFollowUp  bool   `json:"is_follow_up"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userUUID, err := uuid.Parse(input.UserID)
	doctorUUID, err2 := uuid.Parse(input.DoctorID)
	if err != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user or doctor ID"})
		return
	}

	scheduledAt, err := time.Parse(time.RFC3339, input.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid datetime format. Use ISO 8601"})
		return
	}

	appt := models.Appointment{
		ID:          uuid.New(),
		UserID:      userUUID,
		DoctorID:    doctorUUID,
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		ScheduledAt: scheduledAt,
		IsFollowUp:  input.IsFollowUp,
	}

	if err := config.DB.Create(&appt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Appointment scheduled", "appointment": appt})
}

func GetAppointmentsForUser(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var appts []models.Appointment
	if err := config.DB.Where("user_id = ?", userUUID).Order("scheduled_at asc").Find(&appts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve appointments"})
		return
	}

	c.JSON(http.StatusOK, appts)
}
