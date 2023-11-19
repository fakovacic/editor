package app

import "context"

//go:generate moq -out ./mocks/versioning.go -pkg mocks  . Versioning
type Versioning interface {
	Save(context.Context, string, string) error
}
