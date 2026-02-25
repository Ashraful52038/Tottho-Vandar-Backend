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

		// ✅ Public Feed Route
		api.GET("/feed/public", r.feedHandler.GetPublicFeed)

		// Public post routes
		api.GET("/posts", r.postHandler.GetAllPosts)
		api.GET("/posts/search", r.postHandler.SearchPosts)
		api.GET("/posts/:id", r.postHandler.GetPostByID)

		// ✅ Public tag routes
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
			}

			// ✅ Protected tag routes
			tag := protected.Group("/tags")
			{
				tag.POST("", r.tagHandler.CreateTag)
				tag.PUT("/:id", r.tagHandler.UpdateTag)
				tag.DELETE("/:id", r.tagHandler.DeleteTag)
			}

			// Comment routes
			comment := protected.Group("/comments")
			{
				comment.POST("/posts/:postId", r.commentHandler.Create)
				comment.PUT("/:id", r.commentHandler.Update)
				comment.DELETE("/:id", r.commentHandler.Delete)
			}

			// Like routes
			like := protected.Group("/likes")
			{
				like.POST("/posts/:postId", r.likeHandler.TogglePost)
				like.POST("/comments/:commentId", r.likeHandler.ToggleComment)
			}

			// ✅ Feed Routes (Protected)
			feed := protected.Group("/feed")
			{
				feed.GET("", r.feedHandler.GetPersonalizedFeed)            // GET /api/feed
				feed.GET("/tags/followed", r.feedHandler.GetFollowedTags)  // GET /api/feed/tags/followed
				feed.POST("/tags/:id/follow", r.feedHandler.FollowTag)     // POST /api/feed/tags/1/follow
				feed.POST("/tags/:id/unfollow", r.feedHandler.UnfollowTag) // POST /api/feed/tags/1/unfollow
			}
		}
	}
}
