package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mqasimca/nylas-go/events"
)

func TestEventsService_List(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *events.ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			opts:       &events.ListOptions{CalendarID: "cal-123"},
			response:   `{"data": [{"id": "event-1", "title": "Meeting"}, {"id": "event-2", "title": "Lunch"}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  2,
		},
		{
			name:       "empty list",
			grantID:    "grant-123",
			opts:       &events.ListOptions{CalendarID: "cal-123"},
			response:   `{"data": [], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  0,
		},
		{
			name:    "with time range",
			grantID: "grant-123",
			opts: &events.ListOptions{
				CalendarID:      "cal-123",
				Start:           Ptr(int64(1700000000)),
				End:             Ptr(int64(1700100000)),
				ExpandRecurring: Ptr(true),
			},
			response:   `{"data": [{"id": "event-1"}], "request_id": "req-1"}`,
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

			resp, err := client.Events.List(context.Background(), tt.grantID, tt.opts)

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

func TestEventsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		eventID    string
		calendarID string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			eventID:    "event-456",
			calendarID: "cal-123",
			response:   `{"data": {"id": "event-456", "title": "Meeting"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "event-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			eventID:    "event-missing",
			calendarID: "cal-123",
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

			event, err := client.Events.Get(context.Background(), tt.grantID, tt.eventID, tt.calendarID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && event.ID != tt.wantID {
				t.Errorf("Get() ID = %s, want %s", event.ID, tt.wantID)
			}
		})
	}
}

func TestEventsService_Create(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}

		// Verify calendar_id is in query params
		if r.URL.Query().Get("calendar_id") != "cal-123" {
			t.Errorf("calendar_id = %s, want cal-123", r.URL.Query().Get("calendar_id"))
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "New Meeting" {
			t.Errorf("title = %v, want New Meeting", body["title"])
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "event-new", "title": "New Meeting"},
			"request_id": "req-1",
		})
	})

	startTime := int64(1700000000)
	endTime := int64(1700003600)
	event, err := client.Events.Create(context.Background(), "grant-123", "cal-123", &events.CreateRequest{
		Title: "New Meeting",
		When: events.When{
			Object:    "timespan",
			StartTime: &startTime,
			EndTime:   &endTime,
		},
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if event.ID != "event-new" {
		t.Errorf("Create() ID = %s, want event-new", event.ID)
	}
}

func TestEventsService_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}

		// Verify calendar_id is in query params
		if r.URL.Query().Get("calendar_id") != "cal-123" {
			t.Errorf("calendar_id = %s, want cal-123", r.URL.Query().Get("calendar_id"))
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "event-123", "title": "Updated"},
			"request_id": "req-1",
		})
	})

	event, err := client.Events.Update(context.Background(), "grant-123", "event-123", "cal-123", &events.UpdateRequest{
		Title: Ptr("Updated"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if event.ID != "event-123" {
		t.Errorf("Update() ID = %s, want event-123", event.ID)
	}
}

func TestEventsService_Delete(t *testing.T) {
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
				// Verify calendar_id is in query params
				if r.URL.Query().Get("calendar_id") != "cal-123" {
					t.Errorf("calendar_id = %s, want cal-123", r.URL.Query().Get("calendar_id"))
				}
				w.WriteHeader(tt.statusCode)
				if tt.statusCode >= 400 {
					_, _ = w.Write([]byte(`{"message": "error", "type": "error"}`))
				}
			})

			err := client.Events.Delete(context.Background(), "grant-123", "event-123", "cal-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventsService_ListAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "event-1"}, {"id": "event-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "event-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Events.ListAll(context.Background(), "grant-123", &events.ListOptions{CalendarID: "cal-123"})
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}

func TestEventsService_SendRSVP(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/events/event-123/send-rsvp" {
			t.Errorf("Path = %s, want /v3/grants/grant-123/events/event-123/send-rsvp", r.URL.Path)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["status"] != "yes" {
			t.Errorf("status = %v, want yes", body["status"])
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{"request_id": "req-1"})
	})

	err := client.Events.SendRSVP(context.Background(), "grant-123", "event-123", "cal-123", &events.RSVPRequest{
		Status: "yes",
	})

	if err != nil {
		t.Fatalf("SendRSVP() error = %v", err)
	}
}

func TestEvent_Helpers(t *testing.T) {
	startTime := int64(1700000000)
	endTime := int64(1700003600)
	event := &events.Event{
		ID:        "event-1",
		CreatedAt: 1699999999,
		UpdatedAt: 1700000001,
		When: events.When{
			Object:    "timespan",
			StartTime: &startTime,
			EndTime:   &endTime,
		},
	}

	if event.StartDateTime().Unix() != startTime {
		t.Errorf("StartDateTime() = %d, want %d", event.StartDateTime().Unix(), startTime)
	}
	if event.EndDateTime().Unix() != endTime {
		t.Errorf("EndDateTime() = %d, want %d", event.EndDateTime().Unix(), endTime)
	}
	if event.IsAllDay() {
		t.Error("IsAllDay() = true, want false")
	}

	// Test all-day event
	allDayEvent := &events.Event{
		When: events.When{Object: "date", Date: "2023-11-15"},
	}
	if !allDayEvent.IsAllDay() {
		t.Error("IsAllDay() = false, want true")
	}

	// Test recurring event
	recurringEvent := &events.Event{
		Recurrence: &events.Recurrence{RRule: "RRULE:FREQ=WEEKLY"},
	}
	if !recurringEvent.IsRecurring() {
		t.Error("IsRecurring() = false, want true")
	}
}

func TestEventsService_ErrorResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{"bad request", 400, `{"message": "bad request", "type": "error"}`, true},
		{"unauthorized", 401, `{"message": "unauthorized", "type": "error"}`, true},
		{"forbidden", 403, `{"message": "forbidden", "type": "error"}`, true},
		{"not found", 404, `{"message": "not found", "type": "error"}`, true},
		{"rate limited", 429, `{"message": "rate limited", "type": "error"}`, true},
		{"server error", 500, `{"message": "server error", "type": "error"}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			_, err := client.Events.List(context.Background(), "grant-123", &events.ListOptions{CalendarID: "cal-123"})
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventsService_CreateErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{"bad request", 400, `{"message": "invalid event", "type": "error"}`},
		{"forbidden", 403, `{"message": "cannot create event", "type": "error"}`},
		{"conflict", 409, `{"message": "event conflict", "type": "error"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			startTime := int64(1700000000)
			endTime := int64(1700003600)
			_, err := client.Events.Create(context.Background(), "grant-123", "cal-123", &events.CreateRequest{
				Title: "Test",
				When: events.When{
					Object:    "timespan",
					StartTime: &startTime,
					EndTime:   &endTime,
				},
			})
			if err == nil {
				t.Error("Create() expected error, got nil")
			}
		})
	}
}

func TestEventsService_UpdateErrors(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"message": "event not found", "type": "error"}`))
	})

	_, err := client.Events.Update(context.Background(), "grant-123", "event-missing", "cal-123", &events.UpdateRequest{
		Title: Ptr("Updated"),
	})
	if err == nil {
		t.Error("Update() expected error for missing event")
	}
}

