package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) VerifyEmail(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("verified", true).Error
}

func (r *userRepository) FindByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("verification_token = ?", token).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByResetToken(ctx context.Context, token string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}

func (r *userRepository) SetResetToken(ctx context.Context, userID uint, token string, expiry time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"reset_token":        token,
			"reset_token_expiry": expiry,
		}).Error
}

func (r *userRepository) FindMostFollowed(ctx context.Context, limit int) ([]domain.UserWithFollowCount, error) {
	var users []domain.UserWithFollowCount
	err := r.db.WithContext(ctx).
		Table("users").
		Select("users.id, users.name, users.avatar, users.bio, COUNT(f.follower_id) as followers_count").
		Joins("LEFT JOIN followers f ON f.following_id = users.id AND f.status = 'accepted'").
		Where("users.deleted_at IS NULL").
		Group("users.id").
		Order("followers_count DESC, users.created_at ASC").
		Limit(limit).
		Scan(&users).Error
	return users, err
}
