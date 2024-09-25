package user

type User struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	OAuthProvider string `json:"oauth"`
	ID            string `json:"userId"`
}
