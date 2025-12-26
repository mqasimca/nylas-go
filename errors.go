package nylas

import (
	"errors"
	"fmt"
)

var (
	ErrMissingAPIKey = errors.New("nylas: API key required (use WithAPIKey)")
	ErrUnauthorized  = errors.New("nylas: unauthorized")
	ErrNotFound      = errors.New("nylas: not found")
	ErrRateLimited   = errors.New("nylas: rate limited")
	ErrBadRequest    = errors.New("nylas: bad request")
	ErrServerError   = errors.New("nylas: server error")
)

type APIError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("nylas: %s (status=%d, request_id=%s)", e.Message, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("nylas: %s (status=%d)", e.Message, e.StatusCode)
}

func (e *APIError) Is(target error) bool {
	switch e.StatusCode {
	case 400:
		return target == ErrBadRequest
	case 401:
		return target == ErrUnauthorized
	case 404:
		return target == ErrNotFound
	case 429:
		return target == ErrRateLimited
	}
	if e.StatusCode >= 500 {
		return target == ErrServerError
	}
	return false
}
