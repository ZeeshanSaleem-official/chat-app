package models

import (
	"database/sql"
	"time"
)

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Content        string    `json:"content"`
	MessageType    string    `json:"message_type"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	// Joined fields
	SenderName   string `json:"sender_name,omitempty"`
	SenderAvatar string `json:"sender_avatar,omitempty"`
}

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) Create(conversationID, senderID, content, messageType string) (*Message, error) {
	msg := &Message{}
	err := r.DB.QueryRow(
		`INSERT INTO messages (conversation_id, sender_id, content, message_type)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, conversation_id, sender_id, content, message_type, is_read, created_at`,
		conversationID, senderID, content, messageType,
	).Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.Content,
		&msg.MessageType, &msg.IsRead, &msg.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Update conversation timestamp
	_, _ = r.DB.Exec(`UPDATE conversations SET updated_at = NOW() WHERE id = $1`, conversationID)

	return msg, nil
}

func (r *MessageRepository) GetByConversation(conversationID string, limit, offset int) ([]Message, error) {
	rows, err := r.DB.Query(
		`SELECT m.id, m.conversation_id, m.sender_id, m.content, m.message_type, m.is_read, m.created_at,
		        u.display_name, u.avatar_url
		 FROM messages m
		 LEFT JOIN users u ON m.sender_id = u.id
		 WHERE m.conversation_id = $1
		 ORDER BY m.created_at ASC
		 LIMIT $2 OFFSET $3`,
		conversationID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content,
			&m.MessageType, &m.IsRead, &m.CreatedAt, &m.SenderName, &m.SenderAvatar); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *MessageRepository) MarkAsRead(conversationID, userID string) error {
	_, err := r.DB.Exec(
		`UPDATE messages SET is_read = true
		 WHERE conversation_id = $1 AND sender_id != $2 AND is_read = false`,
		conversationID, userID,
	)
	return err
}
