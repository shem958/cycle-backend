package models

import (
	"time"

	"github.com/google/uuid"
)

// PregnancyCheckupFile stores file attachments linked to a pregnancy checkup
type PregnancyCheckupFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CheckupID uuid.UUID `gorm:"type:uuid;not null;index" json:"checkup_id"`

	FileURL    string    `gorm:"type:text;not null" json:"file_url"`
	FileType   string    `gorm:"type:varchar(50)" json:"file_type"` // e.g., "image/png", "application/pdf"
	UploadedBy uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	// ✅ who uploaded (doctor or user)
	Uploader User `gorm:"foreignKey:UploadedBy" json:"uploader"`

	CreatedAt time.Time `json:"created_at"`
}
