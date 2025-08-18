package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PregnancyCheckupFile stores file attachments linked to a pregnancy checkup
type PregnancyCheckupFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CheckupID uuid.UUID `gorm:"type:uuid;not null;index" json:"checkup_id"`

	FileName   string    `gorm:"not null" json:"file_name"`
	FileURL    string    `gorm:"type:text;not null" json:"file_url"`
	FileType   string    `gorm:"type:varchar(50)" json:"file_type"`
	UploadedBy uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`

	// ✅ relation to uploader (User)
	Uploader User `gorm:"foreignKey:UploadedBy" json:"uploader"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
