package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
	"github.com/shem958/cycle-backend/services"
	"gorm.io/gorm"
)

func RegisterPregnancyCheckupRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	// Initialize service & controller
	service := services.NewPregnancyCheckupService(db)
	controller := controllers.NewPregnancyCheckupController(service)

	pregnancyCheckup := rg.Group("/pregnancy-checkups")
	pregnancyCheckup.Use(middleware.AuthMiddleware())

	pregnancyCheckup.POST("/", controller.CreateCheckup)
	pregnancyCheckup.GET("/user/:userID", controller.GetUserCheckups)
	pregnancyCheckup.GET("/:id", controller.GetCheckup)
	pregnancyCheckup.PUT("/:id", controller.UpdateCheckup)
	pregnancyCheckup.DELETE("/:id", controller.DeleteCheckup)
}
