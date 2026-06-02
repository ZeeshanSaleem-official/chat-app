package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/config"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/database"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/router"
	ws "github.com/ZeeshanSaleem-official/chat-app/server/internal/websocket"
)

func main() {
	// Load .env from project root (../.env when CWD is server/)
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := config.Load()

	// Connect to PostgreSQL
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Start WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Setup router
	handler := router.Setup(db, cfg, hub)

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler,
	}

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("🛑 Shutting down server...")
		server.Close()
	}()

	log.Printf("🚀 Server starting on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v", err)
	}
}
