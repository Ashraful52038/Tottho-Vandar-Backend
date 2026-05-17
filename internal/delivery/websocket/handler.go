package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	userID := c.QueryParam("user_id")
	log.Printf("🔌 New WebSocket connection from user: %s", userID)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("❌ Upgrade failed: %v", err)
		return err
	}

	client := &Client{
		ID:     generateClientID(),
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    h.hub,
	}

	h.hub.register <- client

	// ✅ IMPORTANT: Unregister when function returns
	defer func() {
		h.hub.unregister <- client
		conn.Close()
		log.Printf("🔌 WebSocket closed for user: %s", userID)
	}()

	// ✅ Write pump (sends messages to client)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case message, ok := <-client.Send:
				if !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("❌ Write error for user %s: %v", userID, err)
					return
				}
				log.Printf("📤 Sent to user %s: %s", userID, string(message))
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("❌ Ping error for user %s: %v", userID, err)
					return
				}
			}
		}
	}()

	// ✅ Read pump - keeps connection alive (MUST be in main goroutine)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("⚠️ Read error for user %s: %v", userID, err)
			break
		}
	}

	return nil
}

func generateClientID() string {
	return "ws_" + time.Now().Format("20060102150405")
}
