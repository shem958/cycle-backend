package models

import (
	"time"

	"github.com/google/uuid"
)

type Pregnancy struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	DueDate     time.Time `json:"due_date"` // calculated (StartDate + 40 weeks)
	CurrentWeek int       `json:"current_week"`
	Status      string    `gorm:"default:'active'" json:"status"` // "active", "completed", "miscarried"
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
