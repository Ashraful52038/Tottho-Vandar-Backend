package usecase

import (
	"context"
	"errors"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type TagUsecase interface {
	Create(ctx context.Context, name string) (*domain.Tag, error)
	GetByID(ctx context.Context, id uint) (*domain.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Tag, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Tag, int64, error)
	GetPopular(ctx context.Context, limit int) ([]domain.Tag, error)
	GetByPostID(ctx context.Context, postID uint) ([]domain.Tag, error)
	Update(ctx context.Context, id uint, name string) (*domain.Tag, error)
	Delete(ctx context.Context, id uint) error
}

type tagUsecase struct {
	tagRepo  repository.TagRepository
	postRepo repository.PostRepository
}

func NewTagUsecase(
	tagRepo repository.TagRepository,
	postRepo repository.PostRepository,
) TagUsecase {
	return &tagUsecase{
		tagRepo:  tagRepo,
		postRepo: postRepo,
	}
}

func (u *tagUsecase) validateTagName(name string) error {
	if name == "" {
		return errors.New("tag name cannot be empty")
	}
	if len(name) < 2 || len(name) > 50 {
		return errors.New("tag name must be between 2 and 50 characters")
	}
	return nil
}

func (u *tagUsecase) Create(ctx context.Context, name string) (*domain.Tag, error) {
	// Validate name
	if err := u.validateTagName(name); err != nil {
		return nil, err
	}

	// Check if tag already exists
	existingTag, _ := u.tagRepo.FindByName(ctx, name)
	if existingTag != nil {
		return nil, errors.New("tag with this name already exists")
	}

	// Create new tag
	tag := &domain.Tag{
		Name: name,
	}

	err := u.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// GetByID - আইডি দিয়ে ট্যাগ খুঁজে বের করা (Comment/Like pattern অনুসারে)
func (u *tagUsecase) GetByID(ctx context.Context, id uint) (*domain.Tag, error) {
	tag, err := u.tagRepo.FindByID(ctx, id)
	if err != nil || tag == nil {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

// GetBySlug - স্লাগ দিয়ে ট্যাগ খুঁজে বের করা
func (u *tagUsecase) GetBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	tag, err := u.tagRepo.FindBySlug(ctx, slug)
	if err != nil || tag == nil {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

// GetAll - সব ট্যাগ লিস্ট করা (পেজিনেটেড)
func (u *tagUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Tag, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	tags, total, err := u.tagRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return tags, total, nil
}

// GetPopular - জনপ্রিয় ট্যাগ (সবচেয়ে বেশি ব্যবহৃত)
func (u *tagUsecase) GetPopular(ctx context.Context, limit int) ([]domain.Tag, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	tags, err := u.tagRepo.FindPopular(ctx, limit)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetByPostID - নির্দিষ্ট পোস্টের ট্যাগ লিস্ট (Comment/Like pattern অনুসারে)
func (u *tagUsecase) GetByPostID(ctx context.Context, postID uint) ([]domain.Tag, error) {
	// Check if post exists
	_, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	tags, err := u.tagRepo.FindByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Update - ট্যাগ আপডেট (Comment/Like pattern অনুসারে)
func (u *tagUsecase) Update(ctx context.Context, id uint, name string) (*domain.Tag, error) {
	// Validate name
	if err := u.validateTagName(name); err != nil {
		return nil, err
	}

	// Find the tag
	tag, err := u.tagRepo.FindByID(ctx, id)
	if err != nil || tag == nil {
		return nil, errors.New("tag not found")
	}

	// Check if new name already exists (if name changed)
	if tag.Name != name {
		existingTag, _ := u.tagRepo.FindByName(ctx, name)
		if existingTag != nil {
			return nil, errors.New("tag with this name already exists")
		}
	}

	// Update tag
	tag.Name = name
	err = u.tagRepo.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete
func (u *tagUsecase) Delete(ctx context.Context, id uint) error {
	// Check if tag exists
	tag, err := u.tagRepo.FindByID(ctx, id)
	if err != nil || tag == nil {
		return errors.New("tag not found")
	}

	// Check if tag is used by any posts
	posts, _, err := u.postRepo.FindByTagID(ctx, id, 1, 1)
	if err != nil {
		return err
	}
	if len(posts) > 0 {
		return errors.New("cannot delete tag that is being used by posts")
	}

	return u.tagRepo.Delete(ctx, id)
}
