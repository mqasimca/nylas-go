package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mqasimca/nylas-go/scheduler"
)

func TestSchedulerService_ListConfigurations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/scheduling/configurations" {
			t.Errorf("path = %s, want /v3/grants/grant-123/scheduling/configurations", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": []map[string]any{
				{"id": "config-1"},
				{"id": "config-2"},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	resp, err := client.Scheduler.ListConfigurations(context.Background(), "grant-123", nil)
	if err != nil {
		t.Fatalf("ListConfigurations() error = %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("ListConfigurations() returned %d configs, want 2", len(resp.Data))
	}
	if resp.Data[0].ID != "config-1" {
		t.Errorf("ListConfigurations() first config ID = %s, want config-1", resp.Data[0].ID)
	}
}

func TestSchedulerService_GetConfiguration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/scheduling/configurations/config-456" {
			t.Errorf("path = %s, want /v3/grants/grant-123/scheduling/configurations/config-456", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id": "config-456",
				"participants": []map[string]any{
					{"email": "organizer@example.com", "is_organizer": true},
				},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	config, err := client.Scheduler.GetConfiguration(context.Background(), "grant-123", "config-456")
	if err != nil {
		t.Fatalf("GetConfiguration() error = %v", err)
	}

	if config.ID != "config-456" {
		t.Errorf("GetConfiguration() ID = %s, want config-456", config.ID)
	}
	if len(config.Participants) != 1 {
		t.Errorf("GetConfiguration() participants = %d, want 1", len(config.Participants))
	}
}

func TestSchedulerService_CreateConfiguration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		var body scheduler.ConfigurationRequest
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body.Participants) != 1 {
			t.Errorf("body participants = %d, want 1", len(body.Participants))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id": "new-config",
				"participants": []map[string]any{
					{"email": "organizer@example.com", "is_organizer": true},
				},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	config, err := client.Scheduler.CreateConfiguration(context.Background(), "grant-123", &scheduler.ConfigurationRequest{
		Participants: []scheduler.Participant{
			{Email: "organizer@example.com", IsOrganizer: true},
		},
	})
	if err != nil {
		t.Fatalf("CreateConfiguration() error = %v", err)
	}

	if config.ID != "new-config" {
		t.Errorf("CreateConfiguration() ID = %s, want new-config", config.ID)
	}
}

func TestSchedulerService_UpdateConfiguration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/scheduling/configurations/config-456" {
			t.Errorf("path = %s, want /v3/grants/grant-123/scheduling/configurations/config-456", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"id": "config-456",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	config, err := client.Scheduler.UpdateConfiguration(context.Background(), "grant-123", "config-456", &scheduler.ConfigurationRequest{})
	if err != nil {
		t.Fatalf("UpdateConfiguration() error = %v", err)
	}

	if config.ID != "config-456" {
		t.Errorf("UpdateConfiguration() ID = %s, want config-456", config.ID)
	}
}

func TestSchedulerService_DeleteConfiguration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/v3/grants/grant-123/scheduling/configurations/config-456" {
			t.Errorf("path = %s, want /v3/grants/grant-123/scheduling/configurations/config-456", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Scheduler.DeleteConfiguration(context.Background(), "grant-123", "config-456")
	if err != nil {
		t.Fatalf("DeleteConfiguration() error = %v", err)
	}
}

func TestSchedulerService_CreateSession(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/scheduling/sessions" {
			t.Errorf("path = %s, want /v3/scheduling/sessions", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"session_id": "session-789",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	session, err := client.Scheduler.CreateSession(context.Background(), &scheduler.SessionRequest{
		ConfigurationID: "config-456",
		TimeToLive:      3600,
	})
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if session.SessionID != "session-789" {
		t.Errorf("CreateSession() SessionID = %s, want session-789", session.SessionID)
	}
}

func TestSchedulerService_ListBookings(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/v3/scheduling/configurations/config-456/bookings" {
			t.Errorf("path = %s, want /v3/scheduling/configurations/config-456/bookings", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": []map[string]any{
				{"booking_id": "booking-1", "status": "confirmed"},
				{"booking_id": "booking-2", "status": "pending"},
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	resp, err := client.Scheduler.ListBookings(context.Background(), "config-456", nil)
	if err != nil {
		t.Fatalf("ListBookings() error = %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("ListBookings() returned %d bookings, want 2", len(resp.Data))
	}
}

func TestSchedulerService_GetBooking(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"booking_id": "booking-123",
				"status":     "confirmed",
				"title":      "Test Meeting",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	booking, err := client.Scheduler.GetBooking(context.Background(), "config-456", "booking-123")
	if err != nil {
		t.Fatalf("GetBooking() error = %v", err)
	}

	if booking.BookingID != "booking-123" {
		t.Errorf("GetBooking() BookingID = %s, want booking-123", booking.BookingID)
	}
	if booking.Status != "confirmed" {
		t.Errorf("GetBooking() Status = %s, want confirmed", booking.Status)
	}
}

func TestSchedulerService_CreateBooking(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"booking_id": "new-booking",
				"status":     "pending",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	booking, err := client.Scheduler.CreateBooking(context.Background(), "config-456", &scheduler.BookingRequest{
		StartTime: 1700000000,
		EndTime:   1700003600,
		Guest:     scheduler.BookingParticipant{Email: "guest@example.com"},
	})
	if err != nil {
		t.Fatalf("CreateBooking() error = %v", err)
	}

	if booking.BookingID != "new-booking" {
		t.Errorf("CreateBooking() BookingID = %s, want new-booking", booking.BookingID)
	}
}

func TestSchedulerService_ConfirmBooking(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"request_id": "req-123",
			"data": map[string]any{
				"booking_id": "booking-123",
				"status":     "confirmed",
			},
		})
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	booking, err := client.Scheduler.ConfirmBooking(context.Background(), "config-456", "booking-123", &scheduler.ConfirmBookingRequest{
		Status: "confirmed",
	})
	if err != nil {
		t.Fatalf("ConfirmBooking() error = %v", err)
	}

	if booking.Status != "confirmed" {
		t.Errorf("ConfirmBooking() Status = %s, want confirmed", booking.Status)
	}
}

func TestSchedulerService_CancelBooking(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v3/scheduling/configurations/config-456/bookings/booking-123/cancel" {
			t.Errorf("path = %s, want /v3/scheduling/configurations/config-456/bookings/booking-123/cancel", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, _ := NewClient(WithAPIKey("test-key"), WithBaseURL(srv.URL))
	err := client.Scheduler.CancelBooking(context.Background(), "config-456", "booking-123", "Changed plans")
	if err != nil {
		t.Fatalf("CancelBooking() error = %v", err)
	}
}

func TestSchedulerService_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{"not found", http.StatusNotFound, true},
		{"unauthorized", http.StatusUnauthorized, true},
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
			_, err := client.Scheduler.GetConfiguration(context.Background(), "grant-123", "config-456")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
