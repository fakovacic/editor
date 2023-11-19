package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/errors"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(dirpath string, client httpClient, timeFunc func() time.Time) app.Versioning {
	return &versioningFile{
		path:     dirpath,
		client:   client,
		timeFunc: timeFunc,
	}
}

type versioningFile struct {
	path     string
	client   httpClient
	timeFunc func() time.Time
}

func (s *versioningFile) Save(ctx context.Context, _, content string) error {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		s.path,
		strings.NewReader(content),
	)
	if err != nil {
		return errors.New("create http request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return errors.New("http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code not equal 200: %v", resp.StatusCode)
	}

	return nil
}
