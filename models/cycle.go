package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Cycle represents a user's menstrual cycle entry
type Cycle struct {
	gorm.Model
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid"` // foreign key
	StartDate time.Time `json:"start_date"`
	Length    int       `json:"length"`   // duration in days
	Mood      string    `json:"mood"`     // optional mood description
	Symptoms  string    `json:"symptoms"` // comma-separated or JSON string
}
