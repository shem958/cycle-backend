package models

import (
	"time"

	"github.com/google/uuid"
)

type MonitoringRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   *time.Time
	Type      string `gorm:"not null"`  // "pregnancy" or "postpartum"
	Data      string `gorm:"type:text"` // JSON-encoded dynamic field data
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
