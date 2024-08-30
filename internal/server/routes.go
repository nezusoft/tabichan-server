package server

import (
	"net/http"

	"github.com/tabichanorg/tabichan-server/internal/healthcheck"
	"github.com/tabichanorg/tabichan-server/internal/oauth"
)

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {

	mux.HandleFunc("/healthcheck", healthcheck.HealthCheck)

	oauth.RegisterOAuthHandlers(mux)
	// mux.HandleFunc("/signup", user.SignupHandler)
	// mux.HandleFunc("/login", user.LoginHandler)

	return mux
}
