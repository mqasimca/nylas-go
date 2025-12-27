package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/connectors"
	"github.com/mqasimca/nylas-go/credentials"
)

func TestCredentialsService_List(t *testing.T) {
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
					{"id": "cred-1", "name": "Credential 1"},
					{"id": "cred-2", "name": "Credential 2"}
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
				if r.URL.Path != "/v3/connectors/google/creds" {
					t.Errorf("Path = %s, want /v3/connectors/google/creds", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Credentials.List(context.Background(), connectors.ProviderGoogle, nil)

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

func TestCredentialsService_Get(t *testing.T) {
	tests := []struct {
		name         string
		credentialID string
		response     string
		statusCode   int
		wantErr      bool
		wantName     string
	}{
		{
			name:         "success",
			credentialID: "cred-123",
			response: `{
				"request_id": "req-123",
				"data": {"id": "cred-123", "name": "My Credential"}
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantName:   "My Credential",
		},
		{
			name:         "not found",
			credentialID: "cred-404",
			response:     `{"error": {"type": "not_found", "message": "Credential not found"}}`,
			statusCode:   http.StatusNotFound,
			wantErr:      true,
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
			cred, err := client.Credentials.Get(context.Background(), connectors.ProviderGoogle, tt.credentialID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cred.Name != tt.wantName {
				t.Errorf("Name = %s, want %s", cred.Name, tt.wantName)
			}
		})
	}
}

func TestCredentialsService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/connectors/microsoft/creds" {
			t.Errorf("Path = %s, want /v3/connectors/microsoft/creds", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"id": "cred-new", "name": "New Credential"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	req := credentials.CreateMicrosoftRequest("New Credential", "client-id", "client-secret")
	cred, err := client.Credentials.Create(context.Background(), connectors.ProviderMicrosoft, req)

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if cred.Name != "New Credential" {
		t.Errorf("Name = %s, want New Credential", cred.Name)
	}
}

func TestCredentialsService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v3/connectors/google/creds/cred-123" {
			t.Errorf("Path = %s, want /v3/connectors/google/creds/cred-123", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {"id": "cred-123", "name": "Updated Credential"}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	cred, err := client.Credentials.Update(context.Background(), connectors.ProviderGoogle, "cred-123", &credentials.UpdateRequest{
		Name: Ptr("Updated Credential"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if cred.Name != "Updated Credential" {
		t.Errorf("Name = %s, want Updated Credential", cred.Name)
	}
}

func TestCredentialsService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v3/connectors/google/creds/cred-123" {
			t.Errorf("Path = %s, want /v3/connectors/google/creds/cred-123", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Credentials.Delete(context.Background(), connectors.ProviderGoogle, "cred-123")

	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}
