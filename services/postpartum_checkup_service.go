package services

import (
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

// CreatePostpartumCheckup creates a new checkup
func CreatePostpartumCheckup(checkup *models.PostpartumCheckup) error {
	return config.DB.Create(checkup).Error
}

// GetPostpartumCheckupsByUser retrieves all checkups for a user
func GetPostpartumCheckupsByUser(userID uuid.UUID) ([]models.PostpartumCheckup, error) {
	var checkups []models.PostpartumCheckup
	err := config.DB.Preload("Attachments").Where("user_id = ?", userID).Find(&checkups).Error
	return checkups, err
}

// GetPostpartumCheckupByID retrieves a single checkup
func GetPostpartumCheckupByID(id uuid.UUID) (*models.PostpartumCheckup, error) {
	var checkup models.PostpartumCheckup
	err := config.DB.Preload("Attachments").First(&checkup, "id = ?", id).Error
	return &checkup, err
}

// UpdatePostpartumCheckup updates a checkup
func UpdatePostpartumCheckup(checkup *models.PostpartumCheckup) error {
	return config.DB.Save(checkup).Error
}

// DeletePostpartumCheckup deletes a checkup
func DeletePostpartumCheckup(id uuid.UUID) error {
	return config.DB.Delete(&models.PostpartumCheckup{}, "id = ?", id).Error
}
