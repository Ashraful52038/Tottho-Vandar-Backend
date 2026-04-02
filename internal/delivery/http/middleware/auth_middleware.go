package middleware

import (
	"net/http"
	"strings"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func NewAuthMiddleware(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "authorization header required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid authorization header format",
			})
		}

		claims, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid or expired token",
			})
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		return next(c)
	}
}

// Helper function to use in router
func AuthMiddlewareFunc(jwtService *jwt.JWTService) echo.MiddlewareFunc {
	middleware := NewAuthMiddleware(jwtService)
	return middleware.Handle
}
