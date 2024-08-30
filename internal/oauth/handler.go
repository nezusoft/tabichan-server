package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

var (
	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURI  = "http://localhost:8080/oauth/google/callback"

	appleClientID = os.Getenv("APPLE_CLIENT_ID")
	// appleClientSecret = os.Getenv("APPLE_CLIENT_SECRET")
	appleRedirectURI = "http://localhost:8080/oauth/apple/callback"
)

func RegisterOAuthHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/oauth/google/login", googleLoginHandler)
	mux.HandleFunc("/oauth/google/callback", googleCallbackHandler)
	mux.HandleFunc("/oauth/apple/login", appleLoginHandler)
	mux.HandleFunc("/oauth/apple/callback", appleCallbackHandler)
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" + url.Values{
		"client_id":     {googleClientID},
		"redirect_uri":  {googleRedirectURI},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {"state-token"},
	}.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := exchangeGoogleCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Google User Info: %v\n", userInfo)
}

func appleLoginHandler(w http.ResponseWriter, r *http.Request) {
	authURL := "https://appleid.apple.com/auth/authorize?" + url.Values{
		"client_id":     {appleClientID},
		"redirect_uri":  {appleRedirectURI},
		"response_type": {"code"},
		"scope":         {"name email"},
		"response_mode": {"form_post"},
		"state":         {"state-token"},
	}.Encode()

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func appleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// code := r.URL.Query().Get("code")
	// Placeholder for Apple's token exchange logic
	token := "fake-apple-token-for-demo"
	fmt.Fprintf(w, "Apple Token: %v\n", token)
}

func exchangeGoogleCodeForToken(code string) (string, error) {
	data := url.Values{
		"client_id":     {googleClientID},
		"client_secret": {googleClientSecret},
		"redirect_uri":  {googleRedirectURI},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res["access_token"].(string), nil
}

func fetchGoogleUserInfo(token string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
