package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shem958/cycle-backend/models"
	"github.com/shem958/cycle-backend/services"
)

// Create a new postpartum checkup
func CreatePostpartumCheckup(c *gin.Context) {
	var checkup models.PostpartumCheckup
	if err := c.ShouldBindJSON(&checkup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreatePostpartumCheckup(&checkup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, checkup)
}

// Get all checkups for a user
func GetPostpartumCheckupsByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	checkups, err := services.GetPostpartumCheckupsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, checkups)
}

// Get a single checkup by ID
func GetPostpartumCheckupByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}
	checkup, err := services.GetPostpartumCheckupByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "checkup not found"})
		return
	}
	c.JSON(http.StatusOK, checkup)
}

// Update a checkup
func UpdatePostpartumCheckup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}

	var updated models.PostpartumCheckup
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated.ID = id

	if err := services.UpdatePostpartumCheckup(&updated); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete a checkup
func DeletePostpartumCheckup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid checkup ID"})
		return
	}
	if err := services.DeletePostpartumCheckup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "checkup deleted"})
}
