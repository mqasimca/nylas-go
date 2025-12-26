package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/webhooks"
)

func TestWebhooksService_List(t *testing.T) {
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
				"data": [
					{
						"id": "webhook-1",
						"description": "Test webhook",
						"webhook_url": "https://example.com/webhook",
						"status": "active",
						"triggers": ["message.created"]
					}
				],
				"next_cursor": ""
			}`,
			wantErr: false,
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/v3/webhooks" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Webhooks.List(context.Background(), nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(resp.Data) == 0 {
				t.Error("expected webhooks in response")
			}
		})
	}
}

func TestWebhooksService_Get(t *testing.T) {
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
					"id": "webhook-1",
					"description": "Test webhook",
					"webhook_url": "https://example.com/webhook",
					"status": "active",
					"triggers": ["message.created"]
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
				if r.URL.Path != "/v3/webhooks/webhook-1" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			webhook, err := client.Webhooks.Get(context.Background(), "webhook-1")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && webhook.ID != "webhook-1" {
				t.Errorf("expected webhook-1, got %s", webhook.ID)
			}
		})
	}
}

func TestWebhooksService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req webhooks.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		if req.WebhookURL != "https://example.com/webhook" {
			t.Errorf("expected webhook URL, got %s", req.WebhookURL)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "webhook-new",
				"description": "New webhook",
				"webhook_url": "https://example.com/webhook",
				"status": "active",
				"triggers": ["message.created"],
				"webhook_secret": "secret-123"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	webhook, err := client.Webhooks.Create(context.Background(), &webhooks.CreateRequest{
		WebhookURL:   "https://example.com/webhook",
		TriggerTypes: []string{"message.created"},
		Description:  "New webhook",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if webhook.ID != "webhook-new" {
		t.Errorf("expected webhook-new, got %s", webhook.ID)
	}
}

func TestWebhooksService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "webhook-1",
				"description": "Updated webhook",
				"webhook_url": "https://example.com/webhook-updated",
				"status": "active",
				"triggers": ["message.created", "message.updated"]
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	desc := "Updated webhook"
	webhook, err := client.Webhooks.Update(context.Background(), "webhook-1", &webhooks.UpdateRequest{
		Description: &desc,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if webhook.Description != "Updated webhook" {
		t.Errorf("expected 'Updated webhook', got %s", webhook.Description)
	}
}

func TestWebhooksService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v3/webhooks/webhook-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Webhooks.Delete(context.Background(), "webhook-1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWebhooksService_RotateSecret(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v3/webhooks/rotate-secret/webhook-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"webhook_secret": "new-secret-456"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	resp, err := client.Webhooks.RotateSecret(context.Background(), "webhook-1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.WebhookSecret != "new-secret-456" {
		t.Errorf("expected new-secret-456, got %s", resp.WebhookSecret)
	}
}
