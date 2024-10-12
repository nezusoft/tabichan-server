package server

import (
	"net/http"

	"github.com/tabichanorg/tabichan-server/internal/db"
	"github.com/tabichanorg/tabichan-server/internal/healthcheck"
	"github.com/tabichanorg/tabichan-server/internal/user"
)

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/healthcheck", healthcheck.HealthCheck)

	userHandler := initUserHandler()
	mux.HandleFunc("/signup", userHandler.Signup)
	mux.HandleFunc("/login", userHandler.Login)

	return mux
}

func initUserHandler() *user.UserHandler {
	userRepo := &user.UserRepository{Client: db.DynamoClient}
	userService := &user.UserService{Repo: userRepo}
	return &user.UserHandler{Service: userService}
}
