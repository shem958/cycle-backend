package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null"`
	User   User      `gorm:"foreignKey:UserID"`

	DoctorID uuid.UUID `gorm:"type:uuid;not null"`
	Doctor   User      `gorm:"foreignKey:DoctorID"`

	Title       string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:text"`
	Location    string    `gorm:"type:varchar(255)"`
	ScheduledAt time.Time `gorm:"not null"`

	IsFollowUp bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
