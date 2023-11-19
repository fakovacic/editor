package editor

import (
	"context"
	"strings"
	"sync"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
	"github.com/fakovacic/editor/internal/log"
)

//go:generate moq -out ./mocks/io.go -pkg mocks  . IO
type IO interface {
	Read(context.Context) (string, *app.FileMeta, error)
	Write(context.Context, string, string) error
}

func New(io IO) app.Editor {
	return &editor{
		io: io,
	}
}

type editor struct {
	io   IO
	file editorFile
}

type editorFile struct {
	Meta     *app.FileMeta
	Contents string
	sync.Mutex
}

func (s *editor) FileMeta(_ context.Context) *app.FileMeta {
	s.file.Lock()
	defer s.file.Unlock()

	return s.file.Meta
}

func (s *editor) Load(ctx context.Context) error {
	if s.file.Contents != "" {
		return nil
	}

	s.file.Lock()
	defer s.file.Unlock()

	contents, meta, err := s.io.Read(ctx)
	if err != nil {
		return errors.Wrap(err, "io read")
	}

	s.file.Contents = strings.ReplaceAll(contents, "\t", "    ")
	s.file.Meta = meta

	return nil
}

func (s *editor) Unload(_ context.Context) error {
	s.file.Lock()
	defer s.file.Unlock()

	s.file.Contents = ""

	return nil
}

func (s *editor) Read(_ context.Context) (string, error) {
	s.file.Lock()
	defer s.file.Unlock()

	return s.file.Contents, nil
}

func (s *editor) Write(ctx context.Context) error {
	s.file.Lock()
	defer s.file.Unlock()

	if s.file.Contents == "" {
		log.Error(ctx, "contents empty")

		return nil
	}

	err := s.io.Write(ctx, s.file.Meta.Name, s.file.Contents)
	if err != nil {
		return errors.Wrap(err, "io write")
	}

	return nil
}
