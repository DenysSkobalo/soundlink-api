package config

import (
	"github.com/gofiber/fiber/v2"
)

type ServerConfig struct {
	FiberConfig *fiber.Config
}

func LoadServerConfig(env Environment) *ServerConfig {
	cfg := &fiber.Config{
		AppName:           "SoundLinkAPI",
		Prefork:           false,
		ServerHeader:      "SoundLinkAPI/0.1",
		EnablePrintRoutes: env == Dev,
		BodyLimit:         10 * 1024 * 1024,
	}

	return &ServerConfig{FiberConfig: cfg}
}
