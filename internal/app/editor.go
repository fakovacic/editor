package app

import (
	"context"
)

//go:generate moq -out ./mocks/editor.go -pkg mocks  . Editor
type Editor interface {
	FileMeta(context.Context) *FileMeta

	Load(context.Context) error
	Unload(context.Context) error

	Write(context.Context) error
	Read(context.Context) (string, error)

	Change(context.Context, *ChangeMsg) error
}
