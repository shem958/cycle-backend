package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	RoleUser   = "user"
	RoleDoctor = "doctor"
	RoleMod    = "moderator"
	RoleAdmin  = "admin"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"default:user" json:"role"` // "user", "doctor", "moderator", etc.
	Bio       string    `json:"bio,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
