package config

import "github.com/gofiber/fiber/v2"

type Config struct {
	Env           string
	BodyLogMaxLen int
	SampleErrors  bool
}

func NewFiberConfig() *fiber.Config {
	return &fiber.Config{
		AppName:           "SoundLink",
		Prefork:           false,
		ServerHeader:      "SoundLinkAPI/0.1",
		EnablePrintRoutes: true,
		BodyLimit:         10 * 1024 * 1024,
	}
}
