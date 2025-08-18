package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterPostpartumCheckupRoutes(rg *gin.RouterGroup) {
	routes := rg.Group("/postpartum-checkups")
	routes.Use(middleware.AuthMiddleware())

	routes.POST("/", controllers.CreatePostpartumCheckup)
	routes.GET("/user/:user_id", controllers.GetPostpartumCheckupsByUser)
	routes.GET("/:id", controllers.GetPostpartumCheckupByID)
	routes.PUT("/:id", controllers.UpdatePostpartumCheckup)
	routes.DELETE("/:id", controllers.DeletePostpartumCheckup)
}
