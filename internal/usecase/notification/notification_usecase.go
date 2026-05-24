package notification

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Ashraful52038/tottho-vandar-backend/internal/delivery/websocket"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/domain"
	"github.com/Ashraful52038/tottho-vandar-backend/internal/repository"
)

type NotificationUsecase struct {
	wsHub            *websocket.Hub
	notificationRepo repository.NotificationRepository
}

func NewNotificationUsecase(hub *websocket.Hub, notificationRepo repository.NotificationRepository) *NotificationUsecase {
	return &NotificationUsecase{wsHub: hub, notificationRepo: notificationRepo}
}

func (n *NotificationUsecase) SendNotification(userID string, title, message string, eventType domain.WebSocketEventType, data interface{}) {
	log.Printf("🔔 [SEND] userID=%s, title=%s", userID, title)

	if n.notificationRepo != nil {
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			notif := &domain.Notification{
				UserID:  uint(userIDUint),
				Type:    domain.NotificationType(eventType),
				Subject: title,
				Content: message,
				IsRead:  false,
			}
			if err := n.notificationRepo.Create(nil, notif); err != nil {
				log.Printf("❌ DB save failed: %v", err)
			} else {
				log.Printf("✅ Notification saved to DB for user: %s", userID)
			}
		}
	}

	notification := domain.WebSocketMessage{
		Type:      eventType,
		Title:     title,
		Message:   message,
		Data:      data,
		UserID:    userID,
		CreatedAt: time.Now().Unix(),
	}

	jsonData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("❌ [ERROR] Failed to marshal: %v", err)
		return
	}

	log.Printf("📤 [JSON] %s", string(jsonData))

	result := n.wsHub.SendToUser(userID, jsonData)
	if result {
		log.Printf("✅ [SUCCESS] Notification sent to user %s", userID)
	} else {
		log.Printf("❌ [FAILED] Could not send to user %s", userID)
	}
}

// NotifyNewPost - accepts uint and converts to string
func (n *NotificationUsecase) NotifyNewPost(userID uint, postID uint, postTitle string) {
	n.SendNotification(
		strconv.FormatUint(uint64(userID), 10),
		"New Post Published",
		"Your post '"+postTitle+"' is now live",
		domain.EventNewPost,
		map[string]interface{}{"postId": postID},
	)
}

// NotifyNewLike - accepts uint for userID and postID
func (n *NotificationUsecase) NotifyNewLike(postAuthorID uint, postID uint, likerName string) {
	n.SendNotification(
		strconv.FormatUint(uint64(postAuthorID), 10),
		"New Like",
		likerName+" liked your post",
		domain.EventNewLike,
		map[string]interface{}{"postId": postID},
	)
}

// NotifyNewComment - accepts uint for userID and postID
func (n *NotificationUsecase) NotifyNewComment(postAuthorID uint, postID uint, commenterName, comment string) {
	n.SendNotification(
		strconv.FormatUint(uint64(postAuthorID), 10),
		"New Comment",
		commenterName+" commented: "+comment,
		domain.EventNewComment,
		map[string]interface{}{"postId": postID},
	)
}

// NotifyNewFollow - accepts uint for userID
func (n *NotificationUsecase) NotifyNewFollow(followedUserID uint, followerName string) {
	n.SendNotification(
		strconv.FormatUint(uint64(followedUserID), 10),
		"New Follower",
		followerName+" started following you",
		domain.EventNewFollow,
		map[string]interface{}{"follower": followerName},
	)
}
