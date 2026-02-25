// internal/domain/notification.go
package domain

import (
	"time"
)

type NotificationType string

const (
	NotificationTypeCommentReply  NotificationType = "comment_reply"
	NotificationTypePostReply     NotificationType = "post_reply"
	NotificationTypeWelcome       NotificationType = "welcome"
	NotificationTypePasswordReset NotificationType = "password_reset"
)

type Notification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	UserID    uint             `json:"userId" gorm:"not null"`
	User      User             `json:"user,omitempty"`
	Type      NotificationType `json:"type"`
	Subject   string           `json:"subject"`
	Content   string           `json:"content"`
	IsRead    bool             `json:"isRead" gorm:"default:false"`
	CreatedAt time.Time        `json:"createdAt"`
}

type UserNotificationPreference struct {
	UserID         uint `json:"userId" gorm:"primaryKey"`
	EmailReplies   bool `json:"emailReplies" gorm:"default:true"`
	EmailMarketing bool `json:"emailMarketing" gorm:"default:false"`
	PushEnabled    bool `json:"pushEnabled" gorm:"default:false"`
}
