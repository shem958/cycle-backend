package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// VerifyDoctor marks a doctor as verified
func VerifyDoctor(c *gin.Context) {
	doctorID := c.Param("id")
	parsedID, err := uuid.Parse(doctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != models.RoleDoctor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a doctor"})
		return
	}

	user.Verified = true
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify doctor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor verified successfully"})
}

// UnverifyDoctor removes the verified status from a doctor
func UnverifyDoctor(c *gin.Context) {
	doctorID := c.Param("id")
	parsedID, err := uuid.Parse(doctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", parsedID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != models.RoleDoctor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a doctor"})
		return
	}

	user.Verified = false
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unverify doctor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor unverified successfully"})
}

// IssueWarning allows an admin to issue a warning to a doctor
func IssueWarning(c *gin.Context) {
	var payload struct {
		DoctorID string `json:"doctor_id" binding:"required"`
		Reason   string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	doctorUUID, err := uuid.Parse(payload.DoctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var doctor models.User
	if err := config.DB.First(&doctor, "id = ?", doctorUUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	if doctor.Role != models.RoleDoctor {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is not a doctor"})
		return
	}

	// TODO: Replace with actual admin ID from auth context
	adminID := uuid.New() // Placeholder

	warning := models.Warning{
		ID:        uuid.New(),
		DoctorID:  doctorUUID,
		AdminID:   adminID,
		Reason:    payload.Reason,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&warning).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to issue warning"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warning issued"})
}

// GetDoctorWarnings lists all warnings for a given doctor
func GetDoctorWarnings(c *gin.Context) {
	doctorID := c.Param("id")
	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	var warnings []models.Warning
	if err := config.DB.Where("doctor_id = ?", doctorUUID).Order("created_at desc").Find(&warnings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warnings"})
		return
	}

	c.JSON(http.StatusOK, warnings)
}

// BanUser permanently bans a user
func BanUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Banned = true
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

func UnbanUser(c *gin.Context) {
	userID := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Banned = false
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unban user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}
