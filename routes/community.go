package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shem958/cycle-backend/controllers"
	"github.com/shem958/cycle-backend/middleware"
)

func RegisterCommunityRoutes(r *gin.RouterGroup) {
	community := r.Group("/community")
	community.Use(middleware.AuthMiddleware()) // protect all community endpoints

	{
		// Posts
		community.GET("/posts", controllers.GetAllPosts)
		community.GET("/posts/:id", controllers.GetPostByID)

		// Tags
		community.GET("/tags", controllers.GetAllTags)

		// Create/update/delete routes — include BlockSuspendedMiddleware
		protected := community.Use(middleware.BlockSuspendedMiddleware())

		protected.POST("/posts", controllers.CreatePost)
		protected.POST("/comments", controllers.CreateComment)
		protected.POST("/replies", controllers.ReplyToComment)
		protected.POST("/reports", controllers.ReportContent)
		protected.POST("/reactions", controllers.ReactToContent)
		protected.DELETE("/reactions", controllers.RemoveReaction)
	}
}
