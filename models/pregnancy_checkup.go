package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PregnancyCheckup struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	DoctorID      uuid.UUID `gorm:"type:uuid" json:"doctor_id"`
	VisitDate     time.Time `gorm:"not null" json:"visit_date"`
	DoctorNotes   string    `gorm:"type:text" json:"doctor_notes"`
	Weight        float64   `json:"weight"`
	BloodPressure string    `json:"blood_pressure"`
	NextCheckupAt time.Time `json:"next_checkup_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (pc *PregnancyCheckup) BeforeCreate(tx *gorm.DB) (err error) {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return
}
