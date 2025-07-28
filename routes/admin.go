package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/models"
)

func RegisterAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.GET("/reports", GetAllReports)
		admin.PUT("/reports/:id", UpdateReportStatus)

		admin.DELETE("/posts/:id", DeletePost)
		admin.DELETE("/comments/:id", DeleteComment)

		admin.PUT("/users/:id/suspend", SuspendUser)

		admin.PUT("/verify-doctor/:id", controllers.VerifyDoctor)
		admin.PUT("/unverify-doctor/:id", controllers.UnverifyDoctor)

		admin.POST("/warnings", controllers.IssueWarning)
		admin.GET("/warnings/:id", controllers.GetDoctorWarnings)
	}
}

// Get all reports
func GetAllReports(c *gin.Context) {
	var reports []models.Report
	config.DB.Order("created_at desc").Find(&reports)
	c.JSON(http.StatusOK, reports)
}

// Update a report's status
func UpdateReportStatus(c *gin.Context) {
	var payload struct {
		Status string `json:"status"` // "resolved", "ignored"
	}
	reportID := c.Param("id")

	if err := c.BindJSON(&payload); err != nil || payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result := config.DB.Model(&models.Report{}).Where("id = ?", reportID).Update("status", payload.Status)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// Delete a post by ID
func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	result := config.DB.Delete(&models.Post{}, "id = ?", postID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// Delete a comment by ID
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	result := config.DB.Delete(&models.Comment{}, "id = ?", commentID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}

// Suspend a user by ID
func SuspendUser(c *gin.Context) {
	userID := c.Param("id")
	result := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("suspended", true)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User suspended"})
}
