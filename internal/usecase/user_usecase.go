package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type UserUsecase interface {
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error)
	GetUserPosts(ctx context.Context, userID string, page, limit int) ([]domain.Post, int64, error)
	GetUserComments(ctx context.Context, userID string, page, limit int) ([]domain.Comment, int64, error)
	GetUserLikes(ctx context.Context, userID string, page, limit int) ([]domain.Like, int64, error)
	GetFollowers(ctx context.Context, userID string, currentUserID uint, page, limit int) ([]domain.FollowUser, int64, error)
	GetFollowing(ctx context.Context, userID string, currentUserID uint, page, limit int) ([]domain.FollowUser, int64, error)
	CheckFollowStatus(ctx context.Context, currentUserID, targetUserID uint) (bool, error)
	FollowUser(ctx context.Context, currentUserID, targetUserID uint) error
	UnfollowUser(ctx context.Context, currentUserID, targetUserID uint) error
	GetMostFollowedUsers(ctx context.Context, limit int) ([]domain.UserWithFollowCount, error)
}

type userUsecase struct {
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
	likeRepo    repository.LikeRepository
	followRepo  repository.FollowRepository
}

func NewUserUsecase(
	userRepo repository.UserRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	followRepo repository.FollowRepository,
) UserUsecase {
	return &userUsecase{
		userRepo:    userRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		likeRepo:    likeRepo,
		followRepo:  followRepo,
	}
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) Update(ctx context.Context, id uint, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}

	err = u.userRepo.Update(ctx, user)
	return user, err
}

func (u *userUsecase) GetUserProfile(ctx context.Context, userID uint) (*domain.UserProfileResponse, error) {
	if userID == 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	postsCount, _ := u.postRepo.CountByUserID(ctx, userID)
	commentsCount, _ := u.commentRepo.CountByUserID(ctx, userID)
	likesCount, _ := u.likeRepo.CountByUserID(ctx, userID)
	followersCount, _ := u.followRepo.GetFollowersCount(ctx, userID)
	followingCount, _ := u.followRepo.GetFollowingCount(ctx, userID)

	bio := ""
	if user.Bio != nil {
		bio = *user.Bio
	}

	avatar := ""
	if user.Avatar != nil {
		avatar = *user.Avatar
	}

	response := &domain.UserProfileResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    avatar,
		Bio:       bio,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		Stats: domain.UserStats{
			Posts:     postsCount,
			Comments:  commentsCount,
			Likes:     likesCount,
			Followers: followersCount,
			Following: followingCount,
		},
	}

	return response, nil
}

func (u *userUsecase) GetUserPosts(ctx context.Context, userID string, page, limit int) ([]domain.Post, int64, error) {
	// userID string to uint conversion
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	_, err = u.userRepo.FindByID(ctx, uint(id))
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	return u.postRepo.FindByUserID(ctx, uint(id), offset, limit)
}

func (u *userUsecase) GetUserComments(ctx context.Context, userID string, page, limit int) ([]domain.Comment, int64, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	return u.commentRepo.FindByUserID(ctx, uint(id), offset, limit)
}

// GetUserLikes
func (u *userUsecase) GetUserLikes(ctx context.Context, userID string, page, limit int) ([]domain.Like, int64, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	return u.likeRepo.FindByUserID(ctx, uint(id), offset, limit)
}

func (u *userUsecase) GetFollowers(ctx context.Context, userID string, currentUserID uint, page, limit int) ([]domain.FollowUser, int64, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	_, err = u.userRepo.FindByID(ctx, uint(id))
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	followers, total, err := u.followRepo.GetFollowers(ctx, uint(id), offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for i := range followers {
		isFollowing, _ := u.followRepo.CheckFollowStatus(ctx, currentUserID, followers[i].ID)
		followers[i].IsFollowing = isFollowing
	}

	return followers, total, nil
}

func (u *userUsecase) GetFollowing(ctx context.Context, userID string, currentUserID uint, page, limit int) ([]domain.FollowUser, int64, error) {
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	following, total, err := u.followRepo.GetFollowing(ctx, uint(id), offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for i := range following {
		isFollowing, _ := u.followRepo.CheckFollowStatus(ctx, currentUserID, following[i].ID)
		following[i].IsFollowing = isFollowing
	}

	return following, total, nil
}

func (u *userUsecase) CheckFollowStatus(ctx context.Context, currentUserID, targetUserID uint) (bool, error) {
	return u.followRepo.CheckFollowStatus(ctx, currentUserID, targetUserID)
}

func (u *userUsecase) FollowUser(ctx context.Context, currentUserID, targetUserID uint) error {

	if currentUserID == 0 || targetUserID == 0 {
		return errors.New("invalid user ids")
	}

	if currentUserID == targetUserID {
		return errors.New("you cannot follow yourself")
	}

	targetUser, err := u.userRepo.FindByID(ctx, targetUserID)
	if err != nil || targetUser == nil {
		return errors.New("target user not found")
	}

	isFollowing, err := u.followRepo.CheckFollowStatus(ctx, currentUserID, targetUserID)
	if err != nil {
		return err
	}

	if isFollowing {
		return errors.New("already following this user")
	}

	return u.followRepo.Follow(ctx, currentUserID, targetUserID)
}

func (u *userUsecase) UnfollowUser(ctx context.Context, currentUserID, targetUserID uint) error {

	if currentUserID == 0 || targetUserID == 0 {
		return errors.New("invalid user ids")
	}

	isFollowing, err := u.followRepo.CheckFollowStatus(ctx, currentUserID, targetUserID)
	if err != nil {
		return err
	}
	if !isFollowing {
		return errors.New("you are not following this user")
	}

	return u.followRepo.Unfollow(ctx, currentUserID, targetUserID)
}

func (u *userUsecase) GetMostFollowedUsers(ctx context.Context, limit int) ([]domain.UserWithFollowCount, error) {
	return u.userRepo.FindMostFollowed(ctx, limit)
}
