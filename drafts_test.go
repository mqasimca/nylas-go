package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mqasimca/nylas-go/drafts"
)

func TestDraftsService_List(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *drafts.ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			response:   `{"data": [{"id": "draft-1", "subject": "Test 1"}, {"id": "draft-2", "subject": "Test 2"}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  2,
		},
		{
			name:       "empty list",
			grantID:    "grant-123",
			response:   `{"data": [], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  0,
		},
		{
			name:       "with options",
			grantID:    "grant-123",
			opts:       &drafts.ListOptions{Limit: Ptr(10), Starred: Ptr(true)},
			response:   `{"data": [{"id": "draft-1"}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  1,
		},
		{
			name:       "unauthorized",
			grantID:    "grant-123",
			response:   `{"message": "unauthorized", "type": "error"}`,
			statusCode: 401,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Method = %s, want GET", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			resp, err := client.Drafts.List(context.Background(), tt.grantID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(resp.Data) != tt.wantCount {
				t.Errorf("List() count = %d, want %d", len(resp.Data), tt.wantCount)
			}
		})
	}
}

func TestDraftsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		draftID    string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			draftID:    "draft-456",
			response:   `{"data": {"id": "draft-456", "subject": "Test"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "draft-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			draftID:    "draft-missing",
			response:   `{"message": "not found", "type": "error"}`,
			statusCode: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			draft, err := client.Drafts.Get(context.Background(), tt.grantID, tt.draftID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && draft.ID != tt.wantID {
				t.Errorf("Get() ID = %s, want %s", draft.ID, tt.wantID)
			}
		})
	}
}

func TestDraftsService_Create(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		request    *drafts.CreateRequest
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:    "success",
			grantID: "grant-123",
			request: &drafts.CreateRequest{
				To:      []drafts.Participant{{Email: "to@example.com"}},
				Subject: "Test Draft",
				Body:    "Hello",
			},
			response:   `{"data": {"id": "draft-new", "subject": "Test Draft"}, "request_id": "req-1"}`,
			statusCode: 200,
		},
		{
			name:       "bad request",
			grantID:    "grant-123",
			request:    &drafts.CreateRequest{},
			response:   `{"message": "invalid request", "type": "error"}`,
			statusCode: 400,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Method = %s, want POST", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			_, err := client.Drafts.Create(context.Background(), tt.grantID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDraftsService_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "draft-123", "subject": "Updated"},
			"request_id": "req-1",
		})
	})

	draft, err := client.Drafts.Update(context.Background(), "grant-123", "draft-123", &drafts.UpdateRequest{
		Subject: "Updated",
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if draft.ID != "draft-123" {
		t.Errorf("Update() ID = %s, want draft-123", draft.ID)
	}
}

func TestDraftsService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"success", 200, false},
		{"not found", 404, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("Method = %s, want DELETE", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				if tt.statusCode >= 400 {
					_, _ = w.Write([]byte(`{"message": "error", "type": "error"}`))
				}
			})

			err := client.Drafts.Delete(context.Background(), "grant-123", "draft-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDraftsService_Send(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		draftID    string
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			draftID:    "draft-456",
			response:   `{"data": {"id": "msg-sent", "subject": "Test"}, "request_id": "req-1"}`,
			statusCode: 200,
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			draftID:    "draft-missing",
			response:   `{"message": "not found", "type": "error"}`,
			statusCode: 404,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Method = %s, want POST", r.Method)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			_, err := client.Drafts.Send(context.Background(), tt.grantID, tt.draftID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDraftsService_ListAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "draft-1"}, {"id": "draft-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "draft-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Drafts.ListAll(context.Background(), "grant-123", nil)
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}
