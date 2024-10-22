package user

import "time"

type UserLogin struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	OAuthProvider string `json:"oauth"`
	ID            string `json:"user_id"`
}

type User struct {
	Username         string    `json:"username"`
	DisplayName      string    `json:"display_name"`
	UserID           string    `json:"user_id"`
	ProfileImageData string    `json:"profile_image_data"`
	CreatedAt        time.Time `json:"created_at"`
	LastLoginAt      time.Time `json:"last_login_at"`
}
