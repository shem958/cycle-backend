package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
)

func RegisterCommunityRoutes(r *gin.RouterGroup) {
	community := r.Group("/api/community")
	{
		community.POST("/posts", controllers.CreatePost)
		community.GET("/posts", controllers.GetAllPosts)
		community.GET("/posts/:id", controllers.GetPostByID)

		community.POST("/comments", controllers.CreateComment)
		community.POST("/reports", controllers.ReportContent)
	}
}
