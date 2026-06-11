package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/routes"
)

func main() {
	database.Connect()

	routes.RegisterRoutes()

	port := ":8080"
	fmt.Printf("Server is starting up and listening on port %s...\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Fatal: Server failed to start: %v", err)
	}
}
