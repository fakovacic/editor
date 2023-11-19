package middleware

import (
	"context"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/log"
)

func NewVersioningMiddleware(io editor.IO, versioning app.Versioning) editor.IO {
	return &versioningMiddleware{
		next:       io,
		versioning: versioning,
	}
}

type versioningMiddleware struct {
	next       editor.IO
	versioning app.Versioning
}

func (m *versioningMiddleware) Read(ctx context.Context) (string, *app.FileMeta, error) {
	content, file, err := m.next.Read(ctx)
	if err == nil {
		vErr := m.versioning.Save(ctx, file.Name, content)
		if vErr != nil {
			log.Error(ctx, "error while saving file: %v", vErr)
		}
	}

	return content, file, err
}

func (m *versioningMiddleware) Write(ctx context.Context, filename, content string) error {
	err := m.next.Write(ctx, filename, content)
	if err == nil {
		vErr := m.versioning.Save(ctx, filename, content)
		if vErr != nil {
			log.Error(ctx, "error while saving file: %v", vErr)
		}
	}

	return err
}
