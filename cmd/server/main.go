package main

import (
	"log"

	"github.com/DenysSkobalo/soundlink-api/internal/app/config"
	"github.com/DenysSkobalo/soundlink-api/internal/app/handler"
	"github.com/DenysSkobalo/soundlink-api/internal/app/service"
	"github.com/DenysSkobalo/soundlink-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	appConfig := config.LoadAppConfig()

	dbConfig := config.LoadDBConfig()
	if _, err := config.ConnectDB(dbConfig); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	app := fiber.New(*appConfig.Server.FiberConfig)

	app.Use(middleware.NewLogger(&middleware.LoggerConfig{
		Pretty:       appConfig.Logger.Pretty,
		BodyLogMax:   appConfig.Logger.BodyLogMax,
		SampleErrors: appConfig.Logger.SampleErrors,
	}))

	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	app.Post("/signup", authHandler.Signup)

	log.Printf("Starting %s on port %s in %s mode", appConfig.Server.FiberConfig.AppName, appConfig.AppPort, appConfig.Env)
	if err := app.Listen(":" + appConfig.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
