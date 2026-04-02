package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type FollowRepository interface {
	Follow(ctx context.Context, followerID, followingID uint) error
	Unfollow(ctx context.Context, followerID, followingID uint) error
	CheckFollowStatus(ctx context.Context, followerID, followingID uint) (bool, error)
	GetFollowers(ctx context.Context, userID uint, offset, limit int) ([]domain.FollowUser, int64, error)
	GetFollowing(ctx context.Context, userID uint, offset, limit int) ([]domain.FollowUser, int64, error)
	GetFollowersCount(ctx context.Context, userID uint) (int64, error)
	GetFollowingCount(ctx context.Context, userID uint) (int64, error)
}
