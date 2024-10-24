package server

import (
	"net/http"

	"github.com/tabichanorg/tabichan-server/internal/db"
	"github.com/tabichanorg/tabichan-server/internal/healthcheck"
	middleware "github.com/tabichanorg/tabichan-server/internal/middleware/session"
	"github.com/tabichanorg/tabichan-server/internal/user"
)

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/healthcheck", healthcheck.HealthCheck)

	middlewareService := initMiddleware()

	userHandler := initUserHandler()
	mux.HandleFunc("/signup", userHandler.Signup)
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/user/details", middlewareService.SessionMiddleware(userHandler.GetUser))

	return mux
}

func initUserHandler() *user.UserHandler {
	userRepo := &user.UserRepository{Client: db.DynamoClient}
	userService := &user.UserService{Repo: userRepo}
	return &user.UserHandler{Service: userService}
}

func initMiddleware() *middleware.MiddlewareService {
	middlewareRepo := &middleware.MiddlewareRepository{Client: db.DynamoClient}
	return &middleware.MiddlewareService{Repo: middlewareRepo}
}
