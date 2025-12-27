package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/smartcompose"
)

func TestSmartComposeService_ComposeMessage(t *testing.T) {
	tests := []struct {
		name           string
		grantID        string
		response       string
		statusCode     int
		wantErr        bool
		wantSuggestion string
	}{
		{
			name:    "success",
			grantID: "grant-123",
			response: `{
				"request_id": "req-123",
				"data": {"suggestion": "Hello, I hope this email finds you well."}
			}`,
			statusCode:     http.StatusOK,
			wantErr:        false,
			wantSuggestion: "Hello, I hope this email finds you well.",
		},
		{
			name:       "unauthorized",
			grantID:    "grant-123",
			response:   `{"error": {"type": "unauthorized", "message": "Invalid API key"}}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Method = %s, want POST", r.Method)
				}
				expectedPath := "/v3/grants/" + tt.grantID + "/messages/smart-compose"
				if r.URL.Path != expectedPath {
					t.Errorf("Path = %s, want %s", r.URL.Path, expectedPath)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.SmartCompose.ComposeMessage(context.Background(), tt.grantID, &smartcompose.ComposeRequest{
				Prompt: "Write a professional greeting",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("ComposeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.Suggestion != tt.wantSuggestion {
				t.Errorf("Suggestion = %s, want %s", resp.Suggestion, tt.wantSuggestion)
			}
		})
	}
}

func TestSmartComposeService_ComposeReply(t *testing.T) {
	tests := []struct {
		name           string
		grantID        string
		messageID      string
		response       string
		statusCode     int
		wantErr        bool
		wantSuggestion string
	}{
		{
			name:      "success",
			grantID:   "grant-123",
			messageID: "msg-456",
			response: `{
				"request_id": "req-123",
				"data": {"suggestion": "Thank you for your email. I will review and get back to you shortly."}
			}`,
			statusCode:     http.StatusOK,
			wantErr:        false,
			wantSuggestion: "Thank you for your email. I will review and get back to you shortly.",
		},
		{
			name:       "message not found",
			grantID:    "grant-123",
			messageID:  "msg-404",
			response:   `{"error": {"type": "not_found", "message": "Message not found"}}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Method = %s, want POST", r.Method)
				}
				expectedPath := "/v3/grants/" + tt.grantID + "/messages/" + tt.messageID + "/smart-compose"
				if r.URL.Path != expectedPath {
					t.Errorf("Path = %s, want %s", r.URL.Path, expectedPath)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.SmartCompose.ComposeReply(context.Background(), tt.grantID, tt.messageID, &smartcompose.ComposeRequest{
				Prompt: "Write a polite acknowledgment",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("ComposeReply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.Suggestion != tt.wantSuggestion {
				t.Errorf("Suggestion = %s, want %s", resp.Suggestion, tt.wantSuggestion)
			}
		})
	}
}
