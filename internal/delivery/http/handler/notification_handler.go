package handler

import (
	"net/http"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/usecase/notification"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	notifUsecase *notification.NotificationUsecase
	notifRepo    repository.NotificationRepository
}

func NewNotificationHandler(
	notifUsecase *notification.NotificationUsecase,
	notifRepo repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{
		notifUsecase: notifUsecase,
		notifRepo:    notifRepo,
	}
}

// GetMyNotifications - ইউজারের সব নোটিফিকেশন দেখাবে
func (h *NotificationHandler) GetMyNotifications(c echo.Context) error {
	// userID ইউজ না করলে _ (underscore) দিয়ে ইগনোর করো
	userID := c.Get("userID").(uint)
	_ = userID // temporarily ignore

	return c.JSON(http.StatusOK, map[string]interface{}{
		"notifications": []interface{}{},
		"total":         0,
	})
}

// GetUnreadCount - আনরিড নোটিফিকেশনের কাউন্ট
func (h *NotificationHandler) GetUnreadCount(c echo.Context) error {
	userID := c.Get("userID").(uint)
	_ = userID // temporarily ignore

	return c.JSON(http.StatusOK, map[string]interface{}{
		"unread_count": 0,
	})
}

// MarkAsRead - একটি নোটিফিকেশন রিড হিসেবে মার্ক করো
func (h *NotificationHandler) MarkAsRead(c echo.Context) error {
	id := c.Param("id")
	_ = id // temporarily ignore

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// MarkAllAsRead - সব নোটিফিকেশন রিড হিসেবে মার্ক করো
func (h *NotificationHandler) MarkAllAsRead(c echo.Context) error {
	userID := c.Get("userID").(uint)
	_ = userID // temporarily ignore

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// DeleteNotification - একটি নোটিফিকেশন ডিলিট করো
func (h *NotificationHandler) DeleteNotification(c echo.Context) error {
	id := c.Param("id")
	userID := c.Get("userID").(uint)
	_ = id     // temporarily ignore
	_ = userID // temporarily ignore

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
