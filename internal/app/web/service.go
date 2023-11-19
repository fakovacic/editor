package web

import (
	"context"
	"time"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type Service interface {
	Login(context.Context, string) (string, error)
	Index(context.Context, string) error
	Connection(context.Context, string, *websocket.Conn) error
}

func New(editor app.Editor, hub app.Hub, writeValidator app.WriteValidator, connTTL *time.Duration) Service {
	return &service{
		editor:         editor,
		hub:            hub,
		writeValidator: writeValidator,
		connTTL:        connTTL,
	}
}

type service struct {
	editor         app.Editor
	hub            app.Hub
	writeValidator app.WriteValidator
	connTTL        *time.Duration
}

func (s *service) Login(_ context.Context, username string) (string, error) {
	ok := s.hub.GetByUsername(username)
	if ok {
		return "", errors.New("username already exist")
	}

	client := &app.Client{
		ID:       uuid.New().String(),
		Username: username,
	}

	s.hub.Create(client)

	return client.ID, nil
}

func (s *service) Index(_ context.Context, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id not valid")
	}

	_, ok := s.hub.Get(id)
	if !ok {
		return errors.New("id not exist")
	}

	return nil
}
