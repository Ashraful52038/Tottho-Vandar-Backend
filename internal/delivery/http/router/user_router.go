package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupUserRoutes(protected *echo.Group, userHandler *handler.UserHandler) {
	user := protected.Group("/user")
	{
		user.GET("/me", userHandler.GetCurrentUser)
		user.PUT("/me", userHandler.UpdateUser)
		user.POST("/logout", userHandler.Logout)
	}
}

func setupUserProfileRoutes(protected *echo.Group, userHandler *handler.UserHandler) {
	userProfile := protected.Group("/users")
	{
		userProfile.GET("/:userId", userHandler.GetUserProfile)
		userProfile.GET("/:userId/profile", userHandler.GetUserProfile)
		userProfile.GET("/:userId/posts", userHandler.GetUserPosts)
		userProfile.GET("/:userId/comments", userHandler.GetUserComments)
		userProfile.GET("/:userId/likes", userHandler.GetUserLikes)
		userProfile.GET("/:userId/followers", userHandler.GetUserFollowers)
		userProfile.GET("/:userId/following", userHandler.GetUserFollowing)
		userProfile.GET("/:userId/follow/status", userHandler.GetFollowStatus)
		userProfile.POST("/:userId/follow", userHandler.FollowUser)
		userProfile.DELETE("/:userId/follow", userHandler.UnfollowUser)
	}
}
