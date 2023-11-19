package middleware

import (
	"context"
	"fmt"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/app/editor/io"
	"github.com/fakovacic/editor/internal/log"
)

func NewLogMiddleware(next editor.IO, ioType io.Type) editor.IO {
	return &logMiddleware{
		next:    next,
		service: fmt.Sprintf("io-%s", ioType),
	}
}

type logMiddleware struct {
	next    editor.IO
	service string
}

func (m *logMiddleware) Read(ctx context.Context) (string, *app.FileMeta, error) {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "Read"),
		log.String("layer", "part"))

	content, meta, err := m.next.Read(ctx)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "Read"),
		log.String("layer", "part"),
		log.Err(err))

	return content, meta, err
}

func (m *logMiddleware) Write(ctx context.Context, filename, content string) error {
	log.Info(ctx, "part-request",
		log.String("service", m.service),
		log.String("method", "Write"),
		log.String("layer", "part"),
		log.Any("req", map[string]any{
			"filename": filename,
		}))

	err := m.next.Write(ctx, filename, content)

	log.Info(ctx, "part-response",
		log.String("service", m.service),
		log.String("method", "Write"),
		log.String("layer", "part"),
		log.Err(err))

	return err
}
