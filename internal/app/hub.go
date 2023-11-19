package app

import (
	"context"
)

//go:generate moq -out ./mocks/ws_conn.go -pkg mocks  . WSConn
type WSConn interface {
	WriteMessage(messageType int, data []byte) error
}

//go:generate moq -out ./mocks/hub.go -pkg mocks  . Hub
type Hub interface {
	Get(string) (*Client, bool)
	GetByUsername(string) bool

	SetPosition(string, Position)
	SetReady(string, bool)
	SetReadyAll(bool)

	Create(*Client)
	Register(*Client)
	Unregister(*Client)

	CountRegistered() int
	Brodcast(context.Context, MsgType, string, string, string, *FileMeta) error
}

type Client struct {
	ID         string
	Username   string
	Color      string
	Ready      bool
	Registered bool
	Position   *Position
	Conn       WSConn
}

type Position struct {
	Index  int
	Length int
}
