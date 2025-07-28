package models

import (
	"time"

	"github.com/google/uuid"
)

type Pregnancy struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	CurrentWeek int       `json:"current_week"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
