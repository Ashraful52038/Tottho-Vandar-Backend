package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/middleware"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func setupMiddleware(e *echo.Echo, jwtService *jwt.JWTService) *echo.Group {
	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderAccessControlAllowOrigin,
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Static("/uploads", "./uploads")

	// Protected group with auth middleware
	protected := e.Group("/api")
	protected.Use(middleware.AuthMiddlewareFunc(jwtService))

	return protected
}
