package usecase

import (
	"context"
	"errors"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type LikeUsecase interface {
	TogglePostLike(ctx context.Context, userID uint, postID uint) (*domain.Like, error)
	ToggleCommentLike(ctx context.Context, userID uint, commentID uint) (*domain.Like, error)
	GetPostLikes(ctx context.Context, postID uint) ([]domain.Like, error)
	GetPostLikesCount(ctx context.Context, postID uint) (int64, error)
	GetUserLikes(ctx context.Context, userID uint) ([]domain.Like, error)
	GetCommentLikes(ctx context.Context, commentID uint) ([]domain.Like, error)
	GetCommentLikesCount(ctx context.Context, commentID uint) (int64, error)
}

type likeUsecase struct {
	likeRepo    repository.LikeRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

func NewLikeUsecase(
	likeRepo repository.LikeRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
) LikeUsecase {
	return &likeUsecase{
		likeRepo:    likeRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

// TogglePostLike
func (u *likeUsecase) TogglePostLike(ctx context.Context, userID uint, postID uint) (*domain.Like, error) {
	if err := u.validateUser(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := u.postRepo.FindByID(ctx, postID); err != nil {
		return nil, errors.New("post not found")
	}

	existing, _ := u.likeRepo.FindByUserAndPost(ctx, userID, postID)
	if existing != nil {
		return nil, u.likeRepo.Delete(ctx, userID, &postID, nil)
	}

	like := &domain.Like{UserID: userID, PostID: &postID}
	return like, u.likeRepo.Create(ctx, like)
}

// ToggleCommentLike
func (u *likeUsecase) ToggleCommentLike(ctx context.Context, userID uint, commentID uint) (*domain.Like, error) {
	if err := u.validateUser(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := u.commentRepo.FindByID(ctx, commentID); err != nil {
		return nil, errors.New("comment not found")
	}

	existing, _ := u.likeRepo.FindByUserAndComment(ctx, userID, commentID)
	if existing != nil {
		return nil, u.likeRepo.Delete(ctx, userID, nil, &commentID)
	}

	like := &domain.Like{UserID: userID, CommentID: &commentID}
	return like, u.likeRepo.Create(ctx, like)
}

// GetPostLikes
func (u *likeUsecase) GetPostLikes(ctx context.Context, postID uint) ([]domain.Like, error) {
	if _, err := u.postRepo.FindByID(ctx, postID); err != nil {
		return nil, errors.New("post not found")
	}
	return u.likeRepo.FindByPostID(ctx, postID)
}

// GetPostLikesCount
func (u *likeUsecase) GetPostLikesCount(ctx context.Context, postID uint) (int64, error) {
	return u.likeRepo.CountByPostID(ctx, postID)
}

// GetUserLikes
func (u *likeUsecase) GetUserLikes(ctx context.Context, userID uint) ([]domain.Like, error) {
	if err := u.validateUser(ctx, userID); err != nil {
		return nil, err
	}
	return u.likeRepo.FindByUserID(ctx, userID)
}

// validateUser
func (u *likeUsecase) validateUser(ctx context.Context, userID uint) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	if !user.Verified {
		return errors.New("email not verified")
	}
	return nil
}

// GetCommentLikes
func (u *likeUsecase) GetCommentLikes(ctx context.Context, commentID uint) ([]domain.Like, error) {
	if _, err := u.commentRepo.FindByID(ctx, commentID); err != nil {
		return nil, errors.New("comment not found")
	}
	return u.likeRepo.FindByCommentID(ctx, commentID)
}

// GetCommentLikesCount
func (u *likeUsecase) GetCommentLikesCount(ctx context.Context, commentID uint) (int64, error) {
	return u.likeRepo.CountByCommentID(ctx, commentID)
}
