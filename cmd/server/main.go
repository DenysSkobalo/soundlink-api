package main

import (
	"log"
	"os"

	"github.com/DenysSkobalo/soundlink-api/internal/app/handler"
	"github.com/DenysSkobalo/soundlink-api/internal/app/service"
	"github.com/DenysSkobalo/soundlink-api/internal/config"
	"github.com/DenysSkobalo/soundlink-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type User struct {
	Name string `json:"name"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fiberConfig := config.NewFiberConfig()
	app := fiber.New(*fiberConfig)

	app.Use(middleware.NewLogger(config.Config{
		Env:           os.Getenv("APP_ENV"),
		BodyLogMaxLen: 1024,
		SampleErrors:  true,
	}))

	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	app.Post("/signup", authHandler.Signup)

	log.Printf("Server is running on port %s\n", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
