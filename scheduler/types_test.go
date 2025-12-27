package scheduler

import (
	"testing"
	"time"
)

func TestListConfigurationsOptions_Values(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		var o *ListConfigurationsOptions
		if o.Values() != nil {
			t.Error("expected nil for nil options")
		}
	})

	t.Run("empty options", func(t *testing.T) {
		o := &ListConfigurationsOptions{}
		v := o.Values()
		if len(v) != 0 {
			t.Errorf("expected empty map, got %d entries", len(v))
		}
	})

	t.Run("with all options", func(t *testing.T) {
		limit := 25
		o := &ListConfigurationsOptions{
			Limit:     &limit,
			PageToken: "next-page-token",
		}
		v := o.Values()

		if v["limit"] != 25 {
			t.Errorf("expected limit=25, got %v", v["limit"])
		}
		if v["page_token"] != "next-page-token" {
			t.Errorf("expected page_token=next-page-token, got %v", v["page_token"])
		}
	})

	t.Run("partial options - limit only", func(t *testing.T) {
		limit := 50
		o := &ListConfigurationsOptions{
			Limit: &limit,
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["limit"] != 50 {
			t.Errorf("expected limit=50, got %v", v["limit"])
		}
	})

	t.Run("partial options - page_token only", func(t *testing.T) {
		o := &ListConfigurationsOptions{
			PageToken: "token123",
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["page_token"] != "token123" {
			t.Errorf("expected page_token=token123, got %v", v["page_token"])
		}
	})
}

func TestListBookingsOptions_Values(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		var o *ListBookingsOptions
		if o.Values() != nil {
			t.Error("expected nil for nil options")
		}
	})

	t.Run("empty options", func(t *testing.T) {
		o := &ListBookingsOptions{}
		v := o.Values()
		if len(v) != 0 {
			t.Errorf("expected empty map, got %d entries", len(v))
		}
	})

	t.Run("with all options", func(t *testing.T) {
		limit := 10
		o := &ListBookingsOptions{
			ConfigurationID: "config123",
			Limit:           &limit,
			PageToken:       "page-token",
		}
		v := o.Values()

		if v["configuration_id"] != "config123" {
			t.Errorf("expected configuration_id=config123, got %v", v["configuration_id"])
		}
		if v["limit"] != 10 {
			t.Errorf("expected limit=10, got %v", v["limit"])
		}
		if v["page_token"] != "page-token" {
			t.Errorf("expected page_token=page-token, got %v", v["page_token"])
		}
	})

	t.Run("partial options", func(t *testing.T) {
		o := &ListBookingsOptions{
			ConfigurationID: "config456",
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["configuration_id"] != "config456" {
			t.Errorf("expected configuration_id=config456, got %v", v["configuration_id"])
		}
	})
}

func TestBooking_TimeHelpers(t *testing.T) {
	// Use a known Unix timestamp: 2024-01-15 10:30:00 UTC
	startUnix := int64(1705315800)
	endUnix := int64(1705319400) // 1 hour later

	b := &Booking{
		StartTime: startUnix,
		EndTime:   endUnix,
	}

	t.Run("StartDateTime", func(t *testing.T) {
		start := b.StartDateTime()
		expected := time.Unix(startUnix, 0)

		if !start.Equal(expected) {
			t.Errorf("StartDateTime = %v, want %v", start, expected)
		}
	})

	t.Run("EndDateTime", func(t *testing.T) {
		end := b.EndDateTime()
		expected := time.Unix(endUnix, 0)

		if !end.Equal(expected) {
			t.Errorf("EndDateTime = %v, want %v", end, expected)
		}
	})

	t.Run("duration calculation", func(t *testing.T) {
		duration := b.EndDateTime().Sub(b.StartDateTime())
		if duration != time.Hour {
			t.Errorf("expected 1 hour duration, got %v", duration)
		}
	})
}

func TestConfiguration_Structure(t *testing.T) {
	config := Configuration{
		ID: "config123",
		Participants: []Participant{
			{
				Name:        "Test User",
				Email:       "test@example.com",
				IsOrganizer: true,
			},
		},
		Availability: &Availability{
			DurationMinutes: 30,
			IntervalMinutes: 15,
		},
		EventBooking: &EventBooking{
			Title:       "Meeting",
			Description: "Test meeting",
		},
	}

	if config.ID != "config123" {
		t.Errorf("expected ID=config123, got %s", config.ID)
	}
	if len(config.Participants) != 1 {
		t.Errorf("expected 1 participant, got %d", len(config.Participants))
	}
	if config.Participants[0].IsOrganizer != true {
		t.Error("expected first participant to be organizer")
	}
	if config.Availability.DurationMinutes != 30 {
		t.Errorf("expected 30 min duration, got %d", config.Availability.DurationMinutes)
	}
}

func TestBookingRequest_Structure(t *testing.T) {
	req := BookingRequest{
		StartTime: 1705315800,
		EndTime:   1705319400,
		Guest: BookingParticipant{
			Name:  "Guest User",
			Email: "guest@example.com",
		},
		AdditionalGuests: []string{"extra@example.com"},
		AdditionalFields: map[string]string{
			"company": "Acme Inc",
		},
	}

	if req.Guest.Email != "guest@example.com" {
		t.Errorf("expected guest email, got %s", req.Guest.Email)
	}
	if len(req.AdditionalGuests) != 1 {
		t.Errorf("expected 1 additional guest, got %d", len(req.AdditionalGuests))
	}
	if req.AdditionalFields["company"] != "Acme Inc" {
		t.Errorf("expected company field, got %v", req.AdditionalFields["company"])
	}
}
