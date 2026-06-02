package main

import (
	"log"

	"github.com/ZeeshanSaleem-official/chat-app/server/internal/config"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/database"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Try loading from different relative paths
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")
	
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("❌ DB connect error: %v", err)
	}
	defer db.Close()

	userRepo := models.NewUserRepository(db)

	dummyUsers := []struct {
		Username    string
		Email       string
		DisplayName string
		Password    string
	}{
		{"alice", "alice@example.com", "Alice Smith", "password123"},
		{"bob", "bob@example.com", "Bob Johnson", "password123"},
		{"charlie", "charlie@example.com", "Charlie Brown", "password123"},
		{"david", "david@example.com", "David Miller", "password123"},
		{"emma", "emma@example.com", "Emma Watson", "password123"},
		{"frank", "frank@example.com", "Frank Castle", "password123"},
		{"grace", "grace@example.com", "Grace Hopper", "password123"},
		{"zeeshan", "zeeshan@example.com", "Zeeshan Saleem", "password123"}, // For you!
	}

	for _, u := range dummyUsers {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		_, err := userRepo.Create(u.Username, u.Email, string(hash), u.DisplayName)
		if err != nil {
			log.Printf("⚠️ User %s might already exist (or error): %v", u.Username, err)
		} else {
			log.Printf("✅ Created user: %s (%s)", u.DisplayName, u.Email)
		}
	}
	log.Println("🎉 Database seeding complete!")
}
