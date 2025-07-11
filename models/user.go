package models

import (
	"gorm.io/gorm"
)

// User represents an account in the system
type User struct {
	gorm.Model
	Email    string  `gorm:"unique;not null" json:"email"`
	Password string  `gorm:"not null" json:"password"`        // stored as hashed
	Cycles   []Cycle `gorm:"foreignKey:UserID" json:"cycles"` // one-to-many relationship
}
