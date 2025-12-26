//go:build integration

// Events Integration Tests Coverage:
//   - List, ListWithExpandRecurring, Get, ListAll, CRUD, Update, Delete ✓
//   - CreateAllDayEvent, CreateWithParticipants ✓
//
// Intentionally NOT tested (safety reasons):
//   - SendRSVP: Would send real RSVP responses to event organizers

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/mqasimca/nylas-go/events"
)

// getTestEmail returns the test email from NYLAS_TEST_EMAIL env var.
// Used for event participant tests to avoid sending to example.com.
func getTestEmail() string {
	if email := os.Getenv("NYLAS_TEST_EMAIL"); email != "" {
		return email
	}
	return ""
}

func TestEvents_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First get a calendar ID
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		if len(calsResp.Data) == 0 {
			t.Skip("No calendars found for this provider")
		}

		// Find primary calendar or use first
		calID := calsResp.Data[0].ID
		for _, cal := range calsResp.Data {
			if cal.IsPrimary {
				calID = cal.ID
				break
			}
		}

		// List events for next 30 days
		now := time.Now()
		start := now.Unix()
		end := now.Add(30 * 24 * time.Hour).Unix()

		resp, err := client.Events.List(ctx, grantID, &events.ListOptions{
			CalendarID: calID,
			Start:      &start,
			End:        &end,
			Limit:      intPtr(10),
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		t.Logf("Found %d events in calendar %s", len(resp.Data), calID)

		for _, event := range resp.Data {
			if event.ID == "" {
				t.Error("Event ID should not be empty")
			}
			t.Logf("  - %s (id: %s)", event.Title, event.ID)
		}
	})
}

func TestEvents_ListWithExpandRecurring(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get primary calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		if len(calsResp.Data) == 0 {
			t.Skip("No calendars found")
		}

		calID := calsResp.Data[0].ID
		for _, cal := range calsResp.Data {
			if cal.IsPrimary {
				calID = cal.ID
				break
			}
		}

		now := time.Now()
		start := now.Unix()
		end := now.Add(30 * 24 * time.Hour).Unix()
		expand := true

		resp, err := client.Events.List(ctx, grantID, &events.ListOptions{
			CalendarID:      calID,
			Start:           &start,
			End:             &end,
			ExpandRecurring: &expand,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		t.Logf("Found %d events (recurring expanded)", len(resp.Data))
	})
}

func TestEvents_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get primary calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		if len(calsResp.Data) == 0 {
			t.Skip("No calendars found")
		}

		calID := calsResp.Data[0].ID
		for _, cal := range calsResp.Data {
			if cal.IsPrimary {
				calID = cal.ID
				break
			}
		}

		// Get events list
		now := time.Now()
		start := now.Add(-30 * 24 * time.Hour).Unix() // Look back 30 days
		end := now.Add(30 * 24 * time.Hour).Unix()

		listResp, err := client.Events.List(ctx, grantID, &events.ListOptions{
			CalendarID: calID,
			Start:      &start,
			End:        &end,
			Limit:      intPtr(1),
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No events found for this provider")
		}

		// Get the first event
		eventID := listResp.Data[0].ID
		event, err := client.Events.Get(ctx, grantID, eventID, calID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", eventID, err)
		}

		if event.ID != eventID {
			t.Errorf("Get() ID = %s, want %s", event.ID, eventID)
		}

		t.Logf("Got event: %s (id: %s)", event.Title, event.ID)
	})
}

func TestEvents_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get primary calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		if len(calsResp.Data) == 0 {
			t.Skip("No calendars found")
		}

		calID := calsResp.Data[0].ID
		for _, cal := range calsResp.Data {
			if cal.IsPrimary {
				calID = cal.ID
				break
			}
		}

		now := time.Now()
		start := now.Unix()
		end := now.Add(7 * 24 * time.Hour).Unix() // Next 7 days

		iter := client.Events.ListAll(ctx, grantID, &events.ListOptions{
			CalendarID: calID,
			Start:      &start,
			End:        &end,
			Limit:      intPtr(5), // Small page size
		})

		all, err := iter.Collect()
		if err != nil {
			t.Fatalf("Collect() error = %v", err)
		}

		t.Logf("ListAll() found %d events in next 7 days", len(all))
	})
}

