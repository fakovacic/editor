package app

import (
	"strings"

	"github.com/fakovacic/editor/internal/errors"
)

type FileMeta struct {
	Name      string   `json:"name"`
	Extension FileType `json:"extension"`
}

type FileType string

const (
	HTML       FileType = "html"
	CSS        FileType = "css"
	Javascript FileType = "javascript"
)

func (t FileType) String() string {
	return string(t)
}

func (t *FileType) Parse(s string) error {
	s = strings.Trim(s, "\"")
	switch s {
	case ".html":
		*t = HTML
	case ".css":
		*t = CSS
	case ".js":
		*t = Javascript
	default:
		return errors.BadRequest("invalid file type '%s'", s)
	}

	return nil
}
