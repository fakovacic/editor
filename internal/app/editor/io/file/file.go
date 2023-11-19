package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/errors"
)

func New(filepath string) editor.IO {
	return &ioFile{
		Path: filepath,
	}
}

type ioFile struct {
	Path string
}

func (s *ioFile) Read(_ context.Context) (string, *app.FileMeta, error) {
	var file app.FileMeta

	err := file.Extension.Parse(filepath.Ext(s.Path))
	if err != nil {
		return "", nil, errors.Wrap(err, "extension parsing")
	}

	file.Name = filepath.Base(s.Path)

	contents, err := os.ReadFile(s.Path)
	if err != nil {
		return "", nil, errors.Wrap(err, "os read file")
	}

	return string(contents), &file, nil
}

func (s *ioFile) Write(_ context.Context, _, content string) error {
	err := os.WriteFile(s.Path, []byte(content), 0644)
	if err != nil {
		return errors.Wrap(err, "os flush file")
	}

	return nil
}
