package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
)

func RegisterPregnancyRoutes(router *gin.Engine) {
	pregnancy := router.Group("/pregnancy")
	{
		pregnancy.POST("/", controllers.CreatePregnancy)
		pregnancy.GET("/user/:user_id", controllers.GetPregnanciesByUser)

		pregnancy.POST("/symptom", controllers.LogSymptom)
		pregnancy.GET("/symptom/:pregnancy_id", controllers.GetSymptoms)
	}
}
