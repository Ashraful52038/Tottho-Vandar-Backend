package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupPublicCommentRoutes(api *echo.Group, commentHandler *handler.CommentHandler) {
	api.GET("/comments/posts/:postId", commentHandler.GetByPost)
}

func setupProtectedCommentRoutes(protected *echo.Group, commentHandler *handler.CommentHandler) {
	comment := protected.Group("/comments")
	{
		comment.PUT("/:id", commentHandler.Update)
		comment.DELETE("/:id", commentHandler.Delete)
	}
}
