package web

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
	"github.com/fakovacic/editor/internal/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	msgsChan = 10
)

func (s *service) Connection(ctx context.Context, id string, c *websocket.Conn) error {
	client, err := s.register(ctx, id, c)
	if err != nil {
		return errors.Wrap(err, "register")
	}

	defer func() {
		deferErr := s.unregister(ctx, client)
		if deferErr != nil {
			log.Error(ctx, "unregister:", log.Err(deferErr))
		}
	}()

	// send contents to client
	fileContents, err := s.editor.Read(ctx)
	if err != nil {
		return errors.Wrap(err, "editor read content")
	}

	err = s.hub.Brodcast(ctx, app.MsgConnected, client.ID, client.Username, fileContents, s.editor.FileMeta(ctx))
	if err != nil {
		return errors.Wrap(err, "brodcast %s", app.MsgConnected)
	}

	// inform other clients about new client
	err = s.hub.Brodcast(ctx, app.MsgClientsConnected, client.ID, client.Username, "", nil)
	if err != nil {
		return errors.Wrap(err, "brodcast %s", app.MsgClientsConnected)
	}

	// conn ttl
	var timer *time.Timer

	disconectChan := make(chan bool, 1)

	if s.connTTL != nil {
		timer = time.NewTimer(*s.connTTL)
		go func() {
			<-timer.C

			disconectChan <- true
		}()
	}

	msgsChan := make(chan []byte, msgsChan)

	go s.msgReader(ctx, c, msgsChan, disconectChan)

	for {
		select {
		case <-disconectChan:
			return errors.New("connection timeout or disconnected")
		case message := <-msgsChan:
			if timer != nil {
				timer.Reset(*s.connTTL)
			}

			var wsMsg app.WSMsg

			err := json.Unmarshal(message, &wsMsg)
			if err != nil {
				return errors.Wrap(err, "unmarshal text msg")
			}

			log.Info(ctx, fmt.Sprintf("websocket message received: %s", wsMsg.Type))

			msgType, msgContent, closeConn, err := s.IncommingMsg(ctx, wsMsg.Type, message, client.ID)
			if err != nil {
				log.Error(ctx, "handle message:", log.Err(err))
			}

			if closeConn {
				return errors.New("close connection")
			}

			if msgType != app.MsgNil {
				err = s.hub.Brodcast(
					ctx,
					msgType,
					client.ID,
					client.Username,
					msgContent,
					nil,
				)
				if err != nil {
					log.Error(ctx, "brodcast:", log.Err(err))
				}
			}
		}
	}
}

func (s *service) register(ctx context.Context, id string, wsConn *websocket.Conn) (*app.Client, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("id not valid")
	}

	client, ok := s.hub.Get(id)
	if !ok {
		return nil, errors.New("id not exist")
	}

	client.Conn = wsConn

	// check client count
	// if no clients, load content
	if s.hub.CountRegistered() == 0 {
		loadErr := s.editor.Load(ctx)
		if loadErr != nil {
			return nil, errors.Wrap(loadErr, "editor load content")
		}
	}

	s.hub.Register(client)

	err = s.writeValidator.AddClient(ctx, client.ID)
	if err != nil {
		return nil, errors.Wrap(err, "validator add client")
	}

	return client, nil
}

func (s *service) unregister(ctx context.Context, client *app.Client) error {
	err := s.hub.Brodcast(ctx, app.MsgConnDisconnect, client.ID, client.Username, "", nil)
	if err != nil {
		return errors.Wrap(err, "brodcast %s", app.MsgConnDisconnect)
	}

	s.hub.Unregister(client)

	err = s.hub.Brodcast(ctx, app.MsgClientsDisconnected, client.ID, client.Username, "", nil)
	if err != nil {
		return errors.Wrap(err, "brodcast %s", app.MsgClientsDisconnected)
	}

	err = s.writeValidator.RemoveClient(ctx, client.ID)
	if err != nil {
		return errors.Wrap(err, "validator remove client")
	}

	// check client count
	// if no clients, save file
	if s.hub.CountRegistered() == 0 {
		writeErr := s.editor.Write(ctx)
		if writeErr != nil {
			return errors.Wrap(writeErr, "editor write")
		}

		unloadErr := s.editor.Unload(ctx)
		if unloadErr != nil {
			return errors.Wrap(unloadErr, "editor unload")
		}
	}

	return nil
}

func (s *service) msgReader(ctx context.Context, wsConn *websocket.Conn, msgsChan chan []byte, disconectChan chan bool) {
	defer func() {
		close(msgsChan)
	}()

	for {
		messageType, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(ctx, "unexpected close:", log.Err(err))

				disconectChan <- true

				return
			}

			log.Error(ctx, "read error:", log.Err(err))

			disconectChan <- true

			return
		}

		switch messageType {
		case websocket.TextMessage:
			msgsChan <- message
		default:
			log.Info(ctx, "websocket message received of type", messageType)
		}
	}
}
