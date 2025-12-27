package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/grants"
)

func TestGrantsService_List(t *testing.T) {
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
						"id": "grant-1",
						"provider": "google",
						"email": "user@example.com",
						"grant_status": "valid"
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
				if r.URL.Path != "/v3/grants" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Grants.List(context.Background(), nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(resp.Data) == 0 {
				t.Error("expected grants in response")
			}
		})
	}
}

func TestGrantsService_Get(t *testing.T) {
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
					"id": "grant-1",
					"provider": "google",
					"email": "user@example.com",
					"grant_status": "valid",
					"scope": ["email", "calendar"]
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
				if r.URL.Path != "/v3/grants/grant-1" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			grant, err := client.Grants.Get(context.Background(), "grant-1")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && grant.ID != "grant-1" {
				t.Errorf("expected grant-1, got %s", grant.ID)
			}
		})
	}
}

func TestGrantsService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		var req grants.UpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "grant-1",
				"provider": "google",
				"email": "user@example.com",
				"grant_status": "valid"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	grant, err := client.Grants.Update(context.Background(), "grant-1", &grants.UpdateRequest{
		Settings: map[string]any{"key": "value"},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if grant.ID != "grant-1" {
		t.Errorf("expected grant-1, got %s", grant.ID)
	}
}

func TestGrantsService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Grants.Delete(context.Background(), "grant-1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGrantsService_ListWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("provider") != "google" {
			t.Errorf("expected provider=google, got %s", r.URL.Query().Get("provider"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": [{"id": "grant-1", "provider": "google"}],
			"next_cursor": ""
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	limit := 10
	provider := "google"
	resp, err := client.Grants.List(context.Background(), &grants.ListOptions{
		Provider: &provider,
		Limit:    &limit,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 grant, got %d", len(resp.Data))
	}
}

func TestGrantsService_ListAll(t *testing.T) {
	// Grants uses offset-based pagination, not cursor-based
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page++
		offset := r.URL.Query().Get("offset")
		w.WriteHeader(http.StatusOK)
		if offset == "" || offset == "0" {
			// First page: return 2 grants (matching limit of 2)
			_, _ = w.Write([]byte(`{
				"request_id": "req-1",
				"data": [{"id": "grant-1"}, {"id": "grant-2"}]
			}`))
		} else {
			// Second page: return 1 grant (less than limit means last page)
			_, _ = w.Write([]byte(`{
				"request_id": "req-2",
				"data": [{"id": "grant-3"}]
			}`))
		}
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	limit := 2
	iter := client.Grants.ListAll(context.Background(), &grants.ListOptions{Limit: &limit})
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}
