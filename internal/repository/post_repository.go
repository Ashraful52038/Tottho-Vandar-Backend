package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id uint) (*domain.Post, error)
	FindAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error)
	FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Post, int64, error)
	Update(ctx context.Context, post *domain.Post) error
	Delete(ctx context.Context, id uint) error
	FindByTagID(ctx context.Context, tagID uint, page, limit int) ([]domain.Post, int64, error)
	SearchByTags(ctx context.Context, tagIDs []uint, page, limit int) ([]domain.Post, int64, error)
	SearchPosts(ctx context.Context, params *SearchParams) ([]domain.Post, int64, error)
	GetPersonalizedFeed(ctx context.Context, userID uint, params *domain.FeedQueryParams) ([]domain.FeedPost, int64, error)
	GetPublicFeed(ctx context.Context, params *domain.FeedQueryParams) ([]domain.FeedPost, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	IncrementCommentCount(ctx context.Context, postID uint, delta int) error
}

type SearchParams struct {
	Query    string `json:"q"`        // keyword search
	TagIDs   []uint `json:"tagIds"`   // filter by tags
	AuthorID *uint  `json:"authorId"` // filter by author
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type FeedRepository interface {
	ToggleFollowTag(ctx context.Context, userID, tagID uint) error
	GetFollowedTags(ctx context.Context, userID uint) ([]domain.Tag, error)
	GetFollowedTagIDs(ctx context.Context, userID uint) ([]uint, error)
}
