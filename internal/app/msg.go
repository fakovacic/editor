package app

import (
	"strings"

	"github.com/fakovacic/editor/internal/errors"
)

type MsgType string

// MsgType from client to server
const (
	MsgConnSave         MsgType = "conn-save"          // conn save content to file
	MsgConnTextChange   MsgType = "conn-text-change"   // conn text change
	MsgConnCursorChange MsgType = "conn-cursor-change" // conn cursor change
	MsgConnReady        MsgType = "conn-ready"         // conn content ready for writing
	MsgConnUnready      MsgType = "conn-unready"       // conn content not ready for writing
	MsgConnDisconnect   MsgType = "conn-disconnect"    // disconnect conn
)

// MsgType from server to clients
const (
	MsgConnected MsgType = "conn-connected" // conn connected

	MsgClientsConnected    MsgType = "clients-connected"     // client connected
	MsgClientsReady        MsgType = "clients-ready"         // clients state change
	MsgClientsUnready      MsgType = "clients-unready"       // clients state change
	MsgClientsTextChange   MsgType = "clients-text-change"   // clients text change
	MsgClientsCursorChange MsgType = "clients-cursor-change" // clients cursor change
	MsgClientsDisconnected MsgType = "clients-disconnected"  // client disconnected

	MsgServerFileSaved MsgType = "server-file-saved" // file saved
)

// MsgType error from server to client
const (
	MsgConnNotReady   MsgType = "conn-not-ready"   // conn not ready
	MsgConnNotUnready MsgType = "conn-not-unready" // conn not unready

	MsgServerFileNotReady MsgType = "server-file-not-ready" // file not ready
	MsgServerFileNotSaved MsgType = "server-file-not-saved" // file not saved
)

const (
	MsgNil MsgType = "nil" // no msg
)

func (t MsgType) String() string {
	return string(t)
}

func (t *MsgType) UnmarshalJSON(b []byte) error {
	return t.Parse(string(b))
}

func (t *MsgType) Parse(s string) error {
	s = strings.Trim(s, "\"")
	switch s {
	case "conn-connected":
		*t = MsgConnected
	case "conn-save":
		*t = MsgConnSave
	case "conn-text-change":
		*t = MsgConnTextChange
	case "conn-cursor-change":
		*t = MsgConnCursorChange
	case "conn-ready":
		*t = MsgConnReady
	case "conn-unready":
		*t = MsgConnUnready
	case "conn-disconnect":
		*t = MsgConnDisconnect
	case "clients-connected":
		*t = MsgClientsConnected
	case "clients-text-change":
		*t = MsgClientsTextChange
	case "clients-cursor-change":
		*t = MsgClientsCursorChange
	case "clients-disconnected":
		*t = MsgClientsDisconnected
	case "server-file-saved":
		*t = MsgServerFileSaved
	case "conn-not-ready":
		*t = MsgConnNotReady
	case "conn-not-unready":
		*t = MsgConnNotUnready
	case "server-file-not-ready":
		*t = MsgServerFileNotReady
	case "server-file-not-saved":
		*t = MsgServerFileNotSaved
	case "nil":
		*t = MsgNil
	default:
		return errors.BadRequest("invalid msg type '%s'", s)
	}

	return nil
}

// WebSocket message
type WSMsg struct {
	Type MsgType `json:"type"`
	Data any     `json:"data"`
}

// WebSocket change message
type WSMsgTextChange struct {
	Data ChangeMsg `json:"data"`
}

type ChangeMsg struct {
	Action string    `json:"action"`
	Start  ChangeRow `json:"start"`
	End    ChangeRow `json:"end"`
	Lines  []string  `json:"lines"`
}

type ChangeRow struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

// WebSocket message for cursor change
type WSMsgCursorChange struct {
	Data CursorChange `json:"data"`
}

type CursorChange struct {
	Index  int `json:"index"`
	Length int `json:"length"`
}
