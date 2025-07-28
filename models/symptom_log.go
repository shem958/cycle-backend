package models

import (
	"time"

	"github.com/google/uuid"
)

type SymptomLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PregnancyID uuid.UUID `gorm:"type:uuid;not null;index" json:"pregnancy_id"`
	Date        time.Time `gorm:"not null" json:"date"`
	Symptoms    string    `json:"symptoms"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time
}
