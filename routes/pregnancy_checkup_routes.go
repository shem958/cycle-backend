package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
	"github.com/shem958/cycle-backend/services"
)

func RegisterPregnancyCheckupRoutes(rg *gin.RouterGroup) {
	pregnancyCheckup := rg.Group("/pregnancy-checkups")
	pregnancyCheckup.Use(middleware.AuthMiddleware())

	// create service + controller
	service := &services.PregnancyCheckupService{}
	controller := &controllers.PregnancyCheckupController{Service: service}

	pregnancyCheckup.POST("/", controller.CreateCheckup)
	pregnancyCheckup.GET("/user/:userID", controller.GetUserCheckups)
	pregnancyCheckup.GET("/:id", controller.GetCheckup)
	pregnancyCheckup.PUT("/:id", controller.UpdateCheckup)
	pregnancyCheckup.DELETE("/:id", controller.DeleteCheckup)
}
