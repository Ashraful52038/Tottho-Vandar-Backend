package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindByID(ctx context.Context, id uint) (*domain.Comment, error)
	FindByPostID(ctx context.Context, postID uint) ([]domain.Comment, error)
	Update(ctx context.Context, comment *domain.Comment) error
	Delete(ctx context.Context, id uint) error
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Comment, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}
