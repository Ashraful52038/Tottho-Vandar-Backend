package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupPublicTagRoutes(api *echo.Group, tagHandler *handler.TagHandler) {
	api.GET("/tags", tagHandler.GetAllTags)
	api.GET("/tags/popular", tagHandler.GetPopularTags)
	api.GET("/tags/:id", tagHandler.GetTagByID)
}

func setupProtectedTagRoutes(protected *echo.Group, tagHandler *handler.TagHandler) {
	tag := protected.Group("/tags")
	{
		tag.POST("", tagHandler.CreateTag)
		tag.PUT("/:id", tagHandler.UpdateTag)
		tag.DELETE("/:id", tagHandler.DeleteTag)
	}
}
