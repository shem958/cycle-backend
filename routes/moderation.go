package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/middleware"
	"github.com/shem958/cycle-backend/models"
)

// RegisterModerationRoutes sets up the /moderation endpoints
func RegisterModerationRoutes(r *gin.RouterGroup) {
	moderation := r.Group("/moderation")
	moderation.Use(middleware.AuthMiddleware())

	moderation.POST("/report", ReportContent)
}

// ReportContent handles reporting a post or comment
func ReportContent(c *gin.Context) {
	var body struct {
		PostID    *uuid.UUID `json:"post_id,omitempty"`
		CommentID *uuid.UUID `json:"comment_id,omitempty"`
		Reason    string     `json:"reason"`
	}

	if err := c.BindJSON(&body); err != nil || (body.PostID == nil && body.CommentID == nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report payload"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	report := models.Report{
		ReporterID:      userID.(uuid.UUID),
		TargetPostID:    body.PostID,
		TargetCommentID: body.CommentID,
		Reason:          body.Reason,
	}

	db := config.DB
	if err := db.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Report submitted successfully"})
}
