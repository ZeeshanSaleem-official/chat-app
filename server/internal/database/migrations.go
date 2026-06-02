package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,

		`CREATE TABLE IF NOT EXISTS users (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			username      VARCHAR(50) UNIQUE NOT NULL,
			email         VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			display_name  VARCHAR(100) NOT NULL,
			avatar_url    TEXT DEFAULT '',
			is_online     BOOLEAN DEFAULT false,
			last_seen     TIMESTAMP DEFAULT NOW(),
			created_at    TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS conversations (
			id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_at    TIMESTAMP DEFAULT NOW(),
			updated_at    TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS conversation_participants (
			conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
			user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
			joined_at       TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (conversation_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS messages (
			id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
			sender_id       UUID REFERENCES users(id) ON DELETE SET NULL,
			content         TEXT NOT NULL,
			message_type    VARCHAR(20) DEFAULT 'text',
			is_read         BOOLEAN DEFAULT false,
			created_at      TIMESTAMP DEFAULT NOW()
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender_id)`,
		`CREATE INDEX IF NOT EXISTS idx_participants_user ON conversation_participants(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	log.Println("✅ Database migrations completed")
	return nil
}
