package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type LikeRepository interface {
	Create(ctx context.Context, like *domain.Like) error
	Delete(ctx context.Context, userID uint, postID, commentID *uint) error
	DeleteByID(ctx context.Context, id uint) error
	FindByUserAndPost(ctx context.Context, userID uint, postID uint) (*domain.Like, error)
	FindByUserAndComment(ctx context.Context, userID uint, commentID uint) (*domain.Like, error)
	FindByPostID(ctx context.Context, postID uint) ([]domain.Like, error)
	FindByCommentID(ctx context.Context, commentID uint) ([]domain.Like, error)
	CountByPostID(ctx context.Context, postID uint) (int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	CountByCommentID(ctx context.Context, commentID uint) (int64, error)
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Like, int64, error)
}
