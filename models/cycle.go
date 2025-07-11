package models

import (
	"time"

	"gorm.io/gorm"
)

// Cycle represents a user's menstrual cycle entry
type Cycle struct {
	gorm.Model
	UserID    uint      `json:"user_id"` // foreign key
	StartDate time.Time `json:"start_date"`
	Length    int       `json:"length"`   // duration in days
	Mood      string    `json:"mood"`     // optional mood description
	Symptoms  string    `json:"symptoms"` // comma-separated or JSON string
}
