package middleware

import (
	"context"
	"fmt"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/versioning"
	"github.com/fakovacic/editor/internal/log"
)

func NewLogMiddleware(next app.Versioning, versioningType versioning.Type) app.Versioning {
	return &logMiddleware{
		next:    next,
		service: fmt.Sprintf("versioning-%s", versioningType),
	}
}

type logMiddleware struct {
	next    app.Versioning
	service string
}

func (m *logMiddleware) Save(ctx context.Context, filename, content string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "Save"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"filename": filename,
		}))

	err := m.next.Save(ctx, filename, content)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "Save"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}
