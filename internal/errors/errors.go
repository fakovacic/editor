package errors

import (
	goErrors "errors"
	"net/http"

	"github.com/pkg/errors"
)

const (
	BadRequestCode       = 100
	NotFoundCode         = 101
	UnauthorizedCode     = 102
	MethodNotAllowedCode = 103
	InternalCode         = 104
)

type Error struct {
	Code int
	Err  error
}

func (t Error) Error() string {
	return t.Err.Error()
}

func (t Error) GetCode() int {
	return t.Code
}

func (t Error) HTTPStatusCode() int {
	switch t.Code {
	case BadRequestCode:
		return http.StatusBadRequest
	case NotFoundCode:
		return http.StatusNotFound
	case UnauthorizedCode:
		return http.StatusUnauthorized
	case MethodNotAllowedCode:
		return http.StatusMethodNotAllowed
	case InternalCode:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func New(format string, args ...interface{}) error {
	return Error{
		Code: InternalCode,
		Err:  errors.Errorf(format, args...),
	}
}

func Wrap(err error, format string, args ...interface{}) error {
	status := InternalCode

	var e Error

	if goErrors.As(err, &e) {
		status = e.GetCode()
	}

	return Error{
		Code: status,
		Err:  errors.Wrapf(err, format, args...),
	}
}

func Internal(format string, args ...interface{}) error {
	return Error{
		Code: InternalCode,
		Err:  errors.Errorf(format, args...),
	}
}

func InternalWrap(err error, format string, args ...interface{}) error {
	return Error{
		Code: InternalCode,
		Err:  errors.Wrapf(err, format, args...),
	}
}

func BadRequest(format string, args ...interface{}) error {
	return Error{
		Code: BadRequestCode,
		Err:  errors.Errorf(format, args...),
	}
}

func BadRequestWrap(err error, format string, args ...interface{}) error {
	return Error{
		Code: BadRequestCode,
		Err:  errors.Wrapf(err, format, args...),
	}
}

func NotFound(format string, args ...interface{}) error {
	return Error{
		Code: NotFoundCode,
		Err:  errors.Errorf(format, args...),
	}
}

func NotFoundWrap(err error, format string, args ...interface{}) error {
	return Error{
		Code: NotFoundCode,
		Err:  errors.Wrapf(err, format, args...),
	}
}

func Unauthorized(format string, args ...interface{}) error {
	return Error{
		Code: UnauthorizedCode,
		Err:  errors.Errorf(format, args...),
	}
}

func UnauthorizedWrap(err error, format string, args ...interface{}) error {
	return Error{
		Code: UnauthorizedCode,
		Err:  errors.Wrapf(err, format, args...),
	}
}

func MethodNotAllowed(format string, args ...interface{}) error {
	return Error{
		Code: MethodNotAllowedCode,
		Err:  errors.Errorf(format, args...),
	}
}

func MethodNotAllowedWrap(err error, format string, args ...interface{}) error {
	return Error{
		Code: MethodNotAllowedCode,
		Err:  errors.Wrapf(err, format, args...),
	}
}
