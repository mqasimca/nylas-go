package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/messages"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client, err := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL), WithMaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}
	return client
}

func TestMessagesService_List(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *messages.ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			response:   `{"data": [{"id": "msg-1", "subject": "Test 1"}, {"id": "msg-2", "subject": "Test 2"}], "request_id": "req-1"}`,
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
			opts:       &messages.ListOptions{Limit: Ptr(10), Unread: Ptr(true)},
			response:   `{"data": [{"id": "msg-1"}], "request_id": "req-1"}`,
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

			resp, err := client.Messages.List(context.Background(), tt.grantID, tt.opts)

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

func TestMessagesService_Get(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		messageID  string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			messageID:  "msg-456",
			response:   `{"data": {"id": "msg-456", "subject": "Test"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "msg-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			messageID:  "msg-missing",
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

			msg, err := client.Messages.Get(context.Background(), tt.grantID, tt.messageID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && msg.ID != tt.wantID {
				t.Errorf("Get() ID = %s, want %s", msg.ID, tt.wantID)
			}
		})
	}
}

func TestMessagesService_Send(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		request    *messages.SendRequest
		response   string
		statusCode int
		wantErr    bool
	}{
		{
			name:    "success",
			grantID: "grant-123",
			request: &messages.SendRequest{
				To:      []messages.Participant{{Email: "to@example.com"}},
				Subject: "Test",
				Body:    "Hello",
			},
			response:   `{"data": {"id": "msg-new", "subject": "Test"}, "request_id": "req-1"}`,
			statusCode: 200,
		},
		{
			name:       "bad request",
			grantID:    "grant-123",
			request:    &messages.SendRequest{},
			response:   `{"message": "missing recipients", "type": "error"}`,
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

			_, err := client.Messages.Send(context.Background(), tt.grantID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessagesService_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "msg-123", "unread": false},
			"request_id": "req-1",
		})
	})

	msg, err := client.Messages.Update(context.Background(), "grant-123", "msg-123", &messages.UpdateRequest{
		Unread: Ptr(false),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if msg.ID != "msg-123" {
		t.Errorf("Update() ID = %s, want msg-123", msg.ID)
	}
}

func TestMessagesService_Delete(t *testing.T) {
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

			err := client.Messages.Delete(context.Background(), "grant-123", "msg-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessagesService_ListAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "msg-1"}, {"id": "msg-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "msg-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Messages.ListAll(context.Background(), "grant-123", nil)
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}

func TestMessagesService_ListScheduled(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// Nylas API returns wrapped format
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]string{
				{"schedule_id": "sch-1"},
				{"schedule_id": "sch-2"},
			},
			"request_id": "req-1",
		})
	})

	result, err := client.Messages.ListScheduled(context.Background(), "grant-123")

	if err != nil {
		t.Fatalf("ListScheduled() error = %v", err)
	}
	if len(result) != 2 {
		t.Errorf("ListScheduled() count = %d, want 2", len(result))
	}
}

func TestMessagesService_GetScheduled(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		scheduleID string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			scheduleID: "sch-456",
			response:   `{"data": {"schedule_id": "sch-456", "status": "pending"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "sch-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			scheduleID: "sch-missing",
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

			msg, err := client.Messages.GetScheduled(context.Background(), tt.grantID, tt.scheduleID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetScheduled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && msg.ScheduleID != tt.wantID {
				t.Errorf("GetScheduled() ScheduleID = %s, want %s", msg.ScheduleID, tt.wantID)
			}
		})
	}
}

func TestMessagesService_StopScheduled(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(200)
	})

	err := client.Messages.StopScheduled(context.Background(), "grant-123", "sch-1")

	if err != nil {
		t.Errorf("StopScheduled() error = %v", err)
	}
}

func TestMessagesService_Clean(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name: "success",
			response: `{
				"data": [
					{"id": "msg-1", "grant_id": "grant-123", "conversation": "Clean text"},
					{"id": "msg-2", "grant_id": "grant-123", "conversation": "More clean text"}
				],
				"request_id": "req-1"
			}`,
			statusCode: 200,
			wantCount:  2,
		},
		{
			name:       "unauthorized",
			response:   `{"message": "unauthorized", "type": "error"}`,
			statusCode: 401,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("Method = %s, want PUT", r.Method)
				}
				if r.URL.Path != "/v3/grants/grant-123/messages/clean" {
					t.Errorf("Path = %s, want /v3/grants/grant-123/messages/clean", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			result, err := client.Messages.Clean(context.Background(), "grant-123", &messages.CleanRequest{
				MessageID:   []string{"msg-1", "msg-2"},
				IgnoreLinks: Ptr(true),
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Clean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(result) != tt.wantCount {
				t.Errorf("Clean() count = %d, want %d", len(result), tt.wantCount)
			}
		})
	}
}
