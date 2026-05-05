package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupPublicFeedRoutes(api *echo.Group, feedHandler *handler.FeedHandler) {
	api.GET("/feed/public", feedHandler.GetPublicFeed)
}

func setupProtectedFeedRoutes(protected *echo.Group, feedHandler *handler.FeedHandler) {
	feed := protected.Group("/feed")
	{
		feed.GET("", feedHandler.GetPersonalizedFeed)
		feed.GET("/tags/followed", feedHandler.GetFollowedTags)
		feed.POST("/tags/:id/follow", feedHandler.FollowTag)
		feed.POST("/tags/:id/unfollow", feedHandler.UnfollowTag)
	}
}
