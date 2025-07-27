package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// GetAllReports retrieves all user-submitted reports
func GetAllReports(c *gin.Context) {
	var reports []models.Report
	if err := config.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch reports"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// UpdateReportStatus updates the status of a report (e.g., reviewed, dismissed)
func UpdateReportStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status"` // e.g., "reviewed", "dismissed"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Report{}).
		Where("id = ?", id).
		Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// DeletePost removes a post by ID
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Post{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// DeleteComment removes a comment by ID
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Comment{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// SuspendUser sets a user's account status to suspended
func SuspendUser(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("suspended", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User suspended"})
}
