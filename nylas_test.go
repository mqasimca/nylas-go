package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr error
	}{
		{
			name:    "missing API key",
			opts:    nil,
			wantErr: ErrMissingAPIKey,
		},
		{
			name:    "with API key",
			opts:    []Option{WithAPIKey("test-key")},
			wantErr: nil,
		},
		{
			name: "with all options",
			opts: []Option{
				WithAPIKey("test-key"),
				WithRegion(RegionEU),
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.opts...)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("NewClient() unexpected error = %v", err)
				return
			}
			if client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_NewRequest(t *testing.T) {
	client, _ := NewClient(WithAPIKey("test-key"))

	tests := []struct {
		name    string
		method  string
		path    string
		body    any
		wantErr bool
	}{
		{
			name:   "GET request",
			method: http.MethodGet,
			path:   "/v3/grants/123/messages",
			body:   nil,
		},
		{
			name:   "POST request with body",
			method: http.MethodPost,
			path:   "/v3/grants/123/messages/send",
			body:   map[string]string{"subject": "Test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := client.NewRequest(context.Background(), tt.method, tt.path, tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if req == nil {
				t.Error("NewRequest() returned nil request")
				return
			}
			if req.Method != tt.method {
				t.Errorf("NewRequest() method = %v, want %v", req.Method, tt.method)
			}
			if auth := req.Header.Get("Authorization"); auth != "Bearer test-key" {
				t.Errorf("NewRequest() Authorization = %v, want Bearer test-key", auth)
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: 200,
			response:   `{"data": {"id": "123"}, "request_id": "req-1"}`,
			wantErr:    false,
		},
		{
			name:       "not found",
			statusCode: 404,
			response:   `{"message": "not found", "type": "error"}`,
			wantErr:    true,
		},
		{
			name:       "unauthorized",
			statusCode: 401,
			response:   `{"message": "unauthorized", "type": "error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Request-Id", "req-123")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(
				WithAPIKey("test-key"),
				WithBaseURL(srv.URL),
				WithMaxRetries(0),
			)

			req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
			var result map[string]string
			_, err := client.Do(req, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Do_Retry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]string{"id": "123"},
			"request_id": "req-1",
		})
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
		WithMaxRetries(3),
		WithRetryWait(1),
	)

	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	var result map[string]string
	_, err := client.Do(req, &result)

	if err != nil {
		t.Errorf("Do() with retry error = %v", err)
	}
	if attempts != 3 {
		t.Errorf("Do() attempts = %d, want 3", attempts)
	}
}

func TestClient_Do_RateLimitRetry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(429)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"message": "rate limit exceeded",
				"type":    "error",
			})
			return
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]string{"id": "123"},
			"request_id": "req-1",
		})
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
		WithMaxRetries(3),
		WithRetryWait(1),
	)

	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	var result map[string]string
	_, err := client.Do(req, &result)

	if err != nil {
		t.Errorf("Do() with rate limit retry error = %v", err)
	}
	if attempts != 3 {
		t.Errorf("Do() attempts = %d, want 3", attempts)
	}
}

func TestClient_Do_RateLimitExhausted(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(429)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "rate limit exceeded",
			"type":    "error",
		})
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
		WithMaxRetries(2),
		WithRetryWait(1),
	)

	req, _ := client.NewRequest(context.Background(), http.MethodGet, "/test", nil)
	var result map[string]string
	_, err := client.Do(req, &result)

	// Should return error after exhausting retries
	if err == nil {
		t.Error("Do() expected error after exhausting retries")
	}
	if attempts != 3 { // Initial + 2 retries
		t.Errorf("Do() attempts = %d, want 3", attempts)
	}
}

func TestClient_Services(t *testing.T) {
	client, _ := NewClient(WithAPIKey("test-key"))

	if client.Messages == nil {
		t.Error("Messages service is nil")
	}
	if client.Threads == nil {
		t.Error("Threads service is nil")
	}
	if client.Calendars == nil {
		t.Error("Calendars service is nil")
	}
	if client.Events == nil {
		t.Error("Events service is nil")
	}
}
