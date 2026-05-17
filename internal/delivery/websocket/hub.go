package websocket

import (
	"log"
	"sync"
)

type Client struct {
	ID     string
	UserID string
	Send   chan []byte
	Room   string
	Hub    *Hub
}

type Room struct {
	Name    string
	Clients map[string]*Client
	mu      sync.RWMutex
}

type Hub struct {
	clientsByUser map[string]*Client
	clientsByID   map[string]*Client
	clients       map[string]*Client
	register      chan *Client
	unregister    chan *Client
	broadcast     chan []byte

	rooms     map[string]*Room
	joinRoom  chan JoinRoomRequest
	leaveRoom chan LeaveRoomRequest

	mu sync.RWMutex
}

type JoinRoomRequest struct {
	Client   *Client
	RoomName string
}

type LeaveRoomRequest struct {
	Client   *Client
	RoomName string
}

func NewHub() *Hub {
	return &Hub{
		clientsByUser: make(map[string]*Client),
		clientsByID:   make(map[string]*Client),
		clients:       make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		broadcast:     make(chan []byte),
		rooms:         make(map[string]*Room),
		joinRoom:      make(chan JoinRoomRequest),
		leaveRoom:     make(chan LeaveRoomRequest),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if oldClient, exists := h.clientsByUser[client.UserID]; exists {
				log.Printf("⚠️ User %s already connected, closing old connection", client.UserID)
				close(oldClient.Send)
				delete(h.clientsByID, oldClient.ID)
				delete(h.clientsByUser, oldClient.UserID)
			}
			h.clientsByID[client.ID] = client
			if client.UserID != "" {
				h.clientsByUser[client.UserID] = client
				log.Printf("✅ Client registered: UserID=%s", client.UserID)
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				if client.UserID != "" {
					delete(h.clientsByUser, client.UserID)
				}
				delete(h.clientsByID, client.ID)
				close(client.Send)
				log.Printf("❌ WebSocket client unregistered: %s", client.ID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
			}
			h.mu.RUnlock()

		case req := <-h.joinRoom:
			h.handleJoinRoom(req)

		case req := <-h.leaveRoom:
			h.handleLeaveRoom(req)
		}
	}
}

func (h *Hub) handleJoinRoom(req JoinRoomRequest) {
	h.mu.Lock()
	room, exists := h.rooms[req.RoomName]
	if !exists {
		room = &Room{
			Name:    req.RoomName,
			Clients: make(map[string]*Client),
		}
		h.rooms[req.RoomName] = room
	}
	h.mu.Unlock()

	room.mu.Lock()
	room.Clients[req.Client.ID] = req.Client
	req.Client.Room = req.RoomName
	room.mu.Unlock()

	log.Printf("🚪 Client %s joined room: %s", req.Client.ID, req.RoomName)
}

func (h *Hub) handleLeaveRoom(req LeaveRoomRequest) {
	h.mu.RLock()
	room, exists := h.rooms[req.RoomName]
	h.mu.RUnlock()

	if exists {
		room.mu.Lock()
		delete(room.Clients, req.Client.ID)
		if req.Client.Room == req.RoomName {
			req.Client.Room = ""
		}
		room.mu.Unlock()
		log.Printf("🚪 Client %s left room: %s", req.Client.ID, req.RoomName)
	}
}

// Send notification to specific user
func (h *Hub) SendToUser(userID string, message []byte) bool {
	h.mu.RLock()
	client, exists := h.clientsByUser[userID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("⚠️ User %s not connected", userID)
		return false
	}

	select {
	case client.Send <- message:
		log.Printf("📨 Notification sent to user: %s", userID)
		return true
	default:
		log.Printf("⚠️ Failed to send to user %s", userID)
		return false
	}
}

// Send to room
func (h *Hub) SendToRoom(roomName string, message []byte) int {
	h.mu.RLock()
	room, exists := h.rooms[roomName]
	h.mu.RUnlock()

	if !exists {
		return 0
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	sentCount := 0
	for _, client := range room.Clients {
		select {
		case client.Send <- message:
			sentCount++
		default:
		}
	}
	return sentCount
}
