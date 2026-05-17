package domain

type WebSocketEventType string

const (
	EventNewPost      WebSocketEventType = "new_post"
	EventNewLike      WebSocketEventType = "new_like"
	EventNewComment   WebSocketEventType = "new_comment"
	EventNewFollow    WebSocketEventType = "new_follow"
	EventNotification WebSocketEventType = "notification"
)

type WebSocketMessage struct {
	Type      WebSocketEventType `json:"type"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	Data      interface{}        `json:"data"`
	UserID    string             `json:"userId,omitempty"`
	CreatedAt int64              `json:"createdAt"`
}
