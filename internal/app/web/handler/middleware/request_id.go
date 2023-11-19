package middleware

import (
	appCtx "github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ReqID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := uuid.New().String()

		ctx := c.Context()
		ctx.SetUserValue(appCtx.RequestID.String(), reqID)

		log.AddFields(c.Context(),
			log.String(appCtx.RequestID.String(), reqID),
		)

		return c.Next()
	}
}
