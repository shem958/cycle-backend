package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PregnancyCheckup represents a single prenatal visit/checkup
type PregnancyCheckup struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DoctorID uuid.UUID `gorm:"type:uuid;index" json:"doctor_id,omitempty"`

	VisitDate     time.Time `gorm:"not null" json:"visit_date"`
	DoctorNotes   string    `gorm:"type:text" json:"doctor_notes,omitempty"`
	Weight        float64   `gorm:"type:decimal(5,2)" json:"weight,omitempty"`
	BloodPressure string    `gorm:"type:varchar(20)" json:"blood_pressure,omitempty"`
	NextCheckupAt time.Time `json:"next_checkup_at,omitempty"`

	// ✅ Relation to attachments
	Attachments []PregnancyCheckupFile `gorm:"foreignKey:CheckupID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
