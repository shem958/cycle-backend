package services

import (
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/models"
	"gorm.io/gorm"
)

type PregnancyCheckupService struct {
	DB *gorm.DB
}

func NewPregnancyCheckupService(db *gorm.DB) *PregnancyCheckupService {
	return &PregnancyCheckupService{DB: db}
}

func (s *PregnancyCheckupService) CreateCheckup(checkup *models.PregnancyCheckup) error {
	return s.DB.Create(checkup).Error
}

func (s *PregnancyCheckupService) GetCheckupsByUser(userID uuid.UUID) ([]models.PregnancyCheckup, error) {
	var checkups []models.PregnancyCheckup
	err := s.DB.Where("user_id = ?", userID).Find(&checkups).Error
	return checkups, err
}

func (s *PregnancyCheckupService) GetCheckupByID(id uuid.UUID) (*models.PregnancyCheckup, error) {
	var checkup models.PregnancyCheckup
	err := s.DB.First(&checkup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &checkup, nil
}

func (s *PregnancyCheckupService) UpdateCheckup(checkup *models.PregnancyCheckup) error {
	return s.DB.Save(checkup).Error
}

func (s *PregnancyCheckupService) DeleteCheckup(id uuid.UUID) error {
	return s.DB.Delete(&models.PregnancyCheckup{}, "id = ?", id).Error
}
