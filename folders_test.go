package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/folders"
)

func TestFoldersService_List(t *testing.T) {
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
						"id": "folder-1",
						"grant_id": "grant-123",
						"name": "INBOX",
						"system_folder": true
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
				if r.URL.Path != "/v3/grants/grant-123/folders" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			resp, err := client.Folders.List(context.Background(), "grant-123", nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(resp.Data) == 0 {
				t.Error("expected folders in response")
			}
		})
	}
}

func TestFoldersService_Get(t *testing.T) {
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
					"id": "folder-1",
					"grant_id": "grant-123",
					"name": "INBOX",
					"system_folder": true,
					"total_count": 100,
					"unread_count": 5
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
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			folder, err := client.Folders.Get(context.Background(), "grant-123", "folder-1")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && folder.ID != "folder-1" {
				t.Errorf("expected folder-1, got %s", folder.ID)
			}
		})
	}
}

func TestFoldersService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req folders.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		if req.Name != "My Folder" {
			t.Errorf("expected Name='My Folder', got %s", req.Name)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "folder-new",
				"grant_id": "grant-123",
				"name": "My Folder"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	folder, err := client.Folders.Create(context.Background(), "grant-123", &folders.CreateRequest{
		Name: "My Folder",
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if folder.ID != "folder-new" {
		t.Errorf("expected folder-new, got %s", folder.ID)
	}
}

func TestFoldersService_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req-123",
			"data": {
				"id": "folder-1",
				"grant_id": "grant-123",
				"name": "Renamed Folder"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	name := "Renamed Folder"
	folder, err := client.Folders.Update(context.Background(), "grant-123", "folder-1", &folders.UpdateRequest{
		Name: &name,
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if folder.Name != "Renamed Folder" {
		t.Errorf("expected 'Renamed Folder', got %s", folder.Name)
	}
}

func TestFoldersService_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"request_id": "req-123"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Folders.Delete(context.Background(), "grant-123", "folder-1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFoldersService_ListAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"request_id": "req-1",
				"data": [{"id": "folder-1", "grant_id": "grant-123", "name": "INBOX"}],
				"next_cursor": "page2"
			}`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"request_id": "req-2",
				"data": [{"id": "folder-2", "grant_id": "grant-123", "name": "Sent"}],
				"next_cursor": ""
			}`))
		}
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	iter := client.Folders.ListAll(context.Background(), "grant-123", nil)

	all, err := iter.Collect()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 folders, got %d", len(all))
	}
}
