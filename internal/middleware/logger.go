package middleware

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/DenysSkobalo/soundlink-api/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func NewLogger(cfg config.Config) fiber.Handler {
	var logger zerolog.Logger
	if cfg.Env == "dev" {
		out := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		logger = zerolog.New(out).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		rid := c.Get("X-Request-Id")
		if rid == "" {
			rid = uuid.New().String()
			c.Set("X-Request-Id", rid)
		}

		var reqBodySnippet string
		if cfg.Env != "prod" && cfg.BodyLogMaxLen > 0 {
			body := c.Body()
			if len(body) > cfg.BodyLogMaxLen {
				reqBodySnippet = string(body[:cfg.BodyLogMaxLen]) + "...(truncated)"
			} else {
				reqBodySnippet = string(body)
			}
			reqBodySnippet = redactSensitive(reqBodySnippet)
		}

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		clientIP := c.IP()
		method := c.Method()
		path := c.Path()
		ua := c.Get("User-Agent")

		event := logger.Info().
			Str("req_id", rid).
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Str("client_ip", clientIP).
			Str("user_agent", ua).
			Dur("latency", latency)

		if reqBodySnippet != "" {
			event.Str("req_body", reqBodySnippet)
		}

		if err != nil {
			event = logger.Error().
				Err(err).
				Str("req_id", rid).
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Str("client_ip", clientIP).
				Dur("latency", latency)
			if cfg.SampleErrors && cfg.BodyLogMaxLen > 0 {
				body := c.Body()
				if len(body) > cfg.BodyLogMaxLen {
					event.Str("req_body", string(body[:cfg.BodyLogMaxLen])+"...(truncated)")
				} else {
					event.Str("req_body", string(body))
				}
			}
			event.Msg("request error")
			return err
		}

		event.Msg("request completed")
		return nil
	}
}

func redactSensitive(body string) string {
	if body == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		log.Printf("⚠️ failed to parse body for redaction: %v", err)
		return body // якщо помилка — повертаємо оригінал
	}

	sensitiveKeys := []string{"password", "token", "secret", "api_key"}
	for _, key := range sensitiveKeys {
		if _, exists := data[key]; exists {
			data[key] = "[REDACTED]"
		}
	}

	redacted, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal redacted JSON: %v", err)
		return body
	}

	return string(redacted)
}
