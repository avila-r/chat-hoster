package main

import (
	"log"

	"github.com/avila-r/chat-hoster/internal/router"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	app := fiber.New()

	router.EnableRouting(app)
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
