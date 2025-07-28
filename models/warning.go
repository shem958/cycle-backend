// models/warning.go

package models

import (
	"time"

	"github.com/google/uuid"
)

type Warning struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DoctorID  uuid.UUID `gorm:"type:uuid;not null"`
	AdminID   uuid.UUID `gorm:"type:uuid;not null"`
	Reason    string    `gorm:"type:text;not null"`
	CreatedAt time.Time
}
