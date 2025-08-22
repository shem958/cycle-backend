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

// VerifyDoctor marks a doctor as verified
func VerifyDoctor(c *gin.Context) {
	doctorID := utils.ParseUUIDParamOrAbort(c, "id")
	if doctorID == uuid.Nil {
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", doctorID).Error; err != nil {
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

	// Replace this with real admin ID when available
	adminID := uuid.New()
	utils.LogAdminAction(adminID, user.ID, "verify_doctor", "Doctor verified")

	c.JSON(http.StatusOK, gin.H{"message": "Doctor verified successfully"})
}

// UnverifyDoctor removes the verified status from a doctor
func UnverifyDoctor(c *gin.Context) {
	doctorID := utils.ParseUUIDParamOrAbort(c, "id")
	if doctorID == uuid.Nil {
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", doctorID).Error; err != nil {
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

	adminID := uuid.New()
	utils.LogAdminAction(adminID, user.ID, "unverify_doctor", "Doctor unverified")

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

	doctorUUID := utils.ParseUUIDParamOrAbort(c, "doctor_id")
	if doctorUUID == uuid.Nil {
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

	adminID := uuid.New()
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

	utils.LogAdminAction(adminID, doctor.ID, "issue_warning", payload.Reason)

	c.JSON(http.StatusOK, gin.H{"message": "Warning issued"})
}

// GetDoctorWarnings lists all warnings for a given doctor
func GetDoctorWarnings(c *gin.Context) {
	doctorUUID := utils.ParseUUIDParamOrAbort(c, "id")
	if doctorUUID == uuid.Nil {
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
	userID := utils.ParseUUIDParamOrAbort(c, "id")
	if userID == uuid.Nil {
		return
	}
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

	adminID := uuid.New()
	utils.LogAdminAction(adminID, user.ID, "ban_user", "User banned by admin")

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// UnbanUser lifts a ban on a user
func UnbanUser(c *gin.Context) {
	userID := utils.ParseUUIDParamOrAbort(c, "id")
	if userID == uuid.Nil {
		return
	}
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

	adminID := uuid.New()
	utils.LogAdminAction(adminID, user.ID, "unban_user", "User unbanned by admin")

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}

// GetAdminMetrics returns basic system statistics
func GetAdminMetrics(c *gin.Context) {
	var totalUsers int64
	var verifiedDoctors int64
	var bannedUsers int64
	var totalWarnings int64

	config.DB.Model(&models.User{}).Count(&totalUsers)
	config.DB.Model(&models.User{}).Where("role = ? AND verified = ?", models.RoleDoctor, true).Count(&verifiedDoctors)
	config.DB.Model(&models.User{}).Where("banned = ?", true).Count(&bannedUsers)
	config.DB.Model(&models.Warning{}).Count(&totalWarnings)

	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"verified_doctors": verifiedDoctors,
		"banned_users":     bannedUsers,
		"total_warnings":   totalWarnings,
	})
}

// SearchFilterUsers allows admins to filter users by role, status, or search query
func SearchFilterUsers(c *gin.Context) {
	role := c.Query("role")
	search := c.Query("search")
	verified := c.Query("verified")
	suspended := c.Query("suspended")
	banned := c.Query("banned")

	var users []models.User
	query := config.DB.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("username ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}
	if verified != "" {
		if verified == "true" {
			query = query.Where("verified = ?", true)
		} else if verified == "false" {
			query = query.Where("verified = ?", false)
		}
	}
	if suspended != "" {
		if suspended == "true" {
			query = query.Where("suspended = ?", true)
		} else if suspended == "false" {
			query = query.Where("suspended = ?", false)
		}
	}
	if banned != "" {
		if banned == "true" {
			query = query.Where("banned = ?", true)
		} else if banned == "false" {
			query = query.Where("banned = ?", false)
		}
	}

	if err := query.Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
