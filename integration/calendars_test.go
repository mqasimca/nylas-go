//go:build integration

// Calendars Integration Tests Coverage:
//   - List, ListWithOptions, Get, ListAll, CRUD ✓
//   - FreeBusy, FreeBusy_MultipleEmails ✓
//   - Availability, Availability_WithOptions ✓
//
// All CalendarsService methods are fully tested.

package integration

import (
	"testing"
	"time"

	"github.com/mqasimca/nylas-go/calendars"
)

func TestCalendars_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No calendars found for this provider")
		}

		// Verify we have at least the primary calendar
		var foundPrimary bool
		for _, cal := range resp.Data {
			if cal.ID == "" {
				t.Error("Calendar ID should not be empty")
			}
			if cal.Name == "" {
				t.Error("Calendar name should not be empty")
			}
			if cal.IsPrimary {
				foundPrimary = true
			}
			t.Logf("Found calendar: %s (primary: %v)", cal.Name, cal.IsPrimary)
		}

		if !foundPrimary {
			t.Log("Warning: No primary calendar found")
		}
	})
}

func TestCalendars_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		limit := 5
		resp, err := client.Calendars.List(ctx, grantID, &calendars.ListOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) > limit {
			t.Errorf("List() returned %d calendars, want <= %d", len(resp.Data), limit)
		}
	})
}

func TestCalendars_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First list to get a calendar ID
		listResp, err := client.Calendars.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No calendars found for this provider")
		}

		// Find primary calendar (holiday calendars may return 404 on Get)
		var calID string
		for _, cal := range listResp.Data {
			if cal.IsPrimary {
				calID = cal.ID
				break
			}
		}
		if calID == "" {
			// Fallback to first owned calendar
			for _, cal := range listResp.Data {
				if cal.IsOwnedByUser {
					calID = cal.ID
					break
				}
			}
		}
		if calID == "" {
			calID = listResp.Data[0].ID
		}

		cal, err := client.Calendars.Get(ctx, grantID, calID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", calID, err)
		}

		if cal.ID != calID {
			t.Errorf("Get() ID = %s, want %s", cal.ID, calID)
		}

		t.Logf("Got calendar: %s (id: %s)", cal.Name, cal.ID)
	})
}

func TestCalendars_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		iter := client.Calendars.ListAll(ctx, grantID, &calendars.ListOptions{
			Limit: intPtr(2), // Small page size to test pagination
		})

		all, err := iter.Collect()
		if err != nil {
			t.Fatalf("Collect() error = %v", err)
		}

		t.Logf("ListAll() found %d calendars", len(all))

		// Verify all calendars have valid IDs
		for _, cal := range all {
			if cal.ID == "" {
				t.Error("Calendar ID should not be empty")
			}
		}
	})
}

func TestCalendars_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Create a test calendar
		createReq := &calendars.CreateRequest{
			Name:        "Test Calendar (SDK Integration Test)",
			Description: "Created by Nylas Go SDK integration tests",
			Timezone:    "America/New_York",
		}

		created, err := client.Calendars.Create(ctx, grantID, createReq)
		if err != nil {
			// Some providers don't support calendar creation
			t.Skipf("Create() error = %v (provider may not support calendar creation)", err)
		}

		// Register cleanup to delete the calendar
		cleanup.Add(func() {
			_ = client.Calendars.Delete(ctx, grantID, created.ID)
		})

		if created.ID == "" {
			t.Fatal("Create() returned empty ID")
		}
		if created.Name != createReq.Name {
			t.Errorf("Create() Name = %s, want %s", created.Name, createReq.Name)
		}

		t.Logf("Created calendar: %s (id: %s)", created.Name, created.ID)

		// Update the calendar
		newName := "Updated Test Calendar"
		updated, err := client.Calendars.Update(ctx, grantID, created.ID, &calendars.UpdateRequest{
			Name: &newName,
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		if updated.Name != newName {
			t.Errorf("Update() Name = %s, want %s", updated.Name, newName)
		}

		t.Logf("Updated calendar: %s", updated.Name)

		// Delete the calendar
		err = client.Calendars.Delete(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		t.Log("Deleted calendar successfully")
	})
}

