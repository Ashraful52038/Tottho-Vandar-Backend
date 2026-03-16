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
	CheckUserLikedComment(ctx context.Context, userID uint, commentID uint) (bool, error)
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
	// Validate user
	if err := u.validateUser(ctx, userID); err != nil {
		return nil, err
	}

	// Check if post exists
	post, err := u.postRepo.FindByID(ctx, postID)
	if err != nil || post == nil {
		return nil, errors.New("post not found")
	}

	// Check existing like
	existing, err := u.likeRepo.FindByUserAndPost(ctx, userID, postID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Unlike - delete like and decrement count
		err = u.likeRepo.Delete(ctx, userID, &postID, nil)
		if err != nil {
			return nil, err
		}

		// Decrement post likes count
		if post.Likes > 0 {
			post.Likes--
			_ = u.postRepo.Update(ctx, post)
		}

		return nil, nil
	}

	like := &domain.Like{
		UserID: userID,
		PostID: &postID,
	}

	err = u.likeRepo.Create(ctx, like)
	if err != nil {
		return nil, err
	}

	// Increment post likes count
	post.Likes++
	_ = u.postRepo.Update(ctx, post)

	return like, nil
}

// ToggleCommentLike
func (u *likeUsecase) ToggleCommentLike(ctx context.Context, userID uint, commentID uint) (*domain.Like, error) {
	// Validate user
	if err := u.validateUser(ctx, userID); err != nil {
		return nil, err
	}

	// Check if comment exists
	comment, err := u.commentRepo.FindByID(ctx, commentID)
	if err != nil || comment == nil {
		return nil, errors.New("comment not found")
	}

	// Check existing like
	existing, err := u.likeRepo.FindByUserAndComment(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Unlike
		err = u.likeRepo.DeleteByID(ctx, existing.ID)
		if err != nil {
			return nil, err
		}

		// Decrement comment likes count (if you have comment.Likes field)
		// comment.Likes--
		_ = u.commentRepo.Update(ctx, comment)

		return nil, nil
	}

	// Like
	like := &domain.Like{
		UserID:    userID,
		PostID:    nil,
		CommentID: &commentID,
	}

	err = u.likeRepo.Create(ctx, like)
	if err != nil {
		return nil, err
	}

	return like, nil
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

	likes, _, err := u.likeRepo.FindByUserID(ctx, userID, 0, 1000) // offset 0, limit 1000
	if err != nil {
		return nil, err
	}

	return likes, nil
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

func (u *likeUsecase) CheckUserLikedComment(ctx context.Context, userID uint, commentID uint) (bool, error) {
	like, err := u.likeRepo.FindByUserAndComment(ctx, userID, commentID)
	if err != nil {
		return false, err
	}
	return like != nil, nil
}
