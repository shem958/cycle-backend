package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Post represents a community discussion post
type Post struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	AuthorID    uuid.UUID      `gorm:"type:uuid;not null"`
	Author      User           `gorm:"foreignKey:AuthorID"`
	Title       string         `gorm:"not null"`
	Content     string         `gorm:"type:text;not null"`
	Tags        pq.StringArray `gorm:"type:text[]"`
	IsAnonymous bool           `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Comments    []Comment `gorm:"foreignKey:PostID"`
}

// Comment on a post
// Comment on a post
type Comment struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PostID      uuid.UUID  `gorm:"type:uuid;not null"`
	AuthorID    uuid.UUID  `gorm:"type:uuid;not null"`
	Post        Post       `gorm:"foreignKey:PostID"`
	Author      User       `gorm:"foreignKey:AuthorID"`
	Content     string     `gorm:"type:text;not null"`
	IsAnonymous bool       `gorm:"default:false"`
	ParentID    *uuid.UUID `gorm:"type:uuid"`           // ✅ Parent comment (null = top-level)
	Replies     []Comment  `gorm:"foreignKey:ParentID"` // ✅ Children (nested replies)
	CreatedAt   time.Time
}

// Report content (either post or comment)
type Report struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReporterID      uuid.UUID  `gorm:"type:uuid;not null"`
	TargetPostID    *uuid.UUID `gorm:"type:uuid"`
	TargetCommentID *uuid.UUID `gorm:"type:uuid"`
	Reason          string     `gorm:"type:text;not null"`
	CreatedAt       time.Time
}

// Reaction to a post (e.g., 👍, ❤️)
type Reaction struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TargetID   uuid.UUID `gorm:"type:uuid;not null" json:"target_id"`          // post or comment ID
	TargetType string    `gorm:"type:varchar(10);not null" json:"target_type"` // "post" or "comment"
	Type       string    `gorm:"type:varchar(10);not null" json:"type"`        // "like" or "dislike"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ReactionInput struct {
	TargetID   string `json:"target_id" binding:"required"`   // incoming as string
	TargetType string `json:"target_type" binding:"required"` // "post" or "comment"
	Type       string `json:"type" binding:"required"`        // "like" or "dislike"
}
