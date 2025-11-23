package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Dev   Environment = "dev"
	Stage Environment = "stage"
	Prod  Environment = "prod"
)

type AppConfig struct {
	Env     Environment
	AppPort string
	Server  *ServerConfig
	DB      *DBConfig
	Logger  *LoggerConfig
}

func LoadAppConfig() *AppConfig {
	env := Environment(os.Getenv("APP_ENV"))
	if env == "" {
		env = Dev
	}

	if env == Dev {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on system environment")
		}
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &AppConfig{
		Env:     env,
		AppPort: appPort,
		Server:  LoadServerConfig(env),
		DB:      LoadDBConfig(),
		Logger:  LoadLoggerConfig(env),
	}
}
