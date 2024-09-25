package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Addr   string
	Server *http.Server
}

func NewServer(addr string) *Server {
	mux := http.NewServeMux()
	SetupRoutes(mux)

	return &Server{
		Addr: addr,
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
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
