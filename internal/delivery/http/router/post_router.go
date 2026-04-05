package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupPublicPostRoutes(api *echo.Group, postHandler *handler.PostHandler) {
	api.GET("/posts", postHandler.GetAllPosts)
	api.GET("/posts/search", postHandler.SearchPosts)
	api.GET("/posts/:id", postHandler.GetPostByID)
}

func setupProtectedPostRoutes(protected *echo.Group, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler) {
	post := protected.Group("/posts")
	{
		post.POST("", postHandler.CreatePost)
		post.PUT("/:id", postHandler.UpdatePost)
		post.DELETE("/:id", postHandler.DeletePost)
		post.GET("/my-posts", postHandler.GetMyPosts)
		post.POST("/:postId/comments", commentHandler.Create)
		post.GET("/:postId/comments", commentHandler.GetByPost)
	}
}
