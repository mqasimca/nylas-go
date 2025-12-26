package nylas

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAttachmentsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response: `{
				"request_id": "req-123",
				"data": {
					"id": "attach-1",
					"grant_id": "grant-123",
					"filename": "document.pdf",
					"content_type": "application/pdf",
					"size": 12345,
					"is_inline": false
				}
			}`,
			wantErr: false,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			response:   `{"error": "not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/v3/grants/grant-123/attachments/attach-1" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				if r.URL.Query().Get("message_id") != "msg-123" {
					t.Errorf("expected message_id=msg-123, got %s", r.URL.Query().Get("message_id"))
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			attachment, err := client.Attachments.Get(context.Background(), "grant-123", "attach-1", "msg-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && attachment.ID != "attach-1" {
				t.Errorf("expected attach-1, got %s", attachment.ID)
			}
		})
	}
}

func TestAttachmentsService_Download(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		body        string
		contentType string
		wantErr     bool
	}{
		{
			name:        "success",
			statusCode:  http.StatusOK,
			body:        "file content here",
			contentType: "application/pdf",
			wantErr:     false,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			body:       `{"error": "not found"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			body:       `{"error": "internal error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/v3/grants/grant-123/attachments/attach-1/download" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				if tt.contentType != "" {
					w.Header().Set("Content-Type", tt.contentType)
				}
				w.Header().Set("Content-Length", "17")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Attachments.Download(context.Background(), "grant-123", "attach-1", "msg-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if resp.ContentType != tt.contentType {
					t.Errorf("expected content type %s, got %s", tt.contentType, resp.ContentType)
				}
				content, _ := io.ReadAll(resp.Content)
				_ = resp.Content.Close()
				if string(content) != tt.body {
					t.Errorf("expected body %s, got %s", tt.body, string(content))
				}
			}
		})
	}
}
