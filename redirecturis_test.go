package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/redirecturis"
)

func TestRedirectURIsService_List(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantErr    bool
		wantCount  int
	}{
		{
			name: "success",
			response: `{
				"request_id": "req-123",
				"data": [
					{"id": "uri-1", "url": "https://example.com/callback", "platform": "web"},
					{"id": "uri-2", "url": "myapp://callback", "platform": "ios"}
				]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantCount:  2,
		},
		{
			name:       "unauthorized",
			response:   `{"error": {"type": "unauthorized", "message": "Invalid API key"}}`,
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Method = %s, want GET", r.Method)
				}
				if r.URL.Path != "/v3/applications/redirect-uris" {
					t.Errorf("Path = %s, want /v3/applications/redirect-uris", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.RedirectURIs.List(context.Background(), nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(resp.Data) != tt.wantCount {
				t.Errorf("Count = %d, want %d", len(resp.Data), tt.wantCount)
			}
		})
	}
}

func TestRedirectURIsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		uriID      string
		response   string
		statusCode int
		wantErr    bool
		wantURL    string
	}{
		{
			name:  "success",
			uriID: "uri-123",
			response: `{
				"request_id": "req-123",
				"data": {"id": "uri-123", "url": "https://example.com/callback", "platform": "web"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantURL:    "https://example.com/callback",
		},
		{
			name:       "not found",
			uriID:      "uri-404",
			response:   `{"error": {"type": "not_found", "message": "Redirect URI not found"}}`,
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Method = %s, want GET", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			uri, err := client.RedirectURIs.Get(context.Background(), tt.uriID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && uri.URL != tt.wantURL {
				t.Errorf("URL = %s, want %s", uri.URL, tt.wantURL)
			}
		})
	}
}

func TestRedirectURIsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/applications/redirect-uris" {
			t.Errorf("Path = %s, want /v3/applications/redirect-uris", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"id": "uri-new", "url": "https://new.example.com/callback", "platform": "web"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	uri, err := client.RedirectURIs.Create(context.Background(), &redirecturis.CreateRequest{
		URL:      "https://new.example.com/callback",
		Platform: "web",
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if uri.URL != "https://new.example.com/callback" {
		t.Errorf("URL = %s, want https://new.example.com/callback", uri.URL)
	}
}

func TestRedirectURIsService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v3/applications/redirect-uris/uri-123" {
			t.Errorf("Path = %s, want /v3/applications/redirect-uris/uri-123", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"id": "uri-123", "url": "https://updated.example.com/callback", "platform": "web"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	uri, err := client.RedirectURIs.Update(context.Background(), "uri-123", &redirecturis.UpdateRequest{
		URL: Ptr("https://updated.example.com/callback"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if uri.URL != "https://updated.example.com/callback" {
		t.Errorf("URL = %s, want https://updated.example.com/callback", uri.URL)
	}
}

func TestRedirectURIsService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v3/applications/redirect-uris/uri-123" {
			t.Errorf("Path = %s, want /v3/applications/redirect-uris/uri-123", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.RedirectURIs.Delete(context.Background(), "uri-123")

	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}
