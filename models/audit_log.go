package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	AdminID   uuid.UUID `gorm:"type:uuid"`
	Action    string    `gorm:"not null"`  // e.g., "verify_doctor", "ban_user", etc.
	TargetID  uuid.UUID `gorm:"type:uuid"` // ID of the user/post/etc. affected
	Details   string
	CreatedAt time.Time
}
