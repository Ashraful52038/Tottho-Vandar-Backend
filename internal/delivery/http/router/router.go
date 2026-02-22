package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/middleware"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	authHandler *handler.AuthHandler
	userHandler *handler.UserHandler
	postHandler *handler.PostHandler
	jwtService  *jwt.JWTService
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	jwtService *jwt.JWTService,
) *Router {
	return &Router{
		authHandler: authHandler,
		userHandler: userHandler,
		postHandler: postHandler,
		jwtService:  jwtService,
	}
}

func (r *Router) SetupRoutes(e *echo.Echo) {
	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Public routes
	api := e.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.GET("/verify-email", r.authHandler.VerifyEmail)
			auth.POST("/forgot-password", r.authHandler.ForgotPassword)
			auth.POST("/reset-password", r.authHandler.ResetPassword)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddlewareFunc(r.jwtService))
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/me", r.userHandler.GetCurrentUser)
				user.PUT("/me", r.userHandler.UpdateUser)
				user.POST("/logout", r.userHandler.Logout)
			}

			// Post routes
			post := protected.Group("/posts")
			{
				post.POST("", r.postHandler.CreatePost)
				post.PUT("/:id", r.postHandler.UpdatePost)
				post.DELETE("/:id", r.postHandler.DeletePost)
				post.GET("/my-posts", r.postHandler.GetMyPosts)
			}
		}

		// Public post routes
		api.GET("/posts", r.postHandler.GetAllPosts)
		api.GET("/posts/:id", r.postHandler.GetPostByID)
	}
}
