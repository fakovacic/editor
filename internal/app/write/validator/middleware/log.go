package middleware

import (
	"context"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/log"
)

func NewLogMiddleware(next app.WriteValidator) app.WriteValidator {
	return &logMiddleware{
		next:    next,
		service: "write-validator",
	}
}

type logMiddleware struct {
	next    app.WriteValidator
	service string
}

func (m *logMiddleware) Clear(ctx context.Context) {
	m.next.Clear(ctx)
}

func (m *logMiddleware) AddClient(ctx context.Context, id string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "AddClient"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.AddClient(ctx, id)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "AddClient"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}

func (m *logMiddleware) RemoveClient(ctx context.Context, id string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "RemoveClient"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.RemoveClient(ctx, id)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "RemoveClient"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}

func (m *logMiddleware) ReadyClient(ctx context.Context, id string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "ReadyClient"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.ReadyClient(ctx, id)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "ReadyClient"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}

func (m *logMiddleware) UnreadyClient(ctx context.Context, id string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "UnreadyClient"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"id": id,
		}))

	err := m.next.UnreadyClient(ctx, id)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "UnreadyClient"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}

func (m *logMiddleware) IsReady(ctx context.Context) bool {
	ok := m.next.IsReady(ctx)

	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "IsReady"),
		log.String("layer", "part"),
		log.Any("res", map[string]any{
			"ok": ok,
		}))

	return ok
}
