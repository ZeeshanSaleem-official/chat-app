package models

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Joined fields for display
	OtherUserID     string  `json:"other_user_id,omitempty"`
	OtherUserName   string  `json:"other_user_name,omitempty"`
	OtherUserAvatar string  `json:"other_user_avatar,omitempty"`
	OtherUserOnline bool    `json:"other_user_online,omitempty"`
	LastMessage     *string `json:"last_message,omitempty"`
	LastMessageTime *time.Time `json:"last_message_time,omitempty"`
	UnreadCount     int     `json:"unread_count"`
}

type ConversationRepository struct {
	DB *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{DB: db}
}

// FindOrCreateDirect finds an existing 1:1 conversation between two users, or creates one
func (r *ConversationRepository) FindOrCreateDirect(userID1, userID2 string) (*Conversation, error) {
	// Check if a conversation already exists between these two users
	var convID string
	err := r.DB.QueryRow(
		`SELECT cp1.conversation_id
		 FROM conversation_participants cp1
		 JOIN conversation_participants cp2 ON cp1.conversation_id = cp2.conversation_id
		 WHERE cp1.user_id = $1 AND cp2.user_id = $2
		 LIMIT 1`,
		userID1, userID2,
	).Scan(&convID)

	if err == nil {
		// Conversation exists
		conv := &Conversation{ID: convID}
		return conv, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// Create new conversation
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var conv Conversation
	err = tx.QueryRow(
		`INSERT INTO conversations DEFAULT VALUES
		 RETURNING id, created_at, updated_at`,
	).Scan(&conv.ID, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Add both participants
	_, err = tx.Exec(
		`INSERT INTO conversation_participants (conversation_id, user_id) VALUES ($1, $2), ($1, $3)`,
		conv.ID, userID1, userID2,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &conv, nil
}

// GetUserConversations returns all conversations for a user with last message and other user info
func (r *ConversationRepository) GetUserConversations(userID string) ([]Conversation, error) {
	rows, err := r.DB.Query(
		`SELECT
			c.id, c.created_at, c.updated_at,
			other_user.id, other_user.display_name, other_user.avatar_url, other_user.is_online,
			last_msg.content, last_msg.created_at,
			COALESCE(unread.count, 0)
		 FROM conversations c
		 JOIN conversation_participants cp ON c.id = cp.conversation_id AND cp.user_id = $1
		 JOIN conversation_participants cp2 ON c.id = cp2.conversation_id AND cp2.user_id != $1
		 JOIN users other_user ON cp2.user_id = other_user.id
		 LEFT JOIN LATERAL (
			SELECT content, created_at FROM messages
			WHERE conversation_id = c.id
			ORDER BY created_at DESC LIMIT 1
		 ) last_msg ON true
		 LEFT JOIN LATERAL (
			SELECT COUNT(*) as count FROM messages
			WHERE conversation_id = c.id AND sender_id != $1 AND is_read = false
		 ) unread ON true
		 ORDER BY COALESCE(last_msg.created_at, c.updated_at) DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		if err := rows.Scan(
			&conv.ID, &conv.CreatedAt, &conv.UpdatedAt,
			&conv.OtherUserID, &conv.OtherUserName, &conv.OtherUserAvatar, &conv.OtherUserOnline,
			&conv.LastMessage, &conv.LastMessageTime,
			&conv.UnreadCount,
		); err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

// IsParticipant checks if a user is a participant of a conversation
func (r *ConversationRepository) IsParticipant(conversationID, userID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM conversation_participants WHERE conversation_id = $1 AND user_id = $2)`,
		conversationID, userID,
	).Scan(&exists)
	return exists, err
}

// GetOtherParticipant returns the other user's ID in a 1:1 conversation
func (r *ConversationRepository) GetOtherParticipant(conversationID, userID string) (string, error) {
	var otherID string
	err := r.DB.QueryRow(
		`SELECT user_id FROM conversation_participants
		 WHERE conversation_id = $1 AND user_id != $2 LIMIT 1`,
		conversationID, userID,
	).Scan(&otherID)
	return otherID, err
}
