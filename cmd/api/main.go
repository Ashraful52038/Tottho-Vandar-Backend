package main

import (
	"log"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/config"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/handler"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/http/router"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/email"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/jwt"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/pkg/validator"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository/impl/postgres"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load config
	cfg := config.New()

	// Initialize database
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	postRepo := postgres.NewPostRepository(db)
	commentRepo := postgres.NewCommentRepository(db)
	likeRepo := postgres.NewLikeRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	feedRepo := postRepo
	followRepo := postgres.NewFollowRepository(db)

	// Initialize services
	jwtService := jwt.NewJWTService(cfg.JWTSecret, time.Hour*24)

	emailService := email.NewEmailService(
		cfg.SMTPHost,     // "localhost"
		cfg.SMTPPort,     // "1025"
		cfg.SMTPUsername, // ""
		cfg.SMTPPassword, // ""
		cfg.EmailFrom,    // "noreply@totthovandar.com"
	)

	var emailQueue *email.EmailQueueService = nil
	// emailQueue, _ = email.NewEmailQueueService(emailService, "amqp://guest:guest@localhost:5672/")

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService, emailService)
	userUsecase := usecase.NewUserUsecase(userRepo, postRepo, commentRepo, likeRepo, followRepo)
	postUsecase := usecase.NewPostUsecase(postRepo, userRepo, tagRepo, feedRepo)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, postRepo, userRepo, emailQueue)
	likeUsecase := usecase.NewLikeUsecase(likeRepo, postRepo, commentRepo, userRepo)
	tagUsecase := usecase.NewTagUsecase(tagRepo, postRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	postHandler := handler.NewPostHandler(postUsecase)
	commentHandler := handler.NewCommentHandler(commentUsecase, likeUsecase)
	likeHandler := handler.NewLikeHandler(likeUsecase, postUsecase, commentUsecase)
	tagHandler := handler.NewTagHandler(tagUsecase)
	feedHandler := handler.NewFeedHandler(postUsecase)
	uploadHandler := handler.NewUploadHandler()

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = validator.NewCustomValidator()

	// Setup routes
	router := router.NewRouter(
		authHandler,
		userHandler,
		postHandler,
		commentHandler,
		likeHandler,
		tagHandler,
		feedHandler,
		jwtService,
		uploadHandler,
	)
	router.SetupRoutes(e)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
