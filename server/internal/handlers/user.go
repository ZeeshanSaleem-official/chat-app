package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/middleware"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
)

type UserHandler struct {
	UserRepo *models.UserRepository
}

func NewUserHandler(userRepo *models.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	user, err := h.UserRepo.FindByID(userID)
	if err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"error":"Search query is required"}`, http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	users, err := h.UserRepo.Search(query, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to search users"}`, http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []models.User{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
