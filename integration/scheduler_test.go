//go:build integration

// Scheduler Integration Tests Coverage:
//   - ListConfigurations, ListConfigurations_WithOptions ✓
//   - GetConfiguration ✓
//   - CreateConfiguration, UpdateConfiguration, DeleteConfiguration (CRUD) ✓
//   - CreateSession ✓
//   - ListBookings ✓
//
// Note: Booking creation/confirmation requires an active scheduling page
// which may not be available in all test environments.

package integration

import (
	"testing"
	"time"

	"github.com/mqasimca/nylas-go/scheduler"
)

func TestScheduler_ListConfigurations(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		t.Logf("Found %d scheduler configurations", len(resp.Data))

		for _, config := range resp.Data {
			if config.ID == "" {
				t.Error("Configuration ID should not be empty")
			}
			t.Logf("  - Configuration: %s (participants: %d)", config.ID, len(config.Participants))
		}
	})
}

func TestScheduler_ListConfigurations_WithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		limit := 5
		resp, err := client.Scheduler.ListConfigurations(ctx, grantID, &scheduler.ListConfigurationsOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		if len(resp.Data) > limit {
			t.Errorf("ListConfigurations() returned %d configs, want <= %d", len(resp.Data), limit)
		}

		t.Logf("Listed %d configurations (limit: %d)", len(resp.Data), limit)
	})
}

func TestScheduler_GetConfiguration(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First list to get a configuration ID
		listResp, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No scheduler configurations found for this provider")
		}

		configID := listResp.Data[0].ID

		// Get the configuration
		config, err := client.Scheduler.GetConfiguration(ctx, grantID, configID)
		if err != nil {
			t.Fatalf("GetConfiguration(%s) error = %v", configID, err)
		}

		if config.ID != configID {
			t.Errorf("GetConfiguration() ID = %s, want %s", config.ID, configID)
		}

		t.Logf("Got configuration: %s", config.ID)
		if len(config.Participants) > 0 {
			t.Logf("  - Participants: %d", len(config.Participants))
			for _, p := range config.Participants {
				t.Logf("    - %s (organizer: %v)", p.Email, p.IsOrganizer)
			}
		}
		if config.Availability != nil {
			t.Logf("  - Duration: %d minutes", config.Availability.DurationMinutes)
		}
	})
}

func TestScheduler_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get grant email for participant
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Create a test configuration
		createReq := &scheduler.ConfigurationRequest{
			Participants: []scheduler.Participant{
				{
					Email:       grant.Email,
					IsOrganizer: true,
					Availability: &scheduler.AvailabilityRules{
						OpenHours: []scheduler.OpenHours{
							{
								Days:      []int{1, 2, 3, 4, 5}, // Mon-Fri
								Timezone:  "America/New_York",
								StartTime: "09:00",
								EndTime:   "17:00",
							},
						},
					},
					Booking: &scheduler.ParticipantBooking{
						CalendarID: "primary",
					},
				},
			},
			Availability: &scheduler.Availability{
				DurationMinutes: 30,
				IntervalMinutes: 15,
			},
			EventBooking: &scheduler.EventBooking{
				Title:       "SDK Integration Test Meeting",
				Description: "Created by Nylas Go SDK integration tests",
			},
		}

		created, err := client.Scheduler.CreateConfiguration(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("CreateConfiguration() error = %v (scheduler may not be enabled)", err)
		}

		// Register cleanup
		cleanup.Add(func() {
			_ = client.Scheduler.DeleteConfiguration(ctx, grantID, created.ID)
		})

		if created.ID == "" {
			t.Fatal("CreateConfiguration() returned empty ID")
		}

		t.Logf("Created configuration: %s", created.ID)

		// Update the configuration
		updateReq := &scheduler.ConfigurationRequest{
			Participants: createReq.Participants,
			Availability: &scheduler.Availability{
				DurationMinutes: 60, // Change to 60 minutes
				IntervalMinutes: 30,
			},
			EventBooking: &scheduler.EventBooking{
				Title:       "Updated SDK Test Meeting",
				Description: "Updated by Nylas Go SDK integration tests",
			},
		}

		updated, err := client.Scheduler.UpdateConfiguration(ctx, grantID, created.ID, updateReq)
		if err != nil {
			t.Fatalf("UpdateConfiguration() error = %v", err)
		}

		if updated.Availability != nil && updated.Availability.DurationMinutes != 60 {
			t.Errorf("UpdateConfiguration() DurationMinutes = %d, want 60", updated.Availability.DurationMinutes)
		}

		t.Logf("Updated configuration: %s", updated.ID)

		// Delete the configuration
		err = client.Scheduler.DeleteConfiguration(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("DeleteConfiguration() error = %v", err)
		}

		t.Log("Deleted configuration successfully")
	})
}

