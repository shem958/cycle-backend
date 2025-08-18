package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/services"
)

// RegisterPregnancyCheckupRoutes wires pregnancy checkup endpoints
func RegisterPregnancyCheckupRoutes(api *gin.RouterGroup, db *gorm.DB) {
	// Service + Controller
	service := services.NewPregnancyCheckupService(db)
	controller := &controllers.PregnancyCheckupController{Service: service}

	checkups := api.Group("/pregnancy-checkups")
	{
		checkups.POST("/", controller.CreateCheckup)
		checkups.GET("/user/:userID", controller.GetUserCheckups)
		checkups.GET("/:id", controller.GetCheckup)
		checkups.PUT("/:id", controller.UpdateCheckup)
		checkups.DELETE("/:id", controller.DeleteCheckup)
	}
}
