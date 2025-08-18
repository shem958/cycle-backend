package services

import (
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/config"
	"github.com/shem958/cycle-backend/models"
)

type PregnancyCheckupService struct{}

func (s *PregnancyCheckupService) CreateCheckup(checkup *models.PregnancyCheckup) error {
	db := config.GetDB()
	return db.Create(checkup).Error
}

func (s *PregnancyCheckupService) GetCheckupsByUser(userID uuid.UUID) ([]models.PregnancyCheckup, error) {
	db := config.GetDB()
	var checkups []models.PregnancyCheckup
	err := db.Where("user_id = ?", userID).Find(&checkups).Error
	return checkups, err
}

func (s *PregnancyCheckupService) GetCheckupByID(id uuid.UUID) (*models.PregnancyCheckup, error) {
	db := config.GetDB()
	var checkup models.PregnancyCheckup
	if err := db.First(&checkup, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &checkup, nil
}

func (s *PregnancyCheckupService) UpdateCheckup(checkup *models.PregnancyCheckup) error {
	db := config.GetDB()
	return db.Save(checkup).Error
}

func (s *PregnancyCheckupService) DeleteCheckup(id uuid.UUID) error {
	db := config.GetDB()
	return db.Delete(&models.PregnancyCheckup{}, "id = ?", id).Error
}
