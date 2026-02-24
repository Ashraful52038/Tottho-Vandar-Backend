package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type PostUsecase interface {
	Create(ctx context.Context, userID uint, req *domain.CreatePostRequest) (*domain.Post, error)
	Update(ctx context.Context, postID uint, userID uint, req *domain.UpdatePostRequest) (*domain.Post, error)
	Delete(ctx context.Context, postID uint, userID uint) error
	GetByID(ctx context.Context, id uint) (*domain.Post, error)
	GetAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error)
	GetByUserID(ctx context.Context, userID uint) ([]domain.Post, error)
	GetByTagID(ctx context.Context, tagID uint, page, limit int) ([]domain.Post, int64, error)
	GetByTags(ctx context.Context, tagIDs []uint, page, limit int) ([]domain.Post, int64, error)
}

type postUsecase struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
	tagRepo  repository.TagRepository
}

func NewPostUsecase(
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	tagRepo repository.TagRepository,
) PostUsecase {
	return &postUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
		tagRepo:  tagRepo,
	}
}

// ✅ পরিবর্তন ১: Create মেথড - TagIDs সরাসরি ব্যবহার
func (u *postUsecase) Create(ctx context.Context, userID uint, req *domain.CreatePostRequest) (*domain.Post, error) {
	// Check if user exists
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	// Check if user is verified
	if !user.Verified {
		return nil, errors.New("email not verified")
	}

	// ✅ TagIDs সরাসরি ব্যবহার করুন (কনভার্সনের দরকার নেই)
	tagIDs := req.TagIDs

	// Validate at least one tag
	if len(tagIDs) == 0 {
		return nil, errors.New("post must have at least one tag")
	}

	// Validate tags exist
	for _, tagID := range tagIDs {
		tag, err := u.tagRepo.FindByID(ctx, tagID)
		if err != nil || tag == nil {
			return nil, fmt.Errorf("invalid tag id: %d", tagID)
		}
	}

	post := &domain.Post{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID,
		Published: req.Published,
	}

	// Create post
	err = u.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	// Sync tags
	if len(tagIDs) > 0 {
		err = u.tagRepo.SyncPostTags(ctx, post.ID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

// ✅ পরিবর্তন ২: Update মেথড - TagIDs সরাসরি ব্যবহার
func (u *postUsecase) Update(ctx context.Context, postID uint, userID uint, req *domain.UpdatePostRequest) (*domain.Post, error) {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}

	// Check if user is the author
	if post.AuthorID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Published != nil {
		post.Published = *req.Published
	}

	// Update post
	err = u.postRepo.Update(ctx, post)
	if err != nil {
		return nil, err
	}

	// Update tags if provided
	if req.TagIDs != nil { // ✅ TagIDs ব্যবহার
		tagIDs := req.TagIDs

		// Validate tags exist
		for _, tagID := range tagIDs {
			tag, err := u.tagRepo.FindByID(ctx, tagID)
			if err != nil || tag == nil {
				return nil, fmt.Errorf("invalid tag id: %d", tagID)
			}
		}

		// Sync tags
		err = u.tagRepo.SyncPostTags(ctx, post.ID, tagIDs)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

// ✅ পরিবর্তন ৩: Delete মেথড - nil check যোগ
func (u *postUsecase) Delete(ctx context.Context, postID uint, userID uint) error {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return errors.New("post not found")
	}

	if post.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return u.postRepo.Delete(ctx, postID)
}

// ✅ পরিবর্তন ৪: GetByID - nil check যোগ
func (u *postUsecase) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	post, err := u.postRepo.FindByID(ctx, id)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}
	return post, nil
}

// ✅ পরিবর্তন ৫: GetAll - pagination validation
func (u *postUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return u.postRepo.FindAll(ctx, page, limit)
}

// ✅ পরিবর্তন ৬: GetByUserID - user validation
func (u *postUsecase) GetByUserID(ctx context.Context, userID uint) ([]domain.Post, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return u.postRepo.FindByUserID(ctx, userID)
}

// GetByTagID - নির্দিষ্ট ট্যাগের পোস্ট লিস্ট
func (u *postUsecase) GetByTagID(ctx context.Context, tagID uint, page, limit int) ([]domain.Post, int64, error) {
	// Check if tag exists
	tag, err := u.tagRepo.FindByID(ctx, tagID)
	if err != nil || tag == nil {
		return nil, 0, errors.New("tag not found")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return u.postRepo.FindByTagID(ctx, tagID, page, limit)
}

// GetByTags - একাধিক ট্যাগের পোস্ট লিস্ট
func (u *postUsecase) GetByTags(ctx context.Context, tagIDs []uint, page, limit int) ([]domain.Post, int64, error) {
	// Validate tags
	for _, tagID := range tagIDs {
		tag, err := u.tagRepo.FindByID(ctx, tagID)
		if err != nil || tag == nil {
			return nil, 0, fmt.Errorf("invalid tag id: %d", tagID)
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	return u.postRepo.SearchByTags(ctx, tagIDs, page, limit)
}
