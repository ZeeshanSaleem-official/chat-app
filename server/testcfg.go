package main

import (
	"fmt"
	"github.com/ZeeshanSaleem-official/chat-app/server/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	cfg := config.Load()
	fmt.Printf("DBUser: '%s'\n", cfg.DBUser)
	fmt.Printf("DBPassword: '%s'\n", cfg.DBPassword)
	fmt.Printf("DatabaseURL: '%s'\n", cfg.DatabaseURL())
}
