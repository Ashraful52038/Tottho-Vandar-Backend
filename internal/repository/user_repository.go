package repository

import (
	"context"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*domain.User, error)
	FindByResetToken(ctx context.Context, token string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uint) error
	VerifyEmail(ctx context.Context, userID uint) error
	SetResetToken(ctx context.Context, userID uint, token string, expiry time.Time) error
	FindMostFollowed(ctx context.Context, limit int) ([]domain.UserWithFollowCount, error)
}
