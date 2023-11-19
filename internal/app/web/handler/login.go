package handler

import (
	"fmt"

	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) LoginForm() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("templates/login", nil)
	}
}

type LoginRequest struct {
	Username string `json:"username"`
}

func (h *Handler) Login() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var loginReq LoginRequest

		err := c.BodyParser(&loginReq)
		if err != nil {
			log.Error(c.Context(), fmt.Sprintf("body parse: %s", err))

			return c.Redirect("/login")
		}

		id, err := h.service.Login(c.Context(), loginReq.Username)
		if err != nil {
			log.Error(c.Context(), fmt.Sprintf("login: %s", err))

			return c.Redirect("/login")
		}

		return c.Redirect(fmt.Sprintf("/?id=%s", id))
	}
}
