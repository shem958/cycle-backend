package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// GetAllReports allows an admin to view all content reports
func GetAllReports(c *gin.Context) {
	var reports []models.Report
	if err := config.DB.Order("created_at desc").Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

// UpdateReportStatus allows an admin to update a report’s status (e.g., resolved, ignored)
func UpdateReportStatus(c *gin.Context) {
	var payload struct {
		Status string `json:"status"` // e.g., "resolved", "ignored"
	}

	reportIDParam := c.Param("id")
	reportID, err := uuid.Parse(reportIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	if err := c.BindJSON(&payload); err != nil || payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result := config.DB.Model(&models.Report{}).Where("id = ?", reportID).Update("status", payload.Status)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report status updated"})
}
