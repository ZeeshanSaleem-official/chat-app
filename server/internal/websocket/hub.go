package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients mapped by user ID
	clients map[string]*Client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Inbound messages from clients to route
	broadcast chan *WSMessage

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *WSMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("👤 User %s connected (total: %d)", client.UserID, len(h.clients))

			// Notify other users that this user is online
			h.BroadcastToAll(&WSMessage{
				Type: "user_online",
				Data: map[string]interface{}{
					"user_id": client.UserID,
				},
			}, client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("👤 User %s disconnected (total: %d)", client.UserID, len(h.clients))

			// Notify other users that this user is offline
			h.BroadcastToAll(&WSMessage{
				Type: "user_offline",
				Data: map[string]interface{}{
					"user_id": client.UserID,
				},
			}, client.UserID)

		case message := <-h.broadcast:
			h.routeMessage(message)
		}
	}
}

func (h *Hub) routeMessage(msg *WSMessage) {
	if msg.RecipientID == "" {
		return
	}

	h.mu.RLock()
	client, ok := h.clients[msg.RecipientID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.send <- msg:
		default:
			// Client's send buffer is full, disconnect
			h.mu.Lock()
			delete(h.clients, client.UserID)
			close(client.send)
			h.mu.Unlock()
		}
	}
}

// SendToUser sends a message to a specific user if they're online
func (h *Hub) SendToUser(userID string, msg *WSMessage) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.send <- msg:
		default:
		}
	}
}

// BroadcastToAll sends a message to all connected clients except the sender
func (h *Hub) BroadcastToAll(msg *WSMessage, excludeUserID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, client := range h.clients {
		if userID != excludeUserID {
			select {
			case client.send <- msg:
			default:
			}
		}
	}
}

// IsUserOnline checks if a user is currently connected
func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}
