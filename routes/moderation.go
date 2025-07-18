package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config" // ✅ ensure this is correct
	"github.com/shem958/cycle-backend/models"
)

func ReportContent(c *gin.Context) {
	var body struct {
		TargetPostID    *uuid.UUID `json:"post_id,omitempty"`
		TargetCommentID *uuid.UUID `json:"comment_id,omitempty"`
		Reason          string     `json:"reason"`
	}

	// Validate request body
	if err := c.BindJSON(&body); err != nil || (body.TargetPostID == nil && body.TargetCommentID == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report payload"})
		return
	}

	// Extract user_id from context and assert to uuid.UUID
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	reporterID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Create report
	report := models.Report{
		ReporterID:      reporterID,
		TargetPostID:    body.TargetPostID,
		TargetCommentID: body.TargetCommentID,
		Reason:          body.Reason,
	}

	db := config.DB // ✅ Replace with your actual DB instance or method
	if err := db.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Report submitted successfully"})
}
