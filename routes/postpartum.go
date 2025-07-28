package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
)

func RegisterPostpartumRoutes(router *gin.Engine) {
	postpartum := router.Group("/postpartum")
	{
		postpartum.POST("/", controllers.CreatePostpartumLog)
		postpartum.GET("/:id", controllers.GetPostpartumLogs)
	}
}
