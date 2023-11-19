package file

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
	"github.com/fakovacic/editor/internal/log"
)

func New(dirpath string, timeFunc func() time.Time) app.Versioning {
	return &versioningFile{
		Path:     dirpath,
		timeFunc: timeFunc,
	}
}

type versioningFile struct {
	Path     string
	timeFunc func() time.Time
}

func (s *versioningFile) Save(ctx context.Context, filename, content string) error {
	f, err := os.Create(fmt.Sprintf("%s/%d_%s", s.Path, s.timeFunc().Unix(), filename))
	if err != nil {
		return errors.Wrap(err, "create file")
	}

	_, err = f.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "write file")
	}

	defer f.Close()

	log.Info(ctx, fmt.Sprintf("version saved: %s", filename))

	return nil
}
