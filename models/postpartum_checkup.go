package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostpartumCheckup represents a recovery visit/checkup after childbirth
type PostpartumCheckup struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	DoctorID uuid.UUID `gorm:"type:uuid;index" json:"doctor_id,omitempty"`

	VisitDate         time.Time `gorm:"not null" json:"visit_date"`
	MotherHealthNotes string    `gorm:"type:text" json:"mother_health_notes,omitempty"`
	BabyHealthNotes   string    `gorm:"type:text" json:"baby_health_notes,omitempty"`
	Complications     string    `gorm:"type:text" json:"complications,omitempty"`
	MentalHealth      string    `gorm:"type:text" json:"mental_health,omitempty"`
	NextCheckupAt     time.Time `json:"next_checkup_at,omitempty"`

	// File attachments (e.g. prescriptions, scans, reports)
	Attachments []PostpartumCheckupFile `gorm:"foreignKey:CheckupID;constraint:OnDelete:CASCADE" json:"attachments,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
