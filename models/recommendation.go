package models

import (
	"time"

	"github.com/google/uuid"
)

// Recommendation represents health recommendations for users
type Recommendation struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Category   string     `gorm:"not null" json:"category"` // e.g., "Exercise", "Diet", "Mental Health"
	Advice     string     `gorm:"type:text;not null" json:"advice"`
	Source     string     `json:"source,omitempty"`           // optional source/reference
	Priority   int        `gorm:"default:1" json:"priority"`  // 1 (low) to 5 (high)
	Active     bool       `gorm:"default:true" json:"active"` // whether recommendation is currently active
	ValidFrom  time.Time  `gorm:"not null" json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until,omitempty"` // optional expiration date
	CreatedAt  time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"not null" json:"updated_at"`
}
