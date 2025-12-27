package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/connectors"
)

func TestConnectorsService_List(t *testing.T) {
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
					{"provider": "google", "name": "My Google"},
					{"provider": "microsoft", "name": "My Microsoft"}
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
				if r.URL.Path != "/v3/connectors" {
					t.Errorf("Path = %s, want /v3/connectors", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Connectors.List(context.Background(), nil)

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

func TestConnectorsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		provider   connectors.Provider
		response   string
		statusCode int
		wantErr    bool
		wantName   string
	}{
		{
			name:     "success",
			provider: connectors.ProviderGoogle,
			response: `{
				"request_id": "req-123",
				"data": {"provider": "google", "name": "My Google"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantName:   "My Google",
		},
		{
			name:       "not found",
			provider:   connectors.ProviderMicrosoft,
			response:   `{"error": {"type": "not_found", "message": "Connector not found"}}`,
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
			connector, err := client.Connectors.Get(context.Background(), tt.provider)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && connector.Name != tt.wantName {
				t.Errorf("Name = %s, want %s", connector.Name, tt.wantName)
			}
		})
	}
}

func TestConnectorsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/connectors" {
			t.Errorf("Path = %s, want /v3/connectors", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"provider": "google", "name": "New Google"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	req := connectors.CreateGoogleRequest("New Google", "client-id", "client-secret", []string{"email"})
	connector, err := client.Connectors.Create(context.Background(), req)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if connector.Name != "New Google" {
		t.Errorf("Name = %s, want New Google", connector.Name)
	}
}

func TestConnectorsService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v3/connectors/google" {
			t.Errorf("Path = %s, want /v3/connectors/google", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"provider": "google", "name": "Updated Google"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	connector, err := client.Connectors.Update(context.Background(), connectors.ProviderGoogle, &connectors.UpdateRequest{
		Name: Ptr("Updated Google"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if connector.Name != "Updated Google" {
		t.Errorf("Name = %s, want Updated Google", connector.Name)
	}
}

func TestConnectorsService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v3/connectors/google" {
			t.Errorf("Path = %s, want /v3/connectors/google", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Connectors.Delete(context.Background(), connectors.ProviderGoogle)

	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}
