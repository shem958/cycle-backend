package controllers

import (
	"net/http"

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
