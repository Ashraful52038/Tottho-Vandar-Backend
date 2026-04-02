package postgres

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *followRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Follow(ctx context.Context, followerID, followingID uint) error {
	follow := &domain.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followingID uint) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&domain.Follow{}).Error
}

func (r *followRepository) CheckFollowStatus(ctx context.Context, followerID, followingID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepository) GetFollowers(ctx context.Context, userID uint, offset, limit int) ([]domain.FollowUser, int64, error) {
	var users []domain.FollowUser
	var total int64

	// Total count
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Followers list
	err = r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.avatar, users.bio").
		Joins("JOIN follows ON users.id = follows.follower_id").
		Where("follows.following_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("follows.created_at DESC").
		Scan(&users).Error

	return users, total, err
}

func (r *followRepository) GetFollowing(ctx context.Context, userID uint, offset, limit int) ([]domain.FollowUser, int64, error) {
	var users []domain.FollowUser
	var total int64

	// Total count
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Following list
	err = r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.avatar, users.bio").
		Joins("JOIN follows ON users.id = follows.following_id").
		Where("follows.follower_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("follows.created_at DESC").
		Scan(&users).Error

	return users, total, err
}

func (r *followRepository) GetFollowersCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *followRepository) GetFollowingCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}
