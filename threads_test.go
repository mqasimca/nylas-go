package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mqasimca/nylas-go/threads"
)

func TestThreadsService_List(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *threads.ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			response:   `{"data": [{"id": "thread-1", "subject": "Test 1"}, {"id": "thread-2", "subject": "Test 2"}], "request_id": "req-1"}`,
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
			opts:       &threads.ListOptions{Limit: Ptr(10), Unread: Ptr(true)},
			response:   `{"data": [{"id": "thread-1"}], "request_id": "req-1"}`,
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

			resp, err := client.Threads.List(context.Background(), tt.grantID, tt.opts)

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

func TestThreadsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		threadID   string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			threadID:   "thread-456",
			response:   `{"data": {"id": "thread-456", "subject": "Test"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "thread-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			threadID:   "thread-missing",
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

			thread, err := client.Threads.Get(context.Background(), tt.grantID, tt.threadID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && thread.ID != tt.wantID {
				t.Errorf("Get() ID = %s, want %s", thread.ID, tt.wantID)
			}
		})
	}
}

func TestThreadsService_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "thread-123", "unread": false},
			"request_id": "req-1",
		})
	})

	thread, err := client.Threads.Update(context.Background(), "grant-123", "thread-123", &threads.UpdateRequest{
		Unread: Ptr(false),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if thread.ID != "thread-123" {
		t.Errorf("Update() ID = %s, want thread-123", thread.ID)
	}
}

func TestThreadsService_Delete(t *testing.T) {
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

			err := client.Threads.Delete(context.Background(), "grant-123", "thread-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestThreadsService_ListAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "thread-1"}, {"id": "thread-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "thread-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Threads.ListAll(context.Background(), "grant-123", nil)
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}
