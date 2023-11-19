package versioning

import (
	"strings"

	"github.com/fakovacic/editor/internal/errors"
)

type Type string

const (
	HTTP Type = "http"
	File Type = "file"
)

func (t Type) String() string {
	return string(t)
}

func (t *Type) Parse(s string) error {
	s = strings.Trim(s, "\"")
	switch s {
	case "http":
		*t = HTTP
	case "file":
		*t = File
	default:
		return errors.BadRequest("invalid msg type '%s'", s)
	}

	return nil
}
