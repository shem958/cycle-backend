package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
)

func RegisterAppointmentRoutes(router *gin.Engine) {
	appointments := router.Group("/appointments")
	{
		appointments.POST("/", controllers.CreateAppointment)
		appointments.GET("/user/:id", controllers.GetAppointmentsForUser)
	}
}
