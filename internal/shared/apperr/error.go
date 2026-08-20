package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

type Kind string

const (
	KindBadRequest      Kind = "BAD_REQUEST"
	KindUnauthorized    Kind = "UNAUTHORIZED"
	KindForbidden       Kind = "FORBIDDEN"
	KindNotFound        Kind = "NOT_FOUND"
	KindConflict        Kind = "CONFLICT"
	KindValidation      Kind = "VALIDATION"
	KindRateLimited     Kind = "RATE_LIMITED"
	KindPayloadTooLarge Kind = "PAYLOAD_TOO_LARGE"
	KindInternal        Kind = "INTERNAL"
	KindUnavailable     Kind = "UNAVAILABLE"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Kind       Kind         `json:"-"`
	Code       string       `json:"code"`
	Message    string       `json:"message"`
	Details    []FieldError `json:"field_errors,omitempty"`
	HTTPStatus int          `json:"-"`
	cause      error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message, HTTPStatus: StatusFor(kind)}
}

func Wrap(err error, kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message, HTTPStatus: StatusFor(kind), cause: err}
}

func BadRequest(code, message string) *Error {
	return New(KindBadRequest, code, message)
}

func Validation(code string, details []FieldError) *Error {
	return &Error{Kind: KindValidation, Code: code, Message: "请求参数校验失败", HTTPStatus: http.StatusUnprocessableEntity, Details: details}
}

func NotFound(code, message string) *Error {
	return New(KindNotFound, code, message)
}

func Conflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

func Internal(code, message string) *Error {
	return New(KindInternal, code, message)
}

func PayloadTooLarge(err error) *Error {
	return Wrap(err, KindPayloadTooLarge, "REQUEST_BODY_TOO_LARGE", "请求体超过限制")
}

func IsKind(err error, kind Kind) bool {
	var target *Error
	return errors.As(err, &target) && target.Kind == kind
}

func StatusFor(kind Kind) int {
	switch kind {
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindValidation:
		return http.StatusUnprocessableEntity
	case KindRateLimited:
		return http.StatusTooManyRequests
	case KindUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
