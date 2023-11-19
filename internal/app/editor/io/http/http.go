package http

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/errors"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(filepath string, client httpClient) editor.IO {
	return &ioHTTP{
		path:   filepath,
		client: client,
	}
}

type ioHTTP struct {
	client httpClient
	path   string
}

func (s *ioHTTP) Read(ctx context.Context) (string, *app.FileMeta, error) {
	var file app.FileMeta

	err := file.Extension.Parse(filepath.Ext(s.path))
	if err != nil {
		return "", nil, errors.Wrap(err, "parse extension")
	}

	file.Name = filepath.Base(s.path)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		s.path,
		http.NoBody,
	)
	if err != nil {
		return "", nil, errors.New("create http request: %v", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", nil, errors.New("http request:%v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.New("read body response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("status code not equal 200: %v", resp.StatusCode)
	}

	return string(body), &file, nil
}

func (s *ioHTTP) Write(ctx context.Context, _, content string) error {
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
