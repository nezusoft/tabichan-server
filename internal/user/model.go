package user

import (
	"time"

	"github.com/tabichanorg/tabichan-server/internal/utils"
)

type UserLogin struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	OAuthProvider string `json:"oauth"`
	UserID        string `json:"user_id"`
}

type User struct {
	Username         string    `json:"username"`
	DisplayName      string    `json:"displayName"`
	UserID           string    `json:"userId"`
	ProfileImageData string    `json:"profileImageData"`
	CreatedAt        time.Time `json:"createdAt"`
	LastLoginAt      time.Time `json:"lastLoginAt"`
}

type LoginRequestResponse struct {
	Token   string         `json:"token"`
	Session *utils.Session `json:"session"`
}
