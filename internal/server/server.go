package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr   string
	Server *http.Server
}

func NewServer(addr string) *Server {
	mux := mux.NewRouter()
	SetupRoutes(mux)

	return &Server{
		Addr: addr,
		Server: &http.Server{
			Addr:    addr,
			Handler: corsHandler(mux),
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Starting server on %s", s.Addr)
	return s.Server.ListenAndServe()
}

func (s *Server) Shutdown(timeout time.Duration) error {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
