package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/vishwakarmaritu/marketplace/internal/database"
	"github.com/vishwakarmaritu/marketplace/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}
	database.Connect()

	app := fiber.New()

	routes.RegisterRoutes(app)

	log.Println("Fiber server is starting up on port 8080...")

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Fatal: Server failed to start: %v", err)
	}
}
