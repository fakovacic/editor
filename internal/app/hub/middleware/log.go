package middleware

import (
	"context"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/log"
)

func NewLogMiddleware(next app.Hub) app.Hub {
	return &logMiddleware{
		next:    next,
		service: "hub",
	}
}

type logMiddleware struct {
	next    app.Hub
	service string
}

func (m *logMiddleware) Get(id string) (*app.Client, bool) {
	return m.next.Get(id)
}

func (m *logMiddleware) GetByUsername(username string) bool {
	return m.next.GetByUsername(username)
}

func (m *logMiddleware) SetPosition(id string, pos app.Position) {
	m.next.SetPosition(id, pos)
}

func (m *logMiddleware) SetReady(id string, ready bool) {
	m.next.SetReady(id, ready)
}

func (m *logMiddleware) SetReadyAll(ready bool) {
	m.next.SetReadyAll(ready)
}

func (m *logMiddleware) CountRegistered() int {
	return m.next.CountRegistered()
}

func (m *logMiddleware) Create(client *app.Client) {
	m.next.Create(client)
}

func (m *logMiddleware) Register(client *app.Client) {
	m.next.Register(client)
}

func (m *logMiddleware) Unregister(client *app.Client) {
	m.next.Unregister(client)
}

func (m *logMiddleware) Brodcast(ctx context.Context, msgType app.MsgType, clientID, username, msg string, fileMeta *app.FileMeta) error {
	log.Info(ctx, "service-request",
		log.String("service", m.service),
		log.String("method", "Brodcast"),
		log.String("layer", "service"),
		log.Any("req", map[string]any{
			"msgType":  msgType,
			"clientID": clientID,
			"username": username,
		}))

	err := m.next.Brodcast(ctx, msgType, clientID, username, msg, fileMeta)

	log.Info(ctx, "service-response",
		log.String("service", m.service),
		log.String("method", "Brodcast"),
		log.String("layer", "service"),
		log.Err(err))

	return err
}
