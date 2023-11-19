package web

import (
	"context"
	"encoding/json"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
)

func (s *service) IncommingMsg(ctx context.Context, msgType app.MsgType, message []byte, clientID string) (app.MsgType, string, bool, error) {
	switch msgType {
	case app.MsgConnDisconnect:
		return app.MsgNil, "", true, nil
	case app.MsgConnSave:
		ok := s.writeValidator.IsReady(ctx)
		if !ok {
			return app.MsgServerFileNotReady, "", false, errors.New("validator not ready")
		}

		err := s.editor.Write(ctx)
		if err != nil {
			return app.MsgServerFileNotSaved, "", false, errors.Wrap(err, "editor write")
		}

		return app.MsgServerFileSaved, "", false, nil
	case app.MsgConnTextChange:
		var msg app.WSMsgTextChange

		err := json.Unmarshal(message, &msg)
		if err != nil {
			return app.MsgNil, "", true, errors.Wrap(err, "unmarshall text change msg")
		}

		err = s.editor.Change(ctx, &msg.Data)
		if err != nil {
			return app.MsgServerFileNotSaved, "", false, errors.Wrap(err, "editor write")
		}

		return app.MsgClientsTextChange, string(message), false, nil
	case app.MsgConnCursorChange:
		var msg app.WSMsgCursorChange

		err := json.Unmarshal(message, &msg)
		if err != nil {
			return app.MsgNil, "", true, errors.Wrap(err, "unmarshall cursor change msg")
		}

		s.hub.SetPosition(clientID, app.Position{
			Index:  msg.Data.Index,
			Length: msg.Data.Length,
		})

		return app.MsgClientsCursorChange, "", false, nil
	case app.MsgConnReady:
		err := s.writeValidator.ReadyClient(ctx, clientID)
		if err != nil {
			return app.MsgConnNotReady, "", false, errors.Wrap(err, "validator ready client")
		}

		s.hub.SetReady(clientID, true)

		// all clients are ready, save file
		ok := s.writeValidator.IsReady(ctx)
		if ok {
			err = s.editor.Write(ctx)
			if err != nil {
				return app.MsgServerFileNotSaved, "", false, errors.Wrap(err, "editor write")
			}

			s.writeValidator.Clear(ctx)
			s.hub.SetReadyAll(false)

			return app.MsgServerFileSaved, "", false, nil
		}

		return app.MsgClientsReady, "", false, nil
	case app.MsgConnUnready:
		err := s.writeValidator.UnreadyClient(ctx, clientID)
		if err != nil {
			return app.MsgConnNotUnready, "", false, errors.Wrap(err, "validator unready client")
		}

		s.hub.SetReady(clientID, false)

		return app.MsgClientsUnready, "", false, nil
	default:
		return app.MsgNil, "", true, errors.New("unknown message type")
	}
}
