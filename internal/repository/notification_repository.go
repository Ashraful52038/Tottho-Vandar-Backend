package repository

import (
	"context"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]domain.Notification, int64, error)
	GetUnreadByUserID(ctx context.Context, userID uint) ([]domain.Notification, error)
	CountUnread(ctx context.Context, userID uint) (int64, error)
	GetByID(ctx context.Context, id uint) (*domain.Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	MarkAllAsRead(ctx context.Context, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	DeleteAll(ctx context.Context, userID uint) error
}
