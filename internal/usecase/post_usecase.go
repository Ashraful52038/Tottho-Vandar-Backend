package usecase

import (
	"context"
	"errors"

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
}

type postUsecase struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostUsecase(postRepo repository.PostRepository, userRepo repository.UserRepository) PostUsecase {
	return &postUsecase{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

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

	post := &domain.Post{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID,
		Published: req.Published,
	}

	err = u.postRepo.Create(ctx, post)
	return post, err
}

func (u *postUsecase) Update(ctx context.Context, postID uint, userID uint, req *domain.UpdatePostRequest) (*domain.Post, error) {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, err
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

	err = u.postRepo.Update(ctx, post)
	return post, err
}

func (u *postUsecase) Delete(ctx context.Context, postID uint, userID uint) error {
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID != userID {
		return errors.New("unauthorized")
	}

	return u.postRepo.Delete(ctx, postID)
}

func (u *postUsecase) GetByID(ctx context.Context, id uint) (*domain.Post, error) {
	return u.postRepo.FindByID(ctx, id)
}

func (u *postUsecase) GetAll(ctx context.Context, page, limit int) ([]domain.Post, int64, error) {
	return u.postRepo.FindAll(ctx, page, limit)
}

func (u *postUsecase) GetByUserID(ctx context.Context, userID uint) ([]domain.Post, error) {
	return u.postRepo.FindByUserID(ctx, userID)
}
