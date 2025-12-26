package nylas

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with request ID",
			err:  APIError{StatusCode: 404, Message: "not found", RequestID: "req-123"},
			want: "nylas: not found (status=404, request_id=req-123)",
		},
		{
			name: "without request ID",
			err:  APIError{StatusCode: 500, Message: "server error"},
			want: "nylas: server error (status=500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("APIError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_Is(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		target     error
		want       bool
	}{
		{"400 is ErrBadRequest", 400, ErrBadRequest, true},
		{"401 is ErrUnauthorized", 401, ErrUnauthorized, true},
		{"404 is ErrNotFound", 404, ErrNotFound, true},
		{"429 is ErrRateLimited", 429, ErrRateLimited, true},
		{"500 is ErrServerError", 500, ErrServerError, true},
		{"502 is ErrServerError", 502, ErrServerError, true},
		{"404 is not ErrUnauthorized", 404, ErrUnauthorized, false},
		{"200 is not ErrNotFound", 200, ErrNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{StatusCode: tt.statusCode, Message: "test"}
			if got := errors.Is(err, tt.target); got != tt.want {
				t.Errorf("APIError.Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	if ErrMissingAPIKey.Error() != "nylas: API key required (use WithAPIKey)" {
		t.Error("ErrMissingAPIKey message incorrect")
	}
	if ErrUnauthorized.Error() != "nylas: unauthorized" {
		t.Error("ErrUnauthorized message incorrect")
	}
	if ErrNotFound.Error() != "nylas: not found" {
		t.Error("ErrNotFound message incorrect")
	}
	if ErrRateLimited.Error() != "nylas: rate limited" {
		t.Error("ErrRateLimited message incorrect")
	}
}
