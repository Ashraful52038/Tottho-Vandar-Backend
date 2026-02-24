package postgres

import (
	"context"
	"errors"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *postRepository) FindByID(ctx context.Context, id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.WithContext(ctx).Preload("Author").First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &post, err
}

func (r *postRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	offset := (page - 1) * limit

	err := r.db.WithContext(ctx).Model(&domain.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Preload("Author").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) FindByUserID(ctx context.Context, userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.WithContext(ctx).
		Where("author_id = ?", userID).
		Order("created_at desc").
		Find(&posts).Error
	return posts, err
}

func (r *postRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Post{}, id).Error
}

func (r *postRepository) FindByTagID(ctx context.Context, tagID uint, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	offset := (page - 1) * limit

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated posts
	err = r.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id = ?", tagID).
		Offset(offset).
		Limit(limit).
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) SearchByTags(ctx context.Context, tagIDs []uint, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	offset := (page - 1) * limit

	// Count total
	err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id IN ?", tagIDs).
		Group("posts.id").
		Having("COUNT(DISTINCT post_tags.tag_id) = ?", len(tagIDs)).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated posts
	err = r.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id IN ?", tagIDs).
		Group("posts.id").
		Having("COUNT(DISTINCT post_tags.tag_id) = ?", len(tagIDs)).
		Offset(offset).
		Limit(limit).
		Order("posts.created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
