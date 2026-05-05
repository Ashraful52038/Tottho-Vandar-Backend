package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/labstack/echo/v4"
)

func setupAuthRoutes(api *echo.Group, authHandler *handler.AuthHandler) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/verify-email", authHandler.VerifyEmail)
		auth.POST("/forget-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
	}
}

func setupProtectedAuthRoutes(protected *echo.Group, authHandler *handler.AuthHandler) {
	protected.POST("/auth/change-password", authHandler.ChangePassword)
}