func TestEvents_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get primary calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		if len(calsResp.Data) == 0 {
			t.Skip("No calendars found")
		}

		// Find primary or writable calendar
		var calID string
		for _, cal := range calsResp.Data {
			if cal.IsPrimary && !cal.ReadOnly {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			for _, cal := range calsResp.Data {
				if !cal.ReadOnly {
					calID = cal.ID
					break
				}
			}
		}
		if calID == "" {
			t.Skip("No writable calendar found")
		}

		// Create event 1 hour from now
		startTime := time.Now().Add(1 * time.Hour).Unix()
		endTime := time.Now().Add(2 * time.Hour).Unix()

		createReq := &events.CreateRequest{
			Title:       "Test Event (SDK Integration Test)",
			Description: "Created by Nylas Go SDK integration tests",
			Location:    "Virtual",
			When: events.When{
				Object:    "timespan",
				StartTime: &startTime,
				EndTime:   &endTime,
			},
		}

		created, err := client.Events.Create(ctx, grantID, calID, createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// Register cleanup
		cleanup.Add(func() {
			_ = client.Events.Delete(ctx, grantID, created.ID, calID)
		})

		if created.ID == "" {
			t.Fatal("Create() returned empty ID")
		}
		if created.Title != createReq.Title {
			t.Errorf("Create() Title = %s, want %s", created.Title, createReq.Title)
		}

		t.Logf("Created event: %s (id: %s)", created.Title, created.ID)

		// Update the event
		newTitle := "Updated Test Event"
		updated, err := client.Events.Update(ctx, grantID, created.ID, calID, &events.UpdateRequest{
			Title: &newTitle,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.Title != newTitle {
			t.Errorf("Update() Title = %s, want %s", updated.Title, newTitle)
		}

		t.Logf("Updated event: %s", updated.Title)

		// Delete the event
		err = client.Events.Delete(ctx, grantID, created.ID, calID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		t.Log("Deleted event successfully")
	})
}

func TestEvents_Update(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get writable calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		var calID string
		for _, cal := range calsResp.Data {
			if cal.IsPrimary && !cal.ReadOnly {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			for _, cal := range calsResp.Data {
				if !cal.ReadOnly {
					calID = cal.ID
					break
				}
			}
		}
		if calID == "" {
			t.Skip("No writable calendar found")
		}

		// Create event to update
		startTime := time.Now().Add(2 * time.Hour).Unix()
		endTime := time.Now().Add(3 * time.Hour).Unix()

		created, err := client.Events.Create(ctx, grantID, calID, &events.CreateRequest{
			Title:       "Event for Update Test",
			Description: "Original description",
			Location:    "Original location",
			When: events.When{
				Object:    "timespan",
				StartTime: &startTime,
				EndTime:   &endTime,
			},
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		cleanup.Add(func() {
			_ = client.Events.Delete(ctx, grantID, created.ID, calID)
		})

		// Update multiple fields
		newTitle := "Updated Event Title"
		newDescription := "Updated description"
		newLocation := "Updated location"

		updated, err := client.Events.Update(ctx, grantID, created.ID, calID, &events.UpdateRequest{
			Title:       &newTitle,
			Description: &newDescription,
			Location:    &newLocation,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.Title != newTitle {
			t.Errorf("Update() Title = %s, want %s", updated.Title, newTitle)
		}
		if updated.Description != newDescription {
			t.Errorf("Update() Description = %s, want %s", updated.Description, newDescription)
		}
		if updated.Location != newLocation {
			t.Errorf("Update() Location = %s, want %s", updated.Location, newLocation)
		}

		t.Logf("Updated event: %s -> %s", created.Title, updated.Title)

		// Verify by getting the event
		got, err := client.Events.Get(ctx, grantID, created.ID, calID)
		if err != nil {
			t.Fatalf("Get() after update error = %v", err)
		}

		if got.Title != newTitle {
			t.Errorf("Get() after update Title = %s, want %s", got.Title, newTitle)
		}
	})
}

func TestEvents_Delete(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get writable calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		var calID string
		for _, cal := range calsResp.Data {
			if !cal.ReadOnly {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			t.Skip("No writable calendar found")
		}

		// Create event to delete
		startTime := time.Now().Add(3 * time.Hour).Unix()
		endTime := time.Now().Add(4 * time.Hour).Unix()

		created, err := client.Events.Create(ctx, grantID, calID, &events.CreateRequest{
			Title: "Event to Delete (SDK Test)",
			When: events.When{
				Object:    "timespan",
				StartTime: &startTime,
				EndTime:   &endTime,
			},
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		t.Logf("Created event to delete: %s (id: %s)", created.Title, created.ID)

		// Delete the event
		err = client.Events.Delete(ctx, grantID, created.ID, calID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		t.Logf("Deleted event: %s", created.ID)

		// Verify deletion by trying to get the event
		// Note: Some providers (e.g., Google) may return the event with "cancelled" status
		// instead of 404 immediately after deletion
		got, err := client.Events.Get(ctx, grantID, created.ID, calID)
		if err == nil {
			if got.Status == "cancelled" {
				t.Logf("Event returned with cancelled status (expected behavior for some providers)")
			} else {
				t.Logf("Note: Event still accessible after delete (provider-dependent behavior, status: %s)", got.Status)
			}
		} else {
			t.Logf("Event correctly returns error after deletion: %v", err)
		}
	})
}

func TestEvents_CreateAllDayEvent(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get writable calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		var calID string
		for _, cal := range calsResp.Data {
			if !cal.ReadOnly {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			t.Skip("No writable calendar found")
		}

		// Create all-day event for tomorrow
		tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

		createReq := &events.CreateRequest{
			Title: "All Day Test Event (SDK)",
			When: events.When{
				Object: "date",
				Date:   tomorrow,
			},
		}

		created, err := client.Events.Create(ctx, grantID, calID, createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		cleanup.Add(func() {
			_ = client.Events.Delete(ctx, grantID, created.ID, calID)
		})

		if !created.IsAllDay() {
			t.Error("Created event should be all-day")
		}

		t.Logf("Created all-day event: %s (id: %s)", created.Title, created.ID)
	})
}

func TestEvents_CreateWithParticipants(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	testEmail := getTestEmail()
	if testEmail == "" {
		t.Skip("NYLAS_TEST_EMAIL not set, skipping participant test")
	}

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get writable calendar
		calsResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Calendars.List() error = %v", err)
		}

		var calID string
		for _, cal := range calsResp.Data {
			if !cal.ReadOnly {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			t.Skip("No writable calendar found")
		}

		startTime := time.Now().Add(24 * time.Hour).Unix()
		endTime := time.Now().Add(25 * time.Hour).Unix()

		createReq := &events.CreateRequest{
			Title: "Meeting with Participants (SDK Test)",
			When: events.When{
				Object:    "timespan",
				StartTime: &startTime,
				EndTime:   &endTime,
			},
			Participants: []events.Participant{
				{Email: testEmail, Name: "Test User"},
			},
		}

		created, err := client.Events.Create(ctx, grantID, calID, createReq)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		cleanup.Add(func() {
			_ = client.Events.Delete(ctx, grantID, created.ID, calID)
		})

		if len(created.Participants) == 0 {
			t.Log("Warning: Participants not returned (may be provider-specific)")
		}

		t.Logf("Created event with participants: %s (id: %s)", created.Title, created.ID)
	})
}
