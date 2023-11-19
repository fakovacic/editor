package hub

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
	"github.com/gofiber/contrib/websocket"
)

type BrodcastMsg struct {
	Data     string        `json:"data,omitempty"`
	FileMeta *app.FileMeta `json:"fileMeta,omitempty"`
	Type     app.MsgType   `json:"type"`
	Client   string        `json:"client"`
	Clients  []Client      `json:"clients"`
}

func (h *hub) Brodcast(_ context.Context, msgType app.MsgType, clientID, username, msg string, fileMeta *app.FileMeta) error {
	if len(h.clients) == 0 {
		return nil
	}

	h.Lock()
	defer h.Unlock()

	for _, client := range h.clients {
		switch msgType {
		case app.MsgConnected, app.MsgConnDisconnect, app.MsgConnNotReady, app.MsgConnNotUnready, app.MsgServerFileNotReady:
			// only send to the client who sent the message
			if client.ID != clientID {
				continue
			}

		case app.MsgClientsConnected, app.MsgClientsTextChange, app.MsgClientsDisconnected, app.MsgClientsCursorChange:
			// send to all clients except the one that sent the message
			if client.ID == clientID {
				continue
			}

		case app.MsgServerFileNotSaved, app.MsgServerFileSaved, app.MsgClientsReady, app.MsgClientsUnready:
			// send to all clients
		default:
			continue
		}

		msg, err := json.Marshal(BrodcastMsg{
			Data:     msg,
			FileMeta: fileMeta,
			Client:   username,
			Type:     msgType,
			Clients:  h.viewClients(),
		})
		if err != nil {
			return errors.Wrap(err, "marshaling message")
		}

		err = client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return errors.Wrap(err, "writing message to client '%s'", client.ID)
		}
	}

	return nil
}

func (h *hub) viewClients() []Client {
	viewClients := make([]Client, 0)

	for k := range h.clients {
		viewClients = append(viewClients, ToClient(h.clients[k]))
	}

	sort.Slice(viewClients, func(i, j int) bool {
		return viewClients[i].Color < viewClients[j].Color
	})

	return viewClients
}
