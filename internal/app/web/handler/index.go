package handler

import (
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Index() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			log.Error(c.Context(), "id empty")

			return c.Redirect("/login")
		}

		err := h.service.Index(c.Context(), id)
		if err != nil {
			log.Error(c.Context(), err.Error())

			return c.Redirect("/login")
		}

		return c.Render("templates/index", nil)
	}
}
