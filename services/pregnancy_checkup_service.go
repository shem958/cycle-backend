package services

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shem958/cycle-backend/models"
)

type PregnancyCheckupService struct {
	DB *gorm.DB
}

func NewPregnancyCheckupService(db *gorm.DB) *PregnancyCheckupService {
	return &PregnancyCheckupService{DB: db}
}

// Create new pregnancy checkup
func (s *PregnancyCheckupService) CreateCheckup(checkup *models.PregnancyCheckup) error {
	return s.DB.Create(checkup).Error
}

// Get all checkups for a user
func (s *PregnancyCheckupService) GetCheckupsByUser(userID uuid.UUID) ([]models.PregnancyCheckup, error) {
	var checkups []models.PregnancyCheckup
	err := s.DB.Where("user_id = ?", userID).Order("visit_date desc").Find(&checkups).Error
	return checkups, err
}

// Get a single checkup
func (s *PregnancyCheckupService) GetCheckupByID(id uuid.UUID) (*models.PregnancyCheckup, error) {
	var checkup models.PregnancyCheckup
	err := s.DB.First(&checkup, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &checkup, nil
}

// Update checkup
func (s *PregnancyCheckupService) UpdateCheckup(checkup *models.PregnancyCheckup) error {
	return s.DB.Save(checkup).Error
}

// Delete checkup
func (s *PregnancyCheckupService) DeleteCheckup(id uuid.UUID) error {
	return s.DB.Delete(&models.PregnancyCheckup{}, "id = ?", id).Error
}
