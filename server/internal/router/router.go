package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/config"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/handlers"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/middleware"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
	ws "github.com/ZeeshanSaleem-official/chat-app/server/internal/websocket"
)

func Setup(db *sql.DB, cfg *config.Config, hub *ws.Hub) http.Handler {
	r := mux.NewRouter()

	// Repositories
	userRepo := models.NewUserRepository(db)
	msgRepo := models.NewMessageRepository(db)
	convRepo := models.NewConversationRepository(db)

	// Handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userRepo)
	msgHandler := handlers.NewMessageHandler(msgRepo, convRepo)
	convHandler := handlers.NewConversationHandler(convRepo)
	wsHandler := ws.NewHandler(hub, cfg.JWTSecret, msgRepo, convRepo, userRepo)

	// Auth middleware
	authMW := middleware.AuthMiddleware(cfg.JWTSecret)

	// Public routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMW)

	api.HandleFunc("/users/me", userHandler.GetProfile).Methods("GET")
	api.HandleFunc("/users/search", userHandler.SearchUsers).Methods("GET")
	api.HandleFunc("/conversations", convHandler.ListConversations).Methods("GET")
	api.HandleFunc("/conversations", convHandler.CreateConversation).Methods("POST")
	api.HandleFunc("/conversations/{id}/messages", msgHandler.GetMessages).Methods("GET")

	// WebSocket (auth handled inside the handler via token query param)
	r.HandleFunc("/ws", wsHandler.ServeWS)

	// Apply CORS middleware
	handler := middleware.CORSMiddleware(cfg.FrontendURL)(r)

	return handler
}
