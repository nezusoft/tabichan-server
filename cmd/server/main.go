package main

import (
	"log"
	"net/http"

	"github.com/tabichanorg/tabichan-server/internal/healthcheck"
)

func main() {
	http.HandleFunc("/healthcheck", healthcheck.Handler)

	log.Println("Starting TABICHAN on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server Failed to start: %v", err)
	}
}
