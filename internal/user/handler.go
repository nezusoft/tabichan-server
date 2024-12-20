package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type UserHandler struct {
	Service *UserService
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var newUser UserLogin
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := validateUsername(newUser.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateEmail(newUser.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.Service.Signup(newUser, r.Header.Get("User-Agent"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, response.Session.ExpiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    response.Session.SessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		// Secure:   true, // TODO: use Secure when hosting HTTPS
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		UsernameOrEmail string `json:"UsernameOrEmail"`
		Password        string `json:"password"`
		RememberMe      bool   `json:"RememberMe"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response, err := h.Service.Login(loginRequest.UsernameOrEmail, loginRequest.Password, r.Header.Get("User-Agent"), loginRequest.RememberMe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	expiresAt, err := time.Parse(time.RFC3339, response.Session.ExpiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionID",
		Value:    response.Session.SessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		// Secure:   true, // TODO: use Secure when hosting HTTPS
	})

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	userDetails, err := h.Service.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(userDetails)
}

func validateUsername(username string) error {
	for _, char := range username {
		if unicode.IsSpace(char) {
			return fmt.Errorf("username cannot contain spaces")
		}
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return fmt.Errorf("username can only contain letters and numbers")
		}
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return nil
	}

	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return fmt.Errorf("email must contain '@'")
	}

	if atIndex == 0 {
		return fmt.Errorf("email cannot start with '@'")
	}

	domain := email[atIndex+1:]
	if len(domain) == 0 {
		return fmt.Errorf("email must have a domain after '@'")
	}

	if !strings.Contains(domain, ".") {
		return fmt.Errorf("domain part must contain at least one '.'")
	}

	if strings.HasSuffix(domain, ".") {
		return fmt.Errorf("email domain cannot end with a '.'")
	}

	return nil
}
