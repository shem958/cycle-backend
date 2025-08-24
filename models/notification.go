package models

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents different kinds of notifications
type NotificationType string

const (
	NotificationTypeSystem      NotificationType = "system"
	NotificationTypeAppointment NotificationType = "appointment"
	NotificationTypeReminder    NotificationType = "reminder"
	NotificationTypeAlert       NotificationType = "alert"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID        uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Title     string           `gorm:"type:varchar(255);not null" json:"title"`
	Message   string           `gorm:"type:text;not null" json:"message"`
	Link      string           `gorm:"type:varchar(255)" json:"link,omitempty"`
	Read      bool             `gorm:"default:false" json:"read"`
	CreatedAt time.Time        `gorm:"not null" json:"created_at"`
	ReadAt    *time.Time       `json:"read_at,omitempty"`
}

// TableName specifies the table name for the Notification model
func (Notification) TableName() string {
	return "notifications"
}
