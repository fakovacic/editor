package app

import "context"

//go:generate moq -out ./mocks/write_validator.go -pkg mocks  . WriteValidator
type WriteValidator interface {
	Clear(context.Context)

	IsReady(context.Context) bool

	AddClient(context.Context, string) error
	RemoveClient(context.Context, string) error

	ReadyClient(context.Context, string) error
	UnreadyClient(context.Context, string) error
}
