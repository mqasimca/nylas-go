package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mqasimca/nylas-go/calendars"
)

func TestCalendarsService_List(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		opts       *calendars.ListOptions
		response   string
		statusCode int
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			response:   `{"data": [{"id": "cal-1", "name": "Personal"}, {"id": "cal-2", "name": "Work"}], "request_id": "req-1"}`,
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
			opts:       &calendars.ListOptions{Limit: Ptr(10)},
			response:   `{"data": [{"id": "cal-1"}], "request_id": "req-1"}`,
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

			resp, err := client.Calendars.List(context.Background(), tt.grantID, tt.opts)

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

func TestCalendarsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		calendarID string
		response   string
		statusCode int
		wantID     string
		wantErr    bool
	}{
		{
			name:       "success",
			grantID:    "grant-123",
			calendarID: "cal-456",
			response:   `{"data": {"id": "cal-456", "name": "Personal"}, "request_id": "req-1"}`,
			statusCode: 200,
			wantID:     "cal-456",
		},
		{
			name:       "not found",
			grantID:    "grant-123",
			calendarID: "cal-missing",
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

			cal, err := client.Calendars.Get(context.Background(), tt.grantID, tt.calendarID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cal.ID != tt.wantID {
				t.Errorf("Get() ID = %s, want %s", cal.ID, tt.wantID)
			}
		})
	}
}

func TestCalendarsService_Create(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "New Calendar" {
			t.Errorf("name = %v, want New Calendar", body["name"])
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "cal-new", "name": "New Calendar"},
			"request_id": "req-1",
		})
	})

	cal, err := client.Calendars.Create(context.Background(), "grant-123", &calendars.CreateRequest{
		Name: "New Calendar",
	})

	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if cal.ID != "cal-new" {
		t.Errorf("Create() ID = %s, want cal-new", cal.ID)
	}
}

func TestCalendarsService_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %s, want PUT", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data":       map[string]any{"id": "cal-123", "name": "Updated"},
			"request_id": "req-1",
		})
	})

	cal, err := client.Calendars.Update(context.Background(), "grant-123", "cal-123", &calendars.UpdateRequest{
		Name: Ptr("Updated"),
	})

	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if cal.ID != "cal-123" {
		t.Errorf("Update() ID = %s, want cal-123", cal.ID)
	}
}

func TestCalendarsService_Delete(t *testing.T) {
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

			err := client.Calendars.Delete(context.Background(), "grant-123", "cal-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarsService_ListAll(t *testing.T) {
	page := 0
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page++
		var resp map[string]any
		if page == 1 {
			resp = map[string]any{
				"data":        []map[string]string{{"id": "cal-1"}, {"id": "cal-2"}},
				"next_cursor": "page2",
				"request_id":  "req-1",
			}
		} else {
			resp = map[string]any{
				"data":       []map[string]string{{"id": "cal-3"}},
				"request_id": "req-2",
			}
		}
		_ = json.NewEncoder(w).Encode(resp)
	})

	iter := client.Calendars.ListAll(context.Background(), "grant-123", nil)
	all, err := iter.Collect()

	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("Collect() count = %d, want 3", len(all))
	}
}

func TestCalendarsService_Availability(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/calendars/availability" {
			t.Errorf("Path = %s, want /v3/calendars/availability", r.URL.Path)
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"time_slots": []map[string]any{
					{"start_time": 1700000000, "end_time": 1700003600},
				},
			},
			"request_id": "req-1",
		})
	})

	resp, err := client.Calendars.Availability(context.Background(), &calendars.AvailabilityRequest{
		StartTime:       1700000000,
		EndTime:         1700100000,
		DurationMinutes: 60,
		Participants: []calendars.AvailabilityParticipant{
			{Email: "test@example.com"},
		},
	})

	if err != nil {
		t.Fatalf("Availability() error = %v", err)
	}
	if len(resp.TimeSlots) != 1 {
		t.Errorf("Availability() slots = %d, want 1", len(resp.TimeSlots))
	}
}

