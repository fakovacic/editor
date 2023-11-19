package hub

import (
	"sync"

	"github.com/fakovacic/editor/internal/app"
)

type Colors interface {
	Lock() string
	Release(string)
}

func New(colors Colors) app.Hub {
	return &hub{
		colors:  colors,
		clients: make(map[string]*app.Client),
	}
}

type hub struct {
	clients map[string]*app.Client
	colors  Colors
	sync.Mutex
}

func (h *hub) Get(id string) (*app.Client, bool) {
	h.Lock()
	defer h.Unlock()

	client, ok := h.clients[id]
	if !ok {
		return nil, false
	}

	return client, true
}

func (h *hub) GetByUsername(username string) bool {
	h.Lock()
	defer h.Unlock()

	var exist bool

	for _, client := range h.clients {
		if client.Username == username {
			exist = true

			break
		}
	}

	return exist
}

func (h *hub) SetPosition(id string, pos app.Position) {
	h.Lock()
	defer h.Unlock()

	client, ok := h.clients[id]
	if !ok {
		return
	}

	client.Position = &pos

	h.clients[id] = client
}

func (h *hub) SetReady(id string, ready bool) {
	h.Lock()
	defer h.Unlock()

	client, ok := h.clients[id]
	if !ok {
		return
	}

	client.Ready = ready

	h.clients[id] = client
}

func (h *hub) SetReadyAll(ready bool) {
	h.Lock()
	defer h.Unlock()

	for id, cl := range h.clients {
		cl.Ready = ready

		h.clients[id] = cl
	}
}

func (h *hub) Create(client *app.Client) {
	h.Lock()
	defer h.Unlock()

	h.clients[client.ID] = client
}

func (h *hub) Register(client *app.Client) {
	h.Lock()
	defer h.Unlock()

	client.Color = h.colors.Lock()
	client.Registered = true

	h.clients[client.ID] = client
}

func (h *hub) Unregister(client *app.Client) {
	h.Lock()
	defer h.Unlock()

	h.colors.Release(client.Color)

	delete(h.clients, client.ID)
}

func (h *hub) CountRegistered() int {
	i := 0

	for id := range h.clients {
		if h.clients[id].Registered {
			i++
		}
	}

	return i
}
