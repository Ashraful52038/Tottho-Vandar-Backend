package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupUploadRoutes(protected *echo.Group, uploadHandler *handler.UploadHandler) {
	upload := protected.Group("/upload")
	{
		upload.POST("/image", uploadHandler.UploadImage)
	}
}
