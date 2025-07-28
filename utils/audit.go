package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

func LogAdminAction(adminID, targetID uuid.UUID, action, details string) {
	log := models.AuditLog{
		ID:        uuid.New(),
		AdminID:   adminID,
		Action:    action,
		TargetID:  targetID,
		Details:   details,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&log)
}
