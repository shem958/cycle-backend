package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shem958/cycle-backend/models"
	"github.com/shem958/cycle-backend/services"
	"gorm.io/gorm"
)

type PregnancyCheckupController struct {
	Service *services.PregnancyCheckupService
}

// Create new checkup
func (pc *PregnancyCheckupController) CreateCheckup(c *gin.Context) {
	var input struct {
		UserID        uuid.UUID `json:"user_id" binding:"required"`
		DoctorID      uuid.UUID `json:"doctor_id"`
		VisitDate     time.Time `json:"visit_date" binding:"required"`
		DoctorNotes   string    `json:"doctor_notes"`
		Weight        float64   `json:"weight"`
		BloodPressure string    `json:"blood_pressure"`
		NextCheckupAt time.Time `json:"next_checkup_at"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkup := models.PregnancyCheckup{
		UserID:        input.UserID,
		DoctorID:      input.DoctorID,
		VisitDate:     input.VisitDate,
		DoctorNotes:   input.DoctorNotes,
		Weight:        input.Weight,
		BloodPressure: input.BloodPressure,
		NextCheckupAt: input.NextCheckupAt,
	}

	if err := pc.Service.CreateCheckup(&checkup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save checkup"})
		return
	}

	c.JSON(http.StatusCreated, checkup)
}

// Get checkups for a user
func (pc *PregnancyCheckupController) GetUserCheckups(c *gin.Context) {
	userID := c.Param("userID")
	uid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	checkups, err := pc.Service.GetCheckupsByUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch checkups"})
		return
	}

	c.JSON(http.StatusOK, checkups)
}

// Get a single checkup
func (pc *PregnancyCheckupController) GetCheckup(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}

	checkup, err := pc.Service.GetCheckupByID(cid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "checkup not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch checkup"})
		}
		return
	}

	c.JSON(http.StatusOK, checkup)
}

// Update checkup
func (pc *PregnancyCheckupController) UpdateCheckup(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}

	checkup, err := pc.Service.GetCheckupByID(cid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "checkup not found"})
		return
	}

	if err := c.ShouldBindJSON(checkup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.Service.UpdateCheckup(checkup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update checkup"})
		return
	}

	c.JSON(http.StatusOK, checkup)
}

// Delete checkup
func (pc *PregnancyCheckupController) DeleteCheckup(c *gin.Context) {
	id := c.Param("id")
	cid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}

	if err := pc.Service.DeleteCheckup(cid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete checkup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "checkup deleted"})
}
