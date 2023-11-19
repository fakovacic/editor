package hub_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/hub"
	"github.com/fakovacic/editor/internal/app/hub/colors"
	"github.com/fakovacic/editor/internal/app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBrodcast(t *testing.T) {
	cases := []struct {
		it string

		clients []*app.Client

		msgType  app.MsgType
		clientID string
		username string
		msg      string
		fileMeta *app.FileMeta

		expectedMsgs  map[string]hub.BrodcastMsg
		expectedError string
	}{
		{
			it: "send message connected",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgConnected,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgConnected,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message conn-not-ready",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgConnNotReady,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgConnNotReady,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message conn-not-unready",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgConnNotUnready,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgConnNotUnready,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message server-file-not-ready",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgServerFileNotReady,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgServerFileNotReady,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message clients-connected",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgClientsConnected,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id-next": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgClientsConnected,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message clients-text-change",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgClientsTextChange,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id-next": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgClientsTextChange,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message clients-disconnected",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgClientsDisconnected,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id-next": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgClientsDisconnected,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message server-file-not-saved",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgServerFileNotSaved,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,

			expectedMsgs: map[string]hub.BrodcastMsg{
				"mock-id": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgServerFileNotSaved,
					Client:   "mock-username",
				},
				"mock-id-next": {
					Data:     "mock-message",
					FileMeta: nil,
					Type:     app.MsgServerFileNotSaved,
					Client:   "mock-username",
				},
			},
		},
		{
			it: "send message conn-disconnect",

			clients: []*app.Client{
				{
					ID:       "mock-id",
					Username: "mock-username",
					Color:    "mock-color",
				},
				{
					ID:       "mock-id-next",
					Username: "mock-username-next",
					Color:    "mock-color-next",
				},
			},

			msgType:  app.MsgConnDisconnect,
			clientID: "mock-id",
			username: "mock-username",
			msg:      "mock-message",
			fileMeta: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.it, func(t *testing.T) {
			colorList := colors.New()

			hb := hub.New(colorList)

			for cl := range tc.clients {
				conn := &mocks.WSConnMock{
					WriteMessageFunc: func(messageType int, data []byte) error {
						var msg hub.BrodcastMsg

						err := json.Unmarshal(data, &msg)
						if err != nil {
							return err
						}

						msg, ok := tc.expectedMsgs[tc.clients[cl].ID]
						if !ok {
							return nil
						}

						assert.Equal(t, msg.Data, msg.Data)
						assert.Equal(t, msg.FileMeta, msg.FileMeta)
						assert.Equal(t, msg.Type, msg.Type)
						assert.Equal(t, msg.Client, msg.Client)

						return nil
					},
				}

				tc.clients[cl].Conn = conn

				hb.Register(tc.clients[cl])
			}

			err := hb.Brodcast(context.Background(), tc.msgType, tc.clientID, tc.username, tc.msg, tc.fileMeta)
			if err != nil {
				assert.Equal(t, err.Error(), tc.expectedError)
			}
		})
	}
}
