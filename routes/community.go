package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterCommunityRoutes(r *gin.RouterGroup) {
	community := r.Group("/community")         // removed duplicate /api
	community.Use(middleware.AuthMiddleware()) // protect all community endpoints

	{
		// Posts
		community.POST("/posts", controllers.CreatePost)
		community.GET("/posts", controllers.GetAllPosts)
		community.GET("/posts/:id", controllers.GetPostByID)

		// Comments
		community.POST("/comments", controllers.CreateComment)

		// Reporting
		community.POST("/reports", controllers.ReportContent)

		// Tags
		community.GET("/tags", controllers.GetAllTags)
	}

}
