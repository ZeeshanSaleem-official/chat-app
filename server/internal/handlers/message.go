package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/middleware"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
)

type MessageHandler struct {
	MessageRepo      *models.MessageRepository
	ConversationRepo *models.ConversationRepository
}

func NewMessageHandler(msgRepo *models.MessageRepository, convRepo *models.ConversationRepository) *MessageHandler {
	return &MessageHandler{
		MessageRepo:      msgRepo,
		ConversationRepo: convRepo,
	}
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conversationID := vars["id"]
	userID := middleware.GetUserID(r)

	// Check if user is a participant
	isParticipant, err := h.ConversationRepo.IsParticipant(conversationID, userID)
	if err != nil || !isParticipant {
		http.Error(w, `{"error":"Not authorized"}`, http.StatusForbidden)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	messages, err := h.MessageRepo.GetByConversation(conversationID, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch messages"}`, http.StatusInternalServerError)
		return
	}

	if messages == nil {
		messages = []models.Message{}
	}

	// Mark messages as read
	_ = h.MessageRepo.MarkAsRead(conversationID, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
