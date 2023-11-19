package handler

import (
	"context"

	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) WS() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		ctx := context.Background()

		id := c.Query("id")
		if id == "" {
			log.Error(ctx, "id empty")

			_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			closeErr := c.Close()
			if closeErr != nil {
				log.Error(ctx, "closing connection:", log.Err(closeErr))

				return
			}

			return
		}

		err := h.service.Connection(ctx, id, c)
		if err != nil {
			log.Error(ctx, "ws connection", log.Err(err))

			_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			closeErr := c.Close()
			if closeErr != nil {
				log.Error(ctx, "closing connection:", log.Err(closeErr))

				return
			}

			return
		}
	})
}
