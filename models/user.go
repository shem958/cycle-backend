package models

import (
	"time"

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
	Suspended bool      `gorm:"default:false" json:"suspended"`
}

// Block represents a user blocking or muting another user
type Block struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"` // the one doing the blocking/muting
	TargetID  uuid.UUID `gorm:"type:uuid;not null"` // the blocked/muted user
	IsMuted   bool      `gorm:"default:false"`      // true = muted, false = blocked
	CreatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
