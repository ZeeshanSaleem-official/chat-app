package websocket

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	ws "github.com/gorilla/websocket"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type           string                 `json:"type"`
	Data           map[string]interface{} `json:"data,omitempty"`
	SenderID       string                 `json:"sender_id,omitempty"`
	RecipientID    string                 `json:"recipient_id,omitempty"`
	ConversationID string                 `json:"conversation_id,omitempty"`
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, validate origin
	},
}

type Handler struct {
	hub              *Hub
	jwtSecret        string
	messageRepo      *models.MessageRepository
	conversationRepo *models.ConversationRepository
	userRepo         *models.UserRepository
}

func NewHandler(hub *Hub, jwtSecret string, msgRepo *models.MessageRepository, convRepo *models.ConversationRepository, userRepo *models.UserRepository) *Handler {
	return &Handler{
		hub:              hub,
		jwtSecret:        jwtSecret,
		messageRepo:      msgRepo,
		conversationRepo: convRepo,
		userRepo:         userRepo,
	}
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Authenticate via query parameter token
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		// Also check Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse and validate JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(h.hub, conn, userID)
	h.hub.register <- client

	// Set user online
	_ = h.userRepo.SetOnlineStatus(userID, true)

	// Start read/write pumps
	go client.WritePump()
	go client.ReadPump(h.handleMessage)
}

func (h *Handler) handleMessage(client *Client, msg *WSMessage) {
	switch msg.Type {
	case "chat_message":
		h.handleChatMessage(client, msg)
	case "typing":
		h.handleTyping(client, msg)
	case "read_receipt":
		h.handleReadReceipt(client, msg)
	}
}

func (h *Handler) handleChatMessage(client *Client, msg *WSMessage) {
	content, _ := msg.Data["content"].(string)
	conversationID := msg.ConversationID

	if content == "" || conversationID == "" {
		return
	}

	// Verify user is a participant
	isParticipant, err := h.conversationRepo.IsParticipant(conversationID, client.UserID)
	if err != nil || !isParticipant {
		return
	}

	// Persist message to DB
	savedMsg, err := h.messageRepo.Create(conversationID, client.UserID, content, "text")
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	// Get sender info
	sender, _ := h.userRepo.FindByID(client.UserID)
	senderName := ""
	senderAvatar := ""
	if sender != nil {
		senderName = sender.DisplayName
		senderAvatar = sender.AvatarURL
	}

	// Get other participant
	otherUserID, err := h.conversationRepo.GetOtherParticipant(conversationID, client.UserID)
	if err != nil {
		return
	}

	// Build response message
	responseMsg := &WSMessage{
		Type:           "chat_message",
		ConversationID: conversationID,
		SenderID:       client.UserID,
		Data: map[string]interface{}{
			"id":            savedMsg.ID,
			"content":       savedMsg.Content,
			"message_type":  savedMsg.MessageType,
			"is_read":       savedMsg.IsRead,
			"created_at":    savedMsg.CreatedAt,
			"sender_name":   senderName,
			"sender_avatar": senderAvatar,
		},
	}

	// Send to recipient
	h.hub.SendToUser(otherUserID, responseMsg)

	// Send confirmation back to sender
	h.hub.SendToUser(client.UserID, responseMsg)
}

func (h *Handler) handleTyping(client *Client, msg *WSMessage) {
	conversationID := msg.ConversationID
	if conversationID == "" {
		return
	}

	otherUserID, err := h.conversationRepo.GetOtherParticipant(conversationID, client.UserID)
	if err != nil {
		return
	}

	h.hub.SendToUser(otherUserID, &WSMessage{
		Type:           "typing",
		ConversationID: conversationID,
		SenderID:       client.UserID,
	})
}

func (h *Handler) handleReadReceipt(client *Client, msg *WSMessage) {
	conversationID := msg.ConversationID
	if conversationID == "" {
		return
	}

	_ = h.messageRepo.MarkAsRead(conversationID, client.UserID)

	otherUserID, err := h.conversationRepo.GetOtherParticipant(conversationID, client.UserID)
	if err != nil {
		return
	}

	h.hub.SendToUser(otherUserID, &WSMessage{
		Type:           "read_receipt",
		ConversationID: conversationID,
		SenderID:       client.UserID,
	})
}
