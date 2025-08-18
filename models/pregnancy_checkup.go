package models

import (
	"time"

	"github.com/google/uuid"
)

// PregnancyCheckup represents a doctor's visit/checkup during pregnancy
type PregnancyCheckup struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PregnancyID uuid.UUID `gorm:"type:uuid;not null;index" json:"pregnancy_id"`
	DoctorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"doctor_id"`
	Doctor      User      `gorm:"foreignKey:DoctorID" json:"doctor"`

	Date      time.Time  `gorm:"not null" json:"date"`
	Notes     string     `gorm:"type:text" json:"notes,omitempty"`
	NextVisit *time.Time `json:"next_visit,omitempty"`

	// ✅ Relation to file attachments
	Files []PregnancyCheckupFile `gorm:"foreignKey:CheckupID" json:"files"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
