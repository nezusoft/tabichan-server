package utils

import "time"

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    string    `json:"user_id"`
	Device    string    `json:"device"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