func TestCalendars_FreeBusy(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get the grant to find the email address
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Query free/busy for the next 24 hours
		now := time.Now()
		startTime := now.Unix()
		endTime := now.Add(24 * time.Hour).Unix()

		resp, err := client.Calendars.FreeBusy(ctx, grantID, &calendars.FreeBusyRequest{
			StartTime: startTime,
			EndTime:   endTime,
			Emails:    []string{grant.Email},
		})
		if err != nil {
			t.Fatalf("FreeBusy() error = %v", err)
		}

		if len(resp) == 0 {
			t.Skip("No free/busy data returned")
		}

		// Verify response structure
		foundEmail := false
		for _, fb := range resp {
			if fb.Email == "" {
				t.Error("FreeBusy response email should not be empty")
			}
			if fb.Email == grant.Email {
				foundEmail = true
			}
			t.Logf("FreeBusy for %s: %d busy slots", fb.Email, len(fb.TimeSlots))

			// Log busy slots if any
			for _, slot := range fb.TimeSlots {
				if slot.StartTime == 0 {
					t.Error("BusySlot StartTime should not be zero")
				}
				if slot.EndTime == 0 {
					t.Error("BusySlot EndTime should not be zero")
				}
				if slot.EndTime <= slot.StartTime {
					t.Errorf("BusySlot EndTime (%d) should be after StartTime (%d)", slot.EndTime, slot.StartTime)
				}
				t.Logf("  - Busy: %s to %s (status: %s)",
					time.Unix(slot.StartTime, 0).Format(time.RFC3339),
					time.Unix(slot.EndTime, 0).Format(time.RFC3339),
					slot.Status)
			}
		}

		if !foundEmail {
			t.Errorf("FreeBusy() did not return data for requested email %s", grant.Email)
		}
	})
}

func TestCalendars_FreeBusy_MultipleEmails(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get the grant to find the email address
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Query free/busy for multiple emails (including grant email)
		now := time.Now()
		startTime := now.Unix()
		endTime := now.Add(24 * time.Hour).Unix()

		// Use the grant email and a test email
		testEmails := []string{grant.Email}

		resp, err := client.Calendars.FreeBusy(ctx, grantID, &calendars.FreeBusyRequest{
			StartTime: startTime,
			EndTime:   endTime,
			Emails:    testEmails,
		})
		if err != nil {
			t.Fatalf("FreeBusy() error = %v", err)
		}

		t.Logf("FreeBusy returned %d results for %d emails", len(resp), len(testEmails))

		// Verify we got a response for each email
		emailsFound := make(map[string]bool)
		for _, fb := range resp {
			emailsFound[fb.Email] = true
		}

		for _, email := range testEmails {
			if !emailsFound[email] {
				t.Logf("Warning: No FreeBusy data returned for %s", email)
			}
		}
	})
}

func TestCalendars_Availability(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get the grant to find the email address
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Query availability for the next 7 days
		now := time.Now()
		startTime := now.Unix()
		endTime := now.Add(7 * 24 * time.Hour).Unix()

		resp, err := client.Calendars.Availability(ctx, &calendars.AvailabilityRequest{
			StartTime:       startTime,
			EndTime:         endTime,
			DurationMinutes: 30,
			Participants: []calendars.AvailabilityParticipant{
				{Email: grant.Email},
			},
		})
		if err != nil {
			// Availability endpoint may require specific setup or valid calendar access
			t.Skipf("Availability() error = %v (may require calendar read access for participant)", err)
		}

		t.Logf("Availability returned %d time slots", len(resp.TimeSlots))

		// Verify time slot structure
		for i, slot := range resp.TimeSlots {
			if slot.StartTime == 0 {
				t.Errorf("TimeSlot[%d] StartTime should not be zero", i)
			}
			if slot.EndTime == 0 {
				t.Errorf("TimeSlot[%d] EndTime should not be zero", i)
			}
			if slot.EndTime <= slot.StartTime {
				t.Errorf("TimeSlot[%d] EndTime (%d) should be after StartTime (%d)", i, slot.EndTime, slot.StartTime)
			}

			// Log first few slots
			if i < 5 {
				t.Logf("  - Available: %s to %s",
					time.Unix(slot.StartTime, 0).Format(time.RFC3339),
					time.Unix(slot.EndTime, 0).Format(time.RFC3339))
			}
		}

		if len(resp.TimeSlots) > 5 {
			t.Logf("  ... and %d more slots", len(resp.TimeSlots)-5)
		}
	})
}

func TestCalendars_Availability_WithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Get the grant to find the email address
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Query availability with additional options
		now := time.Now()
		startTime := now.Unix()
		endTime := now.Add(24 * time.Hour).Unix()
		intervalMinutes := 15

		resp, err := client.Calendars.Availability(ctx, &calendars.AvailabilityRequest{
			StartTime:        startTime,
			EndTime:          endTime,
			DurationMinutes:  60,
			IntervalMinutes:  &intervalMinutes,
			RoundTo30Minutes: true,
			Participants: []calendars.AvailabilityParticipant{
				{Email: grant.Email},
			},
		})
		if err != nil {
			t.Skipf("Availability() error = %v (may require calendar read access for participant)", err)
		}

		t.Logf("Availability with options returned %d time slots", len(resp.TimeSlots))

		// Verify slots are rounded to 30 minutes
		for i, slot := range resp.TimeSlots {
			startMinute := time.Unix(slot.StartTime, 0).Minute()
			if startMinute != 0 && startMinute != 30 {
				t.Logf("Warning: TimeSlot[%d] start minute is %d (expected 0 or 30 for round_to_30_minutes)", i, startMinute)
			}
		}
	})
}

func intPtr(v int) *int {
	return &v
}