func TestCalendarsService_FreeBusy(t *testing.T) {
	tests := []struct {
		name       string
		grantID    string
		request    *calendars.FreeBusyRequest
		response   string
		statusCode int
		wantCount  int
		wantSlots  int
		wantErr    bool
	}{
		{
			name:    "success with no busy slots",
			grantID: "grant-123",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700000000,
				EndTime:   1700100000,
				Emails:    []string{"test@example.com"},
			},
			response:   `{"data": [{"email": "test@example.com", "time_slots": []}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  1,
			wantSlots:  0,
		},
		{
			name:    "success with busy slots",
			grantID: "grant-123",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700000000,
				EndTime:   1700100000,
				Emails:    []string{"busy@example.com"},
			},
			response:   `{"data": [{"email": "busy@example.com", "time_slots": [{"start_time": 1700010000, "end_time": 1700013600, "status": "busy"}, {"start_time": 1700020000, "end_time": 1700023600, "status": "busy"}]}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  1,
			wantSlots:  2,
		},
		{
			name:    "multiple emails",
			grantID: "grant-123",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700000000,
				EndTime:   1700100000,
				Emails:    []string{"user1@example.com", "user2@example.com"},
			},
			response:   `{"data": [{"email": "user1@example.com", "time_slots": [{"start_time": 1700010000, "end_time": 1700013600}]}, {"email": "user2@example.com", "time_slots": []}], "request_id": "req-1"}`,
			statusCode: 200,
			wantCount:  2,
			wantSlots:  1,
		},
		{
			name:    "bad request",
			grantID: "grant-123",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700100000,
				EndTime:   1700000000, // End before start
				Emails:    []string{"test@example.com"},
			},
			response:   `{"message": "invalid time range", "type": "error"}`,
			statusCode: 400,
			wantErr:    true,
		},
		{
			name:    "unauthorized",
			grantID: "grant-123",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700000000,
				EndTime:   1700100000,
				Emails:    []string{"test@example.com"},
			},
			response:   `{"message": "unauthorized", "type": "error"}`,
			statusCode: 401,
			wantErr:    true,
		},
		{
			name:    "not found",
			grantID: "grant-missing",
			request: &calendars.FreeBusyRequest{
				StartTime: 1700000000,
				EndTime:   1700100000,
				Emails:    []string{"test@example.com"},
			},
			response:   `{"message": "grant not found", "type": "error"}`,
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
				expectedPath := "/v3/grants/" + tt.grantID + "/calendars/free-busy"
				if r.URL.Path != expectedPath {
					t.Errorf("Path = %s, want %s", r.URL.Path, expectedPath)
				}

				// Verify request body
				var body calendars.FreeBusyRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}
				if body.StartTime != tt.request.StartTime {
					t.Errorf("Request StartTime = %d, want %d", body.StartTime, tt.request.StartTime)
				}
				if body.EndTime != tt.request.EndTime {
					t.Errorf("Request EndTime = %d, want %d", body.EndTime, tt.request.EndTime)
				}
				if len(body.Emails) != len(tt.request.Emails) {
					t.Errorf("Request Emails count = %d, want %d", len(body.Emails), len(tt.request.Emails))
				}

				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			resp, err := client.Calendars.FreeBusy(context.Background(), tt.grantID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("FreeBusy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(resp) != tt.wantCount {
				t.Errorf("FreeBusy() count = %d, want %d", len(resp), tt.wantCount)
			}

			// Check total slots across all responses
			totalSlots := 0
			for _, r := range resp {
				totalSlots += len(r.TimeSlots)
			}
			if totalSlots != tt.wantSlots {
				t.Errorf("FreeBusy() total slots = %d, want %d", totalSlots, tt.wantSlots)
			}
		})
	}
}

func TestCalendarsService_FreeBusy_VerifyBusySlotFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"email": "test@example.com",
					"time_slots": []map[string]any{
						{
							"start_time": 1700010000,
							"end_time":   1700013600,
							"status":     "busy",
						},
					},
				},
			},
			"request_id": "req-1",
		})
	})

	resp, err := client.Calendars.FreeBusy(context.Background(), "grant-123", &calendars.FreeBusyRequest{
		StartTime: 1700000000,
		EndTime:   1700100000,
		Emails:    []string{"test@example.com"},
	})

	if err != nil {
		t.Fatalf("FreeBusy() error = %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("FreeBusy() count = %d, want 1", len(resp))
	}

	// Verify email
	if resp[0].Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", resp[0].Email)
	}

	// Verify time slots
	if len(resp[0].TimeSlots) != 1 {
		t.Fatalf("TimeSlots count = %d, want 1", len(resp[0].TimeSlots))
	}

	slot := resp[0].TimeSlots[0]
	if slot.StartTime != 1700010000 {
		t.Errorf("StartTime = %d, want 1700010000", slot.StartTime)
	}
	if slot.EndTime != 1700013600 {
		t.Errorf("EndTime = %d, want 1700013600", slot.EndTime)
	}
	if slot.Status != "busy" {
		t.Errorf("Status = %s, want busy", slot.Status)
	}
}

func TestCalendarsService_ErrorResponses(t *testing.T) {
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

			_, err := client.Calendars.List(context.Background(), "grant-123", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalendarsService_CreateErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{"bad request", 400, `{"message": "invalid name", "type": "error"}`},
		{"forbidden", 403, `{"message": "cannot create calendar", "type": "error"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			})

			_, err := client.Calendars.Create(context.Background(), "grant-123", &calendars.CreateRequest{
				Name: "Test",
			})
			if err == nil {
				t.Error("Create() expected error, got nil")
			}
		})
	}
}

func TestCalendarsService_AvailabilityErrors(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"message": "invalid time range", "type": "error"}`))
	})

	_, err := client.Calendars.Availability(context.Background(), &calendars.AvailabilityRequest{
		StartTime:       1700100000,
		EndTime:         1700000000, // End before start
		DurationMinutes: 60,
		Participants:    []calendars.AvailabilityParticipant{{Email: "test@example.com"}},
	})
	if err == nil {
		t.Error("Availability() expected error for invalid time range")
	}
}
