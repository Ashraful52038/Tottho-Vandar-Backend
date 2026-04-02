package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	err := r.db.WithContext(ctx).
		Preload("Author").
		Preload("Tags").
		First(&post, id).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	log.Printf("Repository FindByID: Post ID=%d, FeaturedImage=%q", post.ID, post.FeaturedImage)

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

func (r *postRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Where("author_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Paginated posts
	err = r.db.WithContext(ctx).
		Where("author_id = ?", userID).
		Preload("Author").
		Preload("Tags").
		Offset(offset).
		Limit(limit).
		Order("created_at desc").
		Find(&posts).Error

	return posts, total, err
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

// SearchPosts
func (r *postRepository) SearchPosts(ctx context.Context, params *repository.SearchParams) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	// Base query with preloads
	query := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Preload("Author").
		Preload("Tags").
		Where("published = ?", true)

	// Keyword search (title, content,author name, tag name)
	if params.Query != "" {
		searchTerm := "%" + strings.ToLower(params.Query) + "%"
		query = query.
			Joins("LEFT JOIN users ON users.id = posts.author_id").
			Joins("LEFT JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("LEFT JOIN tags ON tags.id = post_tags.tag_id").
			Where(
				"LOWER(posts.title) LIKE ? OR LOWER(posts.content) LIKE ? OR LOWER(users.name) LIKE ? OR LOWER(tags.name) LIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm,
			).
			Group("posts.id")
	}

	// Filter by tags
	if len(params.TagIDs) > 0 {
		subQuery := r.db.Table("post_tags").
			Select("post_id").
			Where("tag_id IN ?", params.TagIDs).
			Group("post_id").
			Having("COUNT(DISTINCT tag_id) = ?", len(params.TagIDs))

		query = query.Where("id IN (?)", subQuery)
	}

	// Filter by author
	if params.AuthorID != nil {
		query = query.Where("author_id = ?", *params.AuthorID)
	}

	// Pagination
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

// GetPersonalizedFeed
func (r *postRepository) GetPersonalizedFeed(ctx context.Context, userID uint, params *domain.FeedQueryParams) ([]domain.FeedPost, int64, error) {
	var posts []domain.FeedPost
	var total int64

	offset := (params.Page - 1) * params.Limit

	followedTags := r.db.Table("user_followed_tags").
		Select("tag_id").
		Where("user_id = ?", userID)

	followedUsers := r.db.Table("user_followers").
		Select("following_id").
		Where("follower_id = ? AND status = 'accepted'", userID)

	countQuery := r.db.WithContext(ctx).
		Table("posts p").
		Joins("LEFT JOIN post_tags pt ON pt.post_id = p.id").
		Where("p.deleted_at IS NULL").
		Where("p.published = ?", true).
		Where("(p.author_id IN (?) OR pt.tag_id IN (?) OR p.author_id = ?)",
			followedUsers, followedTags, userID).
		Distinct("p.id")

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sortField string
	switch params.SortBy {
	case "likes_count":
		sortField = "likes_count"
	case "comments_count":
		sortField = "comments_count"
	default:
		sortField = "p.created_at"
	}

	query := r.db.WithContext(ctx).
		Table("posts p").
		Select(`
			p.id, p.title, p.content, p.author_id, p.created_at, p.published,
			u.name as author_name, u.email as author_email,
			COUNT(DISTINCT l.id) as likes_count,
			COUNT(DISTINCT c.id) as comments_count,
			EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = p.id) as is_liked,
			EXISTS(SELECT 1 FROM bookmarks WHERE user_id = ? AND post_id = p.id) as is_bookmarked,
			(
				SELECT string_agg(t.name, ',') 
				FROM tags t 
				JOIN post_tags pt ON pt.tag_id = t.id 
				WHERE pt.post_id = p.id
			) as tags
		`, userID, userID).
		Joins("LEFT JOIN users u ON u.id = p.author_id").
		Joins("LEFT JOIN likes l ON l.post_id = p.id").
		Joins("LEFT JOIN comments c ON c.post_id = p.id").
		Joins("LEFT JOIN post_tags pt ON pt.post_id = p.id").
		Where("p.deleted_at IS NULL").
		Where("p.published = ?", true).
		Where("(p.author_id IN (?) OR pt.tag_id IN (?) OR p.author_id = ?)",
			followedUsers, followedTags, userID).
		Group("p.id, u.name, u.email, p.title, p.content, p.author_id, p.created_at, p.published").
		Order(fmt.Sprintf("%s %s", sortField, params.SortOrder)).
		Offset(offset).
		Limit(params.Limit)

	if err := query.Scan(&posts).Error; err != nil {
		return nil, 0, err
	}

	for i := range posts {
		if len(posts[i].Tags) > 0 {
			posts[i].Tags = strings.Split(posts[i].Tags[0], ",")
		} else {
			posts[i].Tags = []string{}
		}
	}

	return posts, total, nil
}

// GetPublicFeed
func (r *postRepository) GetPublicFeed(ctx context.Context, params *domain.FeedQueryParams) ([]domain.FeedPost, int64, error) {
	var posts []domain.FeedPost
	var total int64

	offset := (params.Page - 1) * params.Limit

	if err := r.db.WithContext(ctx).Model(&domain.Post{}).Where("deleted_at IS NULL AND published = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var sortField string
	switch params.SortBy {
	case "likes_count":
		sortField = "likes_count"
	case "comments_count":
		sortField = "comments_count"
	default:
		sortField = "p.created_at"
	}

	query := r.db.WithContext(ctx).
		Table("posts p").
		Select(`
			p.id, p.title, p.content, p.author_id, p.created_at, p.published,
			u.name as author_name, u.email as author_email,
			COUNT(DISTINCT l.id) as likes_count,
			COUNT(DISTINCT c.id) as comments_count,
			false as is_liked,
			false as is_bookmarked,
			(
				SELECT string_agg(t.name, ',') 
				FROM tags t 
				JOIN post_tags pt ON pt.tag_id = t.id 
				WHERE pt.post_id = p.id
			) as tags
		`).
		Joins("LEFT JOIN users u ON u.id = p.author_id").
		Joins("LEFT JOIN likes l ON l.post_id = p.id").
		Joins("LEFT JOIN comments c ON c.post_id = p.id").
		Where("p.deleted_at IS NULL").
		Where("p.published = ?", true).
		Group("p.id, u.name, u.email, p.title, p.content, p.author_id, p.created_at, p.published").
		Order(fmt.Sprintf("%s %s", sortField, params.SortOrder)).
		Offset(offset).
		Limit(params.Limit)

	if err := query.Scan(&posts).Error; err != nil {
		return nil, 0, err
	}

	for i := range posts {
		if len(posts[i].Tags) > 0 {
			posts[i].Tags = strings.Split(posts[i].Tags[0], ",")
		} else {
			posts[i].Tags = []string{}
		}
	}

	return posts, total, nil
}

// ToggleFollowTag
func (r *postRepository) ToggleFollowTag(ctx context.Context, userID, tagID uint) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("user_followed_tags").
		Where("user_id = ? AND tag_id = ?", userID, tagID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return r.db.WithContext(ctx).
			Table("user_followed_tags").
			Where("user_id = ? AND tag_id = ?", userID, tagID).
			Delete(nil).Error
	}

	return r.db.WithContext(ctx).
		Table("user_followed_tags").
		Create(map[string]interface{}{
			"user_id": userID,
			"tag_id":  tagID,
		}).Error
}

// GetFollowedTags
func (r *postRepository) GetFollowedTags(ctx context.Context, userID uint) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.WithContext(ctx).
		Table("tags t").
		Joins("JOIN user_followed_tags uft ON uft.tag_id = t.id").
		Where("uft.user_id = ?", userID).
		Find(&tags).Error
	return tags, err
}

// GetFollowedTagIDs
func (r *postRepository) GetFollowedTagIDs(ctx context.Context, userID uint) ([]uint, error) {
	var tagIDs []uint
	err := r.db.WithContext(ctx).
		Table("user_followed_tags").
		Where("user_id = ?", userID).
		Pluck("tag_id", &tagIDs).Error
	return tagIDs, err
}

func (r *postRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Where("author_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

func (r *postRepository) IncrementCommentCount(ctx context.Context, postID uint, delta int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Post{}).
		Where("id = ?", postID).
		Update("comments", gorm.Expr("comments + ?", delta)).
		Error
}