func TestEventsService_SendRSVPErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{"bad request", 400, `{"message": "invalid status", "type": "error"}`},
		{"not found", 404, `{"message": "event not found", "type": "error"}`},
		{"forbidden", 403, `{"message": "cannot RSVP", "type": "error"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			err := client.Events.SendRSVP(context.Background(), "grant-123", "event-123", "cal-123", &events.RSVPRequest{
				Status: "yes",
			})
			if err == nil {
				t.Error("SendRSVP() expected error, got nil")
			}
		})
	}
}

func TestEventsService_Import(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *events.ImportOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			opts:       &events.ImportOptions{CalendarID: "cal-123"},
			response:   `{"data": [{"id": "event-1", "title": "Meeting"}, {"id": "event-2", "title": "Lunch"}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  2,
		},
		{
			name:       "empty list",
			grantID:    "grant-123",
			opts:       &events.ImportOptions{CalendarID: "cal-123"},
			response:   `{"data": [], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  0,
		},
		{
			name:    "with time range",
			grantID: "grant-123",
			opts: &events.ImportOptions{
				CalendarID: "cal-123",
				Start:      Ptr(int64(1700000000)),
				End:        Ptr(int64(1700100000)),
			},
			response:   `{"data": [{"id": "event-1"}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  1,
		},
		{
			name:       "unauthorized",
			grantID:    "grant-123",
			opts:       &events.ImportOptions{CalendarID: "cal-123"},
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
				if r.URL.Path != "/v3/grants/grant-123/events/import" {
					t.Errorf("Path = %s, want /v3/grants/grant-123/events/import", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			resp, err := client.Events.Import(context.Background(), tt.grantID, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("Import() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(resp.Data) != tt.wantCount {
				t.Errorf("Import() count = %d, want %d", len(resp.Data), tt.wantCount)
			}
		})
	}
}

func TestEventsService_ImportAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/grants/grant-123/events/import" {
			t.Errorf("Path = %s, want /v3/grants/grant-123/events/import", r.URL.Path)
		}
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "event-1"}, {"id": "event-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "event-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Events.ImportAll(context.Background(), "grant-123", &events.ImportOptions{CalendarID: "cal-123"})
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}

func TestImportOptions_Values(t *testing.T) {
	tests := []struct {
		name string
		opts *events.ImportOptions
		want map[string]any
	}{
		{
			name: "nil options",
			opts: nil,
			want: nil,
		},
		{
			name: "calendar_id only",
			opts: &events.ImportOptions{CalendarID: "cal-123"},
			want: map[string]any{"calendar_id": "cal-123"},
		},
		{
			name: "all fields",
			opts: &events.ImportOptions{
				CalendarID: "cal-123",
				Start:      Ptr(int64(1700000000)),
				End:        Ptr(int64(1700100000)),
				Limit:      Ptr(50),
				PageToken:  "token123",
			},
			want: map[string]any{
				"calendar_id": "cal-123",
				"start":       int64(1700000000),
				"end":         int64(1700100000),
				"limit":       50,
				"page_token":  "token123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Values()
			if tt.want == nil && got != nil {
				t.Errorf("Values() = %v, want nil", got)
				return
			}
			if tt.want != nil {
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("Values()[%s] = %v, want %v", k, got[k], v)
					}
				}
			}
		})
	}
}
