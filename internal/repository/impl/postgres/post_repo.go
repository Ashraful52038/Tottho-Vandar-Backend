package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
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

// ✅ SearchPosts - সার্চ ও ফিল্টার ফাংশন
func (r *postRepository) SearchPosts(ctx context.Context, params *repository.SearchParams) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Base query with preloads
	query := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Preload("Author").
		Preload("Tags").
		Where("published = ?", true)

	// 🔍 Keyword search (title or content)
	if params.Query != "" {
		searchTerm := "%" + strings.ToLower(params.Query) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(content) LIKE ?",
			searchTerm, searchTerm,
		)
	}

	// 🏷️ Filter by tags
	if len(params.TagIDs) > 0 {
		subQuery := r.db.Table("post_tags").
			Select("post_id").
			Where("tag_id IN ?", params.TagIDs).
			Group("post_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(params.TagIDs))

		query = query.Where("id IN (?)", subQuery)
	}

	// 👤 Filter by author
	if params.AuthorID != nil {
		query = query.Where("author_id = ?", *params.AuthorID)
	}

	// 📊 Pagination
	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err = query.
		Offset(offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&posts).Error

	return posts, total, err
}
