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
	}
}

func (r *Router) SetupRoutes(e *echo.Echo) {
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

	// Public routes
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

		// PUBLIC comment routes
		api.GET("/comments/posts/:postId", r.commentHandler.GetByPost)

		// PUBLIC like routes
		api.GET("/likes/posts/:postId", r.likeHandler.GetPostLikes)
		api.GET("/likes/comments/:commentId", r.likeHandler.GetCommentLikes)

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

			//Feed Routes (Protected)
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
			}

			userProfile := protected.Group("/users")
			{
				userProfile.GET("/:userId", r.userHandler.GetUserProfile)
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

			protected.POST("/auth/change-password", r.authHandler.ChangePassword)
		}
	}
}
