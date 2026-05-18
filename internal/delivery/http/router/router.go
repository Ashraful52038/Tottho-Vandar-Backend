package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/middleware"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct {
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	commentHandler *handler.CommentHandler
	likeHandler    *handler.LikeHandler
	feedHandler    *handler.FeedHandler
	tagHandler     *handler.TagHandler
	jwtService     *jwt.JWTService
	uploadHandler  *handler.UploadHandler
	allowedOrigins []string
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	likeHandler *handler.LikeHandler,
	tagHandler *handler.TagHandler,
	feedHandler *handler.FeedHandler,
	jwtService *jwt.JWTService,
	uploadHandler *handler.UploadHandler,
	allowedOrigins []string,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		postHandler:    postHandler,
		commentHandler: commentHandler,
		likeHandler:    likeHandler,
		tagHandler:     tagHandler,
		feedHandler:    feedHandler,
		jwtService:     jwtService,
		uploadHandler:  uploadHandler,
		allowedOrigins: allowedOrigins,
	}
}

func (r *Router) SetupRoutes(e *echo.Echo) {

	e.Static("/uploads", "./uploads")

	// ✅ CORS middleware - প্রথমেই যোগ করো
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: r.allowedOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Setup middleware and get groups
	api := e.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.GET("/verify-email", r.authHandler.VerifyEmail)
			auth.POST("/forget-password", r.authHandler.ForgotPassword)
			auth.POST("/reset-password", r.authHandler.ResetPassword)
			auth.POST("/resend-verification", r.authHandler.ResendVerification)
		}

		// Public Feed Route
		api.GET("/feed/public", r.feedHandler.GetPublicFeed)

		// Public post routes
		api.GET("/posts", r.postHandler.GetAllPosts)
		api.GET("/posts/search", r.postHandler.SearchPosts)
		api.GET("/posts/:id", r.postHandler.GetPostByID)

		// Public tag routes
		api.GET("/tags", r.tagHandler.GetAllTags)
		api.GET("/tags/popular", r.tagHandler.GetPopularTags)
		api.GET("/tags/:id", r.tagHandler.GetTagByID)

		// Public comment routes
		api.GET("/comments/posts/:postId", r.commentHandler.GetByPost)

		// Public like routes
		api.GET("/likes/posts/:postId", r.likeHandler.GetPostLikes)
		api.GET("/likes/comments/:commentId", r.likeHandler.GetCommentLikes)

		// Public user routes
		api.GET("/users/most-followed", r.userHandler.GetMostFollowedUsers)
		api.GET("/users/:userId", r.userHandler.GetUserProfile)

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

			// Post routes (protected)
			post := protected.Group("/posts")
			{
				post.POST("", r.postHandler.CreatePost)
				post.PUT("/:id", r.postHandler.UpdatePost)
				post.DELETE("/:id", r.postHandler.DeletePost)
				post.GET("/my-posts", r.postHandler.GetMyPosts)
				post.POST("/:postId/comments", r.commentHandler.Create)
				post.GET("/:postId/comments", r.commentHandler.GetByPost)
			}

			// Protected tag routes
			tag := protected.Group("/tags")
			{
				tag.POST("", r.tagHandler.CreateTag)
				tag.PUT("/:id", r.tagHandler.UpdateTag)
				tag.DELETE("/:id", r.tagHandler.DeleteTag)
			}

			// Comment routes
			comment := protected.Group("/comments")
			{
				comment.PUT("/:id", r.commentHandler.Update)
				comment.DELETE("/:id", r.commentHandler.Delete)
			}

			// Like routes
			like := protected.Group("/likes")
			{
				like.POST("/posts/:postId", r.likeHandler.TogglePost)
				like.POST("/comments/:commentId", r.likeHandler.ToggleComment)
			}

			// Feed Routes (Protected)
			feed := protected.Group("/feed")
			{
				feed.GET("", r.feedHandler.GetPersonalizedFeed)
				feed.GET("/tags/followed", r.feedHandler.GetFollowedTags)
				feed.POST("/tags/:id/follow", r.feedHandler.FollowTag)
				feed.POST("/tags/:id/unfollow", r.feedHandler.UnfollowTag)
			}

			// Upload routes
			upload := protected.Group("/upload")
			{
				upload.POST("/image", r.uploadHandler.UploadImage)
				upload.POST("/avatar", r.userHandler.UploadAvatar)
			}

			// User profile routes
			userProfile := protected.Group("/users")
			{
				userProfile.GET("/:userId/profile", r.userHandler.GetUserProfile)
				userProfile.GET("/:userId/posts", r.userHandler.GetUserPosts)
				userProfile.GET("/:userId/comments", r.userHandler.GetUserComments)
				userProfile.GET("/:userId/likes", r.userHandler.GetUserLikes)
				userProfile.GET("/:userId/followers", r.userHandler.GetUserFollowers)
				userProfile.GET("/:userId/following", r.userHandler.GetUserFollowing)
				userProfile.GET("/:userId/follow/status", r.userHandler.GetFollowStatus)
				userProfile.POST("/:userId/follow", r.userHandler.FollowUser)
				userProfile.DELETE("/:userId/follow", r.userHandler.UnfollowUser)
			}

			// Change password
			protected.POST("/auth/change-password", r.authHandler.ChangePassword)
		}
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":    "ok",
			"websocket": "enabled",
		})
	})
}
