package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterAnalyticsRoutes(rg *gin.RouterGroup) {
	gr := rg.Group("/analytics")
	gr.Use(middleware.AuthMiddleware())

	// User: self analytics (JSON + CSV)
	gr.GET("/user/:user_id/pregnancy-postpartum", controllers.GetPregnancyPostpartumAnalytics)
	gr.GET("/user/:user_id/pregnancy-postpartum.csv", controllers.ExportPregnancyPostpartumCSV)

	// Doctor/Admin: view patient analytics (JSON + CSV)
	gr.GET("/doctor/patient/:patient_id/pregnancy-postpartum", controllers.GetPatientAnalyticsForDoctor)
	gr.GET("/doctor/patient/:patient_id/pregnancy-postpartum.csv", controllers.ExportPatientAnalyticsCSVForDoctor)
}
