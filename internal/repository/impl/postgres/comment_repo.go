package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

// NewCommentRepository
func NewCommentRepository(db *gorm.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	log.Printf("Creating comment - AuthorID: %d, PostID: %d, Content: %s",
		comment.AuthorID, comment.PostID, comment.Content)

	db := r.db.WithContext(ctx).Debug()

	result := db.Select("Content", "PostID", "AuthorID", "CreatedAt", "UpdatedAt").Create(comment)

	if result.Error != nil {
		log.Printf("DB Error: %v", result.Error)
		return result.Error
	}

	log.Printf("Comment created with ID: %d", comment.ID)
	return nil
}

func (r *commentRepository) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.WithContext(ctx).Preload("User").First(&comment, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &comment, err
}

func (r *commentRepository) FindByPostID(ctx context.Context, postID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.WithContext(ctx).
		Where("post_id = ?", postID).
		Preload("User").
		Order("created_at asc").
		Find(&comments).Error
	return comments, err
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}

func (r *commentRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]domain.Comment, int64, error) {
	var comments []domain.Comment
	var total int64

	db := r.db.WithContext(ctx).Model(&domain.Comment{}).Where("author_id = ?", userID)

	// Total count
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Comments with pagination
	err = db.Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	// Post titles আলাদা করে fetch করতে হবে
	for i := range comments {
		var post domain.Post
		r.db.WithContext(ctx).Select("title").First(&post, comments[i].PostID)
		comments[i].PostTitle = post.Title
	}

	return comments, total, nil
}

func (r *commentRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("author_id = ?", userID).
		Count(&count).Error
	return count, err
}
