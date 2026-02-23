package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *likeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, like *domain.Like) error {
	like.CreatedAt = time.Now()
	like.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *likeRepository) Delete(ctx context.Context, userID uint, postID, commentID *uint) error {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if postID != nil {
		query = query.Where("post_id = ?", *postID)
	}
	if commentID != nil {
		query = query.Where("comment_id = ?", *commentID)
	}

	return query.Delete(&domain.Like{}).Error
}

func (r *likeRepository) FindByUserAndPost(ctx context.Context, userID uint, postID uint) (*domain.Like, error) {
	var like domain.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ? AND comment_id IS NULL", userID, postID).
		First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &like, err
}

func (r *likeRepository) FindByUserAndComment(ctx context.Context, userID uint, commentID uint) (*domain.Like, error) {
	var like domain.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &like, err
}

func (r *likeRepository) FindByPostID(ctx context.Context, postID uint) ([]domain.Like, error) {
	var likes []domain.Like
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Preload("User").
		Order("created_at desc").
		Find(&likes).Error
	return likes, err
}

func (r *likeRepository) CountByPostID(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	return count, err
}

func (r *likeRepository) FindByCommentID(ctx context.Context, commentID uint) ([]domain.Like, error) {
	var likes []domain.Like
	err := r.db.WithContext(ctx).
		Where("comment_id = ?", commentID).
		Preload("User").
		Order("created_at desc").
		Find(&likes).Error
	return likes, err
}

func (r *likeRepository) CountByCommentID(ctx context.Context, commentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Like{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error
	return count, err
}

func (r *likeRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.Like, error) {
	var likes []domain.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Post").
		Preload("Comment").
		Order("created_at desc").
		Find(&likes).Error
	return likes, err
}
