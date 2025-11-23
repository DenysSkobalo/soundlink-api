package middleware

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoggerConfig struct {
	Pretty       bool
	BodyLogMax   int
	SampleErrors bool
}

func NewLogger(cfg *LoggerConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		var body string
		if cfg.SampleErrors && len(c.Body()) > 0 {
			body = string(c.Body())
			if len(body) > cfg.BodyLogMax {
				body = body[:cfg.BodyLogMax] + "...[truncated]"
			}
			body = redactSensitiveJSON(body)
		}

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		if cfg.Pretty {
			log.Printf("[%s] %s %s %d %s body: %s",
				time.Now().Format("02-Jan-2006 15:04:05"),
				c.Method(),
				c.Path(),
				status,
				latency,
				body,
			)
		} else {
			entry := map[string]interface{}{
				"time":    time.Now().Format(time.RFC3339),
				"method":  c.Method(),
				"path":    c.Path(),
				"status":  status,
				"latency": latency.String(),
			}
			if body != "" {
				entry["body"] = body
			}
			jsonEntry, _ := json.Marshal(entry)
			log.Println(string(jsonEntry))
		}

		return err
	}
}

func redactSensitiveJSON(body string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return body
	}

	sensitiveKeys := []string{"password", "token", "secret", "api_key"}
	for _, key := range sensitiveKeys {
		if _, ok := data[key]; ok {
			data[key] = "[REDACTED]"
		}
	}

	result, err := json.Marshal(data)
	if err != nil {
		return body
	}

	return string(result)
}
