package validator

import (
	"context"
	"sync"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
)

func New() app.WriteValidator {
	return &writeValidator{
		clients: make(map[string]bool),
	}
}

type writeValidator struct {
	clients map[string]bool
	sync.Mutex
}

func (s *writeValidator) Clear(context.Context) {
	s.Lock()
	defer s.Unlock()

	for id := range s.clients {
		s.clients[id] = false
	}
}

func (s *writeValidator) AddClient(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.clients[id]
	if ok {
		return nil
	}

	s.clients[id] = false

	return nil
}

func (s *writeValidator) RemoveClient(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.clients[id]
	if ok {
		delete(s.clients, id)

		return nil
	}

	return errors.New("client not found")
}

func (s *writeValidator) ReadyClient(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.clients[id]
	if ok {
		s.clients[id] = true

		return nil
	}

	return errors.New("client not found")
}

func (s *writeValidator) UnreadyClient(_ context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.clients[id]
	if ok {
		s.clients[id] = false

		return nil
	}

	return errors.New("client not found")
}

func (s *writeValidator) IsReady(_ context.Context) bool {
	s.Lock()
	defer s.Unlock()

	if len(s.clients) == 0 {
		return false
	}

	if len(s.clients) == 1 {
		return true
	}

	for _, ready := range s.clients {
		if !ready {
			return false
		}
	}

	return true
}
