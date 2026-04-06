package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupPublicLikeRoutes(api *echo.Group, likeHandler *handler.LikeHandler) {
	api.GET("/likes/posts/:postId", likeHandler.GetPostLikes)
	api.GET("/likes/comments/:commentId", likeHandler.GetCommentLikes)
}

func setupProtectedLikeRoutes(protected *echo.Group, likeHandler *handler.LikeHandler) {
	like := protected.Group("/likes")
	{
		like.POST("/posts/:postId", likeHandler.TogglePost)
		like.POST("/comments/:commentId", likeHandler.ToggleComment)
	}
}
