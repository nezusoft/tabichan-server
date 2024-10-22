package utils

type Session struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Device    string `json:"device"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
