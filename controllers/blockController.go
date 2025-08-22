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

func BlockOrMuteUser(c *gin.Context) {
	var input struct {
		TargetID uuid.UUID `json:"target_id"`
		IsMuted  bool      `json:"is_muted"` // true = mute, false = block
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}

	block := models.Block{
		ID:        uuid.New(),
		UserID:    userID,
		TargetID:  input.TargetID,
		IsMuted:   input.IsMuted,
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&block).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block/mute user"})
		return
	}

	action := "blocked"
	if input.IsMuted {
		action = "muted"
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User " + action + " successfully"})
}

func UnblockUser(c *gin.Context) {
	targetID := utils.ParseUUIDParamOrAbort(c, "target_id")
	if targetID == uuid.Nil {
		return
	}
	userID := utils.GetUserIDFromContextOrAbort(c)
	if userID == uuid.Nil {
		return
	}

	if err := config.DB.Where("user_id = ? AND target_id = ?", userID, targetID).Delete(&models.Block{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}
