package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type TagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) error
	FindByID(ctx context.Context, id uint) (*domain.Tag, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Tag, error)
	FindByName(ctx context.Context, name string) (*domain.Tag, error)
	FindAll(ctx context.Context, page, limit int) ([]domain.Tag, int64, error)
	FindPopular(ctx context.Context, limit int) ([]domain.Tag, error)
	FindByPostID(ctx context.Context, postID uint) ([]domain.Tag, error)
	Update(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id uint) error
	SyncPostTags(ctx context.Context, postID uint, tagIDs []uint) error
}
