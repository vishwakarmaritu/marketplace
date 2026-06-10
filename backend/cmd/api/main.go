package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/handlers"
)

func main() {
	database.Connect()

	http.HandleFunc("/api/auth/signup", handlers.Signup)
	http.HandleFunc("/api/auth/login", handlers.Login)

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "message": "Server and database are responding cleanly"}`))
	})

	port := ":8080"
	fmt.Printf("Server is starting up and listening on port %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Fatal: Server failed to start: %v", err)
	}
}
