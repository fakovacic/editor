package middleware

import (
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Info(c.Context(), "http-request",
			log.String("url", c.OriginalURL()),
			log.String("method", c.Route().Method),
		)

		return c.Next()
	}
}
