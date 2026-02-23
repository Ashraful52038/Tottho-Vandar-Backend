package usecase

import (
	"context"
	"errors"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type CommentUsecase interface {
	Create(ctx context.Context, userID uint, postID uint, content string) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]domain.Comment, error)
	Update(ctx context.Context, commentID uint, userID uint, content string) (*domain.Comment, error)
	Delete(ctx context.Context, commentID uint, userID uint) error
}

type commentUsecase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
}

func NewCommentUsecase(
	commentRepo repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) CommentUsecase {
	return &commentUsecase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
	}
}

func (u *commentUsecase) Create(ctx context.Context, userID uint, postID uint, content string) (*domain.Comment, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if !user.Verified {
		return nil, errors.New("email not verified")
	}

	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}

	comment := &domain.Comment{
		Content: content,
		PostID:  postID,
		UserID:  userID,
	}

	err = u.commentRepo.Create(ctx, comment)
	return comment, err
}

func (u *commentUsecase) GetByPostID(ctx context.Context, postID uint) ([]domain.Comment, error) {
	_, err := u.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return u.commentRepo.FindByPostID(ctx, postID)
}

func (u *commentUsecase) Update(ctx context.Context, commentID uint, userID uint, content string) (*domain.Comment, error) {
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil || comment == nil {
		return nil, errors.New("comment not found")
	}

	if comment.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	comment.Content = content
	err = u.commentRepo.Update(ctx, comment)
	return comment, err
}

func (u *commentUsecase) Delete(ctx context.Context, commentID uint, userID uint) error {
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil || comment == nil {
		return errors.New("comment not found")
	}

	if comment.UserID != userID {
		return errors.New("unauthorized")
	}

	return u.commentRepo.Delete(ctx, commentID)
}
