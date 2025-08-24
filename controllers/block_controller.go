package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// BlockOrUnblockUser handles blocking or unblocking a user
func BlockOrUnblockUser(c *gin.Context) {
	var input struct {
		TargetID string `json:"target_id" binding:"required"`
		Block    bool   `json:"block"`    // true to block, false to unblock
		IsMuted  bool   `json:"is_muted"` // true for mute, false for block
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user ID from context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse target user ID
	targetID, err := uuid.Parse(input.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	// Prevent self-blocking
	if userID == targetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot block yourself"})
		return
	}

	// Check if target user exists
	var targetUser models.User
	if err := config.DB.First(&targetUser, targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	if input.Block {
		// Create block
		block := models.Block{
			ID:        uuid.New(),
			UserID:    userID,
			TargetID:  targetID,
			IsMuted:   input.IsMuted,
			CreatedAt: time.Now(),
		}

		// Check if already blocked/muted
		var existingBlock models.Block
		result := config.DB.Where("user_id = ? AND target_id = ? AND is_muted = ?", userID, targetID, input.IsMuted).First(&existingBlock)
		if result.Error == nil {
			action := "blocked"
			if input.IsMuted {
				action = "muted"
			}
			c.JSON(http.StatusConflict, gin.H{"error": "User is already " + action})
			return
		}

		if err := config.DB.Create(&block).Error; err != nil {
			action := "block"
			if input.IsMuted {
				action = "mute"
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to " + action + " user"})
			return
		}

		action := "blocked"
		if input.IsMuted {
			action = "muted"
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User " + action + " successfully"})
	} else {
		// Remove block or mute
		result := config.DB.Where("user_id = ? AND target_id = ? AND is_muted = ?", userID, targetID, input.IsMuted).Delete(&models.Block{})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unblock user"})
			return
		}

		if result.RowsAffected == 0 {
			action := "blocked"
			if input.IsMuted {
				action = "muted"
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "User was not " + action})
			return
		}

		action := "unblocked"
		if input.IsMuted {
			action = "unmuted"
		}
		c.JSON(http.StatusOK, gin.H{"message": "User " + action + " successfully"})
	}
}

// GetBlockedUsers returns a list of users blocked by the current user
func GetBlockedUsers(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var blocks []models.Block
	if err := config.DB.Where("user_id = ?", userID).Find(&blocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blocked users"})
		return
	}

	// Get user details for blocked users
	type BlockedUserInfo struct {
		ID        uuid.UUID `json:"id"`
		Username  string    `json:"username"`
		BlockedAt time.Time `json:"blocked_at"`
		IsMuted   bool      `json:"is_muted"`
	}

	blockedUsers := make([]BlockedUserInfo, 0, len(blocks))
	for _, block := range blocks {
		var user models.User
		if err := config.DB.Select("id", "username").First(&user, block.TargetID).Error; err != nil {
			continue // Skip if user not found
		}

		blockedUsers = append(blockedUsers, BlockedUserInfo{
			ID:        user.ID,
			Username:  user.Username,
			BlockedAt: block.CreatedAt,
			IsMuted:   block.IsMuted,
		})
	}

	c.JSON(http.StatusOK, blockedUsers)
}
