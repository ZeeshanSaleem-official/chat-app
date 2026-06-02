package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"display_name"`
	AvatarURL    string    `json:"avatar_url"`
	IsOnline     bool      `json:"is_online"`
	LastSeen     time.Time `json:"last_seen"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(username, email, passwordHash, displayName string) (*User, error) {
	user := &User{}
	err := r.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash, display_name)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, username, email, display_name, avatar_url, is_online, last_seen, created_at`,
		username, email, passwordHash, displayName,
	).Scan(&user.ID, &user.Username, &user.Email, &user.DisplayName,
		&user.AvatarURL, &user.IsOnline, &user.LastSeen, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	user := &User{}
	err := r.DB.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, is_online, last_seen, created_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.DisplayName, &user.AvatarURL, &user.IsOnline, &user.LastSeen, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	user := &User{}
	err := r.DB.QueryRow(
		`SELECT id, username, email, display_name, avatar_url, is_online, last_seen, created_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email,
		&user.DisplayName, &user.AvatarURL, &user.IsOnline, &user.LastSeen, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Search(query string, currentUserID string) ([]User, error) {
	rows, err := r.DB.Query(
		`SELECT id, username, email, display_name, avatar_url, is_online, last_seen, created_at
		 FROM users
		 WHERE (username ILIKE $1 OR display_name ILIKE $1) AND id != $2
		 ORDER BY display_name
		 LIMIT 20`,
		"%"+query+"%", currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.DisplayName,
			&u.AvatarURL, &u.IsOnline, &u.LastSeen, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) SetOnlineStatus(userID string, online bool) error {
	_, err := r.DB.Exec(
		`UPDATE users SET is_online = $1, last_seen = NOW() WHERE id = $2`,
		online, userID,
	)
	return err
}
