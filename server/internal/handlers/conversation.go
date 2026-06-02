package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/middleware"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
)

type ConversationHandler struct {
	ConversationRepo *models.ConversationRepository
}

func NewConversationHandler(convRepo *models.ConversationRepository) *ConversationHandler {
	return &ConversationHandler{ConversationRepo: convRepo}
}

type CreateConversationRequest struct {
	UserID string `json:"user_id"`
}

func (h *ConversationHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	conversations, err := h.ConversationRepo.GetUserConversations(userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch conversations"}`, http.StatusInternalServerError)
		return
	}

	if conversations == nil {
		conversations = []models.Conversation{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	var req CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	currentUserID := middleware.GetUserID(r)
	if req.UserID == "" || req.UserID == currentUserID {
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	conversation, err := h.ConversationRepo.FindOrCreateDirect(currentUserID, req.UserID)
	if err != nil {
		http.Error(w, `{"error":"Failed to create conversation"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(conversation)
}