func TestScheduler_CreateSession(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get grant email for participant
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// First create a configuration
		createReq := &scheduler.ConfigurationRequest{
			Participants: []scheduler.Participant{
				{
					Email:       grant.Email,
					IsOrganizer: true,
				},
			},
			Availability: &scheduler.Availability{
				DurationMinutes: 30,
			},
			EventBooking: &scheduler.EventBooking{
				Title: "Session Test Meeting",
			},
			RequiresSessionAuth: true,
		}

		config, err := client.Scheduler.CreateConfiguration(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("CreateConfiguration() error = %v (scheduler may not be enabled)", err)
		}

		cleanup.Add(func() {
			_ = client.Scheduler.DeleteConfiguration(ctx, grantID, config.ID)
		})

		// Create a session for this configuration
		session, err := client.Scheduler.CreateSession(ctx, &scheduler.SessionRequest{
			ConfigurationID: config.ID,
			TimeToLive:      300, // 5 minutes
		})
		if err != nil {
			t.Fatalf("CreateSession() error = %v", err)
		}

		if session.SessionID == "" {
			t.Error("CreateSession() returned empty SessionID")
		}

		t.Logf("Created session: %s (for config: %s)", session.SessionID, config.ID)
	})
}

func TestScheduler_ListBookings(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First get a configuration ID
		listResp, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No scheduler configurations found for this provider")
		}

		configID := listResp.Data[0].ID

		// List bookings for this configuration
		bookings, err := client.Scheduler.ListBookings(ctx, configID, nil)
		if err != nil {
			t.Fatalf("ListBookings() error = %v", err)
		}

		t.Logf("Found %d bookings for configuration %s", len(bookings.Data), configID)

		for _, booking := range bookings.Data {
			if booking.BookingID == "" {
				t.Error("Booking ID should not be empty")
			}
			startTime := time.Unix(booking.StartTime, 0)
			endTime := time.Unix(booking.EndTime, 0)
			t.Logf("  - Booking: %s (status: %s, %s - %s)",
				booking.BookingID,
				booking.Status,
				startTime.Format(time.RFC3339),
				endTime.Format(time.RFC3339))
		}
	})
}

func TestScheduler_ListBookings_WithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First get a configuration ID
		listResp, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No scheduler configurations found for this provider")
		}

		configID := listResp.Data[0].ID

		// List bookings with options
		limit := 5
		bookings, err := client.Scheduler.ListBookings(ctx, configID, &scheduler.ListBookingsOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("ListBookings() error = %v", err)
		}

		if len(bookings.Data) > limit {
			t.Errorf("ListBookings() returned %d bookings, want <= %d", len(bookings.Data), limit)
		}

		t.Logf("Listed %d bookings (limit: %d)", len(bookings.Data), limit)
	})
}

func TestScheduler_GetBooking(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First get a configuration ID
		configResp, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("ListConfigurations() error = %v", err)
		}

		if len(configResp.Data) == 0 {
			t.Skip("No scheduler configurations found for this provider")
		}

		configID := configResp.Data[0].ID

		// Get bookings
		bookingsResp, err := client.Scheduler.ListBookings(ctx, configID, nil)
		if err != nil {
			t.Fatalf("ListBookings() error = %v", err)
		}

		if len(bookingsResp.Data) == 0 {
			t.Skip("No bookings found for this configuration")
		}

		bookingID := bookingsResp.Data[0].BookingID

		// Get specific booking
		booking, err := client.Scheduler.GetBooking(ctx, configID, bookingID)
		if err != nil {
			t.Fatalf("GetBooking(%s) error = %v", bookingID, err)
		}

		if booking.BookingID != bookingID {
			t.Errorf("GetBooking() BookingID = %s, want %s", booking.BookingID, bookingID)
		}

		t.Logf("Got booking: %s (status: %s, title: %s)", booking.BookingID, booking.Status, booking.Title)
	})
}
