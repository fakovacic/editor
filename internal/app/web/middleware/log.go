package middleware

import (
	"context"

	"github.com/fakovacic/editor/internal/app/web"
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/contrib/websocket"
)

func NewLogMiddleware(next web.Service) web.Service {
	return &logMiddleware{
		next:    next,
		service: "web",
	}
}

type logMiddleware struct {
	next    web.Service
	service string
}

func (m *logMiddleware) Login(ctx context.Context, username string) (string, error) {
	log.Info(ctx, "service-request",
		log.String("service", m.service),
		log.String("method", "Login"),
		log.String("layer", "service"),
		log.Any("req", map[string]any{
			"username": username,
		}))

	id, err := m.next.Login(ctx, username)

	log.Info(ctx, "service-response",
		log.String("service", m.service),
		log.String("method", "Login"),
		log.String("layer", "service"),
		log.Any("res", map[string]any{
			"id": id,
		}),
		log.Err(err))

	return id, err
}

func (m *logMiddleware) Index(ctx context.Context, id string) error {
	log.Info(ctx, "service-request",
		log.String("service", m.service),
		log.String("method", "Index"),
		log.String("layer", "service"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.Index(ctx, id)

	log.Info(ctx, "service-response",
		log.String("service", m.service),
		log.String("method", "Index"),
		log.String("layer", "service"),
		log.Err(err))

	return err
}

func (m *logMiddleware) Connection(ctx context.Context, id string, c *websocket.Conn) error {
	log.Info(ctx, "service-request",
		log.String("service", m.service),
		log.String("method", "Connection"),
		log.String("layer", "service"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.Connection(ctx, id, c)

	log.Info(ctx, "service-response",
		log.String("service", m.service),
		log.String("method", "Connection"),
		log.String("layer", "service"),
		log.Err(err))

	return err
}
