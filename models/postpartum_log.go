package models

import (
	"time"

	"github.com/google/uuid"
)

type PostpartumLog struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID"`

	Date      time.Time `gorm:"not null"`
	Mood      string    `gorm:"type:varchar(255)"`
	PainLevel int       `gorm:"type:int"` // 0–10 scale
	Notes     string    `gorm:"type:text"`

	Breastfeeding     bool
	SleepHours        float64
	AppetiteLevel     string // "low", "moderate", "good"
	FollowUpScheduled bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
