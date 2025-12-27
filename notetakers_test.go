package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/notetakers"
)

func TestNotetakersService_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": []map[string]any{
				{"id": "nt-1", "state": "scheduled"},
				{"id": "nt-2", "state": "completed"},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	resp, err := client.Notetakers.List(context.Background(), "grant-123", nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("List() returned %d notetakers, want 2", len(resp.Data))
	}
	if resp.Data[0].ID != "nt-1" {
		t.Errorf("List() first notetaker ID = %s, want nt-1", resp.Data[0].ID)
	}
}

func TestNotetakersService_ListWithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "scheduled" {
			t.Errorf("state = %s, want scheduled", r.URL.Query().Get("state"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit = %s, want 10", r.URL.Query().Get("limit"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data":       []map[string]any{},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	_, err := client.Notetakers.List(context.Background(), "grant-123", &notetakers.ListOptions{
		Limit: Ptr(10),
		State: "scheduled",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
}

func TestNotetakersService_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers/nt-456" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers/nt-456", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":               "nt-456",
				"name":             "Meeting Bot",
				"state":            "recording",
				"meeting_provider": "google_meet",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	nt, err := client.Notetakers.Get(context.Background(), "grant-123", "nt-456")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if nt.ID != "nt-456" {
		t.Errorf("Get() ID = %s, want nt-456", nt.ID)
	}
	if nt.Name != "Meeting Bot" {
		t.Errorf("Get() Name = %s, want Meeting Bot", nt.Name)
	}
	if nt.State != "recording" {
		t.Errorf("Get() State = %s, want recording", nt.State)
	}
}

func TestNotetakersService_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers", r.URL.Path)
		}

		var body notetakers.CreateRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.MeetingLink == "" {
			t.Error("body.MeetingLink should not be empty")
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":           "new-nt",
				"meeting_link": body.MeetingLink,
				"state":        "scheduled",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	nt, err := client.Notetakers.Create(context.Background(), "grant-123", &notetakers.CreateRequest{
		MeetingLink: "https://meet.google.com/abc-defg-hij",
		Name:        "Test Bot",
		MeetingSettings: &notetakers.MeetingSettings{
			Transcription: true,
			Summary:       true,
		},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if nt.ID != "new-nt" {
		t.Errorf("Create() ID = %s, want new-nt", nt.ID)
	}
	if nt.State != "scheduled" {
		t.Errorf("Create() State = %s, want scheduled", nt.State)
	}
}

func TestNotetakersService_CreateWithJoinTime(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body notetakers.CreateRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.JoinTime == nil {
			t.Error("body.JoinTime should not be nil")
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id":        "scheduled-nt",
				"join_time": *body.JoinTime,
				"state":     "scheduled",
			},
		})
	}))
	defer srv.Close()

	joinTime := int64(1700000000)
	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	nt, err := client.Notetakers.Create(context.Background(), "grant-123", &notetakers.CreateRequest{
		MeetingLink: "https://zoom.us/j/123456789",
		JoinTime:    &joinTime,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if nt.ID != "scheduled-nt" {
		t.Errorf("Create() ID = %s, want scheduled-nt", nt.ID)
	}
}

func TestNotetakersService_Cancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers/nt-456/cancel" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers/nt-456/cancel", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Notetakers.Cancel(context.Background(), "grant-123", "nt-456")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
}

func TestNotetakersService_Leave(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers/nt-456/leave" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers/nt-456/leave", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Notetakers.Leave(context.Background(), "grant-123", "nt-456")
	if err != nil {
		t.Fatalf("Leave() error = %v", err)
	}
}

func TestNotetakersService_GetHistory(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers/nt-456/history" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers/nt-456/history", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"events": []map[string]any{
					{"event_type": "state_changed", "created_at": 1700000000},
					{"event_type": "media_available", "created_at": 1700003600},
				},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	history, err := client.Notetakers.GetHistory(context.Background(), "grant-123", "nt-456")
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}

	if len(history.Events) != 2 {
		t.Errorf("GetHistory() returned %d events, want 2", len(history.Events))
	}
	if history.Events[0].EventType != "state_changed" {
		t.Errorf("GetHistory() first event type = %s, want state_changed", history.Events[0].EventType)
	}
}

func TestNotetakersService_GetMedia(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/notetakers/nt-456/media" {
			t.Errorf("path = %s, want /v3/grants/grant-123/notetakers/nt-456/media", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": []map[string]any{
				{"id": "media-1", "type": "video", "status": "available"},
				{"id": "media-2", "type": "transcript", "status": "available"},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	media, err := client.Notetakers.GetMedia(context.Background(), "grant-123", "nt-456")
	if err != nil {
		t.Fatalf("GetMedia() error = %v", err)
	}

	if len(media) != 2 {
		t.Errorf("GetMedia() returned %d items, want 2", len(media))
	}
	if media[0].Type != "video" {
		t.Errorf("GetMedia() first type = %s, want video", media[0].Type)
	}
}

func TestNotetakersService_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"not found", http.StatusNotFound, true},
		{"unauthorized", http.StatusUnauthorized, true},
		{"bad request", http.StatusBadRequest, true},
		{"server error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"message": tt.name,
					},
				})
			}))
			defer srv.Close()

			client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
			_, err := client.Notetakers.Get(context.Background(), "grant-123", "nt-456")
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotetaker_DateTime(t *testing.T) {
	nt := &notetakers.Notetaker{
		JoinTime:  1700000000,
		CreatedAt: 1699900000,
	}

	joinDT := nt.JoinDateTime()
	if joinDT.Unix() != 1700000000 {
		t.Errorf("JoinDateTime() = %d, want 1700000000", joinDT.Unix())
	}

	createdDT := nt.CreatedDateTime()
	if createdDT.Unix() != 1699900000 {
		t.Errorf("CreatedDateTime() = %d, want 1699900000", createdDT.Unix())
	}
}

func TestNotetakersService_ListAll(t *testing.T) {
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page++
		w.WriteHeader(http.StatusOK)
		if page == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id":  "req-1",
				"data":        []map[string]string{{"id": "nt-1"}, {"id": "nt-2"}},
				"next_cursor": "page2",
			})
		} else {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"request_id":  "req-2",
				"data":        []map[string]string{{"id": "nt-3"}},
				"next_cursor": "",
			})
		}
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	iter := client.Notetakers.ListAll(context.Background(), "grant-123", nil)
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}
