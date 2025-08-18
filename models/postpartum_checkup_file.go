package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostpartumCheckupFile stores uploaded files linked to a postpartum checkup
type PostpartumCheckupFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CheckupID uuid.UUID `gorm:"type:uuid;not null;index" json:"checkup_id"`

	FileName string `gorm:"not null" json:"file_name"`
	FileURL  string `gorm:"not null" json:"file_url"`
	FileType string `gorm:"type:varchar(50)" json:"file_type,omitempty"`

	UploadedAt time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
