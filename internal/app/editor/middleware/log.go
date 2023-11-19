package middleware

import (
	"context"

	"github.com/fakovacic/editor/internal/app"
)

func NewLogMiddleware(next app.Editor) app.Editor {
	return &logMiddleware{
		next:    next,
		service: "editor",
	}
}

type logMiddleware struct {
	next    app.Editor
	service string
}

func (m *logMiddleware) FileMeta(ctx context.Context) *app.FileMeta {
	return m.next.FileMeta(ctx)
}

func (m *logMiddleware) Load(ctx context.Context) error {
	return m.next.Load(ctx)
}

func (m *logMiddleware) Unload(ctx context.Context) error {
	return m.next.Unload(ctx)
}

func (m *logMiddleware) Read(ctx context.Context) (string, error) {
	return m.next.Read(ctx)
}

func (m *logMiddleware) Write(ctx context.Context) error {
	return m.next.Write(ctx)
}

func (m *logMiddleware) Change(ctx context.Context, msg *app.ChangeMsg) error {
	return m.next.Change(ctx, msg)
}
