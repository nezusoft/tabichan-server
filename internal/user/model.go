package user

type User struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	OAuthProvider string `json:"oauth"`
	ID            string `json:"userId"`
}

type UserDetails struct {
	Username         string `json:"username"`
	DisplayName      string `json:"displayName"`
	UserID           string `json:"userId"`
	ProfileImageData string `json:"profileImageData"`
}
