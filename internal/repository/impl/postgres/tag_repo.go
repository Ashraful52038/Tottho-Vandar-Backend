package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *tagRepository {
	return &tagRepository{
		db: db,
	}
}

// Create - নতুন ট্যাগ তৈরি (Like repository pattern)
func (r *tagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(tag).Error
}

// FindByID - আইডি দিয়ে ট্যাগ খুঁজে বের করা (Post repository pattern - nil return)
func (r *tagRepository) FindByID(ctx context.Context, id uint) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).First(&tag, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Post/Like repository মত nil return
	}
	return &tag, err
}

// FindBySlug - স্লাগ দিয়ে ট্যাগ খুঁজে বের করা
func (r *tagRepository) FindBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// FindByName - নাম দিয়ে ট্যাগ খুঁজে বের করা
func (r *tagRepository) FindByName(ctx context.Context, name string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tag, err
}

// FindAll - সব ট্যাগ লিস্ট করা (পেজিনেটেড) (Post repository pattern)
func (r *tagRepository) FindAll(ctx context.Context, page, limit int) ([]domain.Tag, int64, error) {
	var tags []domain.Tag
	var total int64

	offset := (page - 1) * limit

	// Total count
	if err := r.db.WithContext(ctx).Model(&domain.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated tags
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("name ASC").
		Find(&tags).Error

	return tags, total, err
}

// FindPopular - জনপ্রিয় ট্যাগ (সবচেয়ে বেশি ব্যবহৃত)
func (r *tagRepository) FindPopular(ctx context.Context, limit int) ([]domain.Tag, error) {
	var tags []domain.Tag

	err := r.db.WithContext(ctx).
		Model(&domain.Tag{}).
		Select("tags.*, COUNT(post_tags.post_id) as posts_count").
		Joins("LEFT JOIN post_tags ON post_tags.tag_id = tags.id").
		Group("tags.id").
		Order("posts_count DESC").
		Limit(limit).
		Find(&tags).Error

	return tags, err
}

// FindByPostID - নির্দিষ্ট পোস্টের ট্যাগ লিস্ট
func (r *tagRepository) FindByPostID(ctx context.Context, postID uint) ([]domain.Tag, error) {
	var tags []domain.Tag

	err := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID).
		Order("tags.name ASC").
		Find(&tags).Error

	return tags, err
}

// Update - ট্যাগ আপডেট
func (r *tagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	tag.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete - ট্যাগ ডিলিট (soft delete) - Post repository pattern
func (r *tagRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Tag{}, id).Error
}

// SyncPostTags - পোস্টের ট্যাগ সিঙ্ক্রোনাইজ করা
func (r *tagRepository) SyncPostTags(ctx context.Context, postID uint, tagIDs []uint) error {
	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Remove existing tags
		if err := tx.Where("post_id = ?", postID).Delete(&domain.PostTag{}).Error; err != nil {
			return err
		}

		// Add new tags if any
		if len(tagIDs) > 0 {
			postTags := make([]domain.PostTag, len(tagIDs))
			for i, tagID := range tagIDs {
				postTags[i] = domain.PostTag{
					PostID:    postID,
					TagID:     tagID,
					CreatedAt: time.Now(),
				}
			}
			if err := tx.Create(&postTags).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetTagsByPostIDs - একাধিক পোস্টের ট্যাগ একসাথে পাওয়ার জন্য (optional)
func (r *tagRepository) GetTagsByPostIDs(ctx context.Context, postIDs []uint) (map[uint][]domain.Tag, error) {
	var postTags []struct {
		PostID uint
		domain.Tag
	}

	err := r.db.WithContext(ctx).
		Table("post_tags").
		Select("post_tags.post_id, tags.*").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id IN ?", postIDs).
		Order("tags.name ASC").
		Scan(&postTags).Error

	if err != nil {
		return nil, err
	}

	// Group tags by post ID
	result := make(map[uint][]domain.Tag)
	for _, pt := range postTags {
		result[pt.PostID] = append(result[pt.PostID], pt.Tag)
	}

	return result, nil
}
