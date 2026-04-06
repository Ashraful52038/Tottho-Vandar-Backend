package router

import (
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/labstack/echo/v4"
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
	// Setup middleware and get groups
	api := e.Group("/api")
	protected := setupMiddleware(e, r.jwtService, r.allowedOrigins)

	// Public routes
	setupAuthRoutes(api, r.authHandler)
	setupPublicFeedRoutes(api, r.feedHandler)
	setupPublicPostRoutes(api, r.postHandler)
	setupPublicTagRoutes(api, r.tagHandler)
	setupPublicCommentRoutes(api, r.commentHandler)
	setupPublicLikeRoutes(api, r.likeHandler)

	// Protected routes
	setupProtectedAuthRoutes(protected, r.authHandler)
	setupUserRoutes(protected, r.userHandler)
	setupUserProfileRoutes(protected, r.userHandler)
	setupProtectedPostRoutes(protected, r.postHandler, r.commentHandler)
	setupProtectedTagRoutes(protected, r.tagHandler)
	setupProtectedCommentRoutes(protected, r.commentHandler)
	setupProtectedLikeRoutes(protected, r.likeHandler)
	setupProtectedFeedRoutes(protected, r.feedHandler)
	setupUploadRoutes(protected, r.uploadHandler)
}
