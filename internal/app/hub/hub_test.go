package hub_test

import (
	"testing"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/hub"
	"github.com/fakovacic/editor/internal/app/hub/colors"
)

func TestHub(t *testing.T) {
	colorList := colors.New()

	hub := hub.New(colorList)

	hub.Create(&app.Client{
		ID:       "mock-id",
		Username: "mock-username",
	})

	client, ok := hub.Get("mock-id")
	if !ok {
		t.Fatal("client not found")
	}

	if client.Username != "mock-username" {
		t.Fatal("client username not valid")
	}

	ok = hub.GetByUsername("mock-username")
	if !ok {
		t.Fatal("client not found")
	}

	registered := hub.CountRegistered()
	if registered != 0 {
		t.Fatal("registered count not valid")
	}

	hub.Register(&app.Client{
		ID:       "mock-id",
		Username: "mock-username",
	})

	registered = hub.CountRegistered()
	if registered != 1 {
		t.Fatal("registered count not valid")
	}

	hub.SetReady("mock-id", true)

	hub.Unregister(&app.Client{
		ID:       "mock-id",
		Username: "mock-username",
	})

	registered = hub.CountRegistered()
	if registered != 0 {
		t.Fatal("registered count not valid")
	}

	ok = hub.GetByUsername("mock-username")
	if ok {
		t.Fatal("client deleted but still found")
	}

	_, ok = hub.Get("mock-id")
	if ok {
		t.Fatal("client deleted but still found")
	}
}
