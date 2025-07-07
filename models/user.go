package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
