package editor_test

import (
	"context"
	"testing"

	"github.com/fakovacic/editor/internal/app"
	"github.com/fakovacic/editor/internal/app/editor"
	"github.com/fakovacic/editor/internal/app/editor/mocks"
)

func TestEditor(t *testing.T) {
	ctx := context.Background()
	ioContent := "mock-content"

	io := &mocks.IOMock{
		ReadFunc: func(ctx context.Context) (string, *app.FileMeta, error) {
			return ioContent, &app.FileMeta{
				Name:      "mock-name",
				Extension: "mock-extension",
			}, nil
		},
		WriteFunc: func(ctx context.Context, contents string, name string) error {
			return nil
		},
	}

	editor := editor.New(io)

	// meta must be nil
	meta := editor.FileMeta(ctx)
	if meta != nil {
		t.Errorf("meta must be nil")
	}

	// content must be empty
	content, err := editor.Read(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	if content != "" {
		t.Errorf("content must be empty")
	}

	// load content
	err = editor.Load(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	// meta must be set
	meta = editor.FileMeta(ctx)
	if meta == nil {
		t.Errorf("meta must be set")
	}

	// content must be set
	content, err = editor.Read(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	if content == "" {
		t.Errorf("content must be set")
	}

	// write content
	err = editor.Write(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	// unload content
	err = editor.Unload(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	// contents must be empty
	content, err = editor.Read(ctx)
	if err != nil {
		t.Errorf("error must be nil")
	}

	if content != "" {
		t.Errorf("content must be empty")
	}
}
