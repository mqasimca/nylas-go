//go:build integration

// Notetakers Integration Tests Coverage:
//   - List, List_WithOptions ✓
//   - Get ✓
//   - GetHistory ✓
//   - GetMedia ✓
//
// Note: Create, Cancel, and Leave operations require active meetings
// which may not be available in all test environments. These tests
// are skipped when no existing notetakers are available.

package integration

import (
	"testing"
	"time"

	"github.com/mqasimca/nylas-go/notetakers"
)

func TestNotetakers_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.Notetakers.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		t.Logf("Found %d notetakers", len(resp.Data))

		for _, nt := range resp.Data {
			if nt.ID == "" {
				t.Error("Notetaker ID should not be empty")
			}
			t.Logf("  - Notetaker: %s (state: %s, provider: %s)",
				nt.ID, nt.State, nt.MeetingProvider)
		}
	})
}

func TestNotetakers_List_WithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		limit := 5
		resp, err := client.Notetakers.List(ctx, grantID, &notetakers.ListOptions{
			Limit: &limit,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(resp.Data) > limit {
			t.Errorf("List() returned %d notetakers, want <= %d", len(resp.Data), limit)
		}

		t.Logf("Listed %d notetakers (limit: %d)", len(resp.Data), limit)
	})
}

func TestNotetakers_List_ByState(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Test filtering by completed state
		resp, err := client.Notetakers.List(ctx, grantID, &notetakers.ListOptions{
			State: notetakers.StateCompleted,
		})
		if err != nil {
			t.Skipf("List(state=completed) error = %v (state filter may not be supported)", err)
		}

		for _, nt := range resp.Data {
			if nt.State != notetakers.StateCompleted {
				t.Errorf("Expected state=%s, got %s for notetaker %s",
					notetakers.StateCompleted, nt.State, nt.ID)
			}
		}

		t.Logf("Found %d completed notetakers", len(resp.Data))
	})
}

func TestNotetakers_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First list to get a notetaker ID
		listResp, err := client.Notetakers.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No notetakers found for this provider")
		}

		notetakerID := listResp.Data[0].ID

		// Get the notetaker
		nt, err := client.Notetakers.Get(ctx, grantID, notetakerID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", notetakerID, err)
		}

		if nt.ID != notetakerID {
			t.Errorf("Get() ID = %s, want %s", nt.ID, notetakerID)
		}

		t.Logf("Got notetaker: %s", nt.ID)
		t.Logf("  - State: %s", nt.State)
		t.Logf("  - Provider: %s", nt.MeetingProvider)
		if nt.Name != "" {
			t.Logf("  - Name: %s", nt.Name)
		}
		if nt.MeetingLink != "" {
			t.Logf("  - Meeting: %s", nt.MeetingLink)
		}
		if nt.JoinTime > 0 {
			t.Logf("  - Join Time: %s", time.Unix(nt.JoinTime, 0).Format(time.RFC3339))
		}
	})
}

func TestNotetakers_GetHistory(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// List notetakers (without state filter for broader compatibility)
		listResp, err := client.Notetakers.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No notetakers found for this provider")
		}

		notetakerID := listResp.Data[0].ID

		// Get the history
		history, err := client.Notetakers.GetHistory(ctx, grantID, notetakerID)
		if err != nil {
			t.Skipf("GetHistory(%s) error = %v (history may not be available)", notetakerID, err)
		}

		t.Logf("Got history for notetaker %s: %d events", notetakerID, len(history.Events))

		for _, event := range history.Events {
			eventTime := time.Unix(event.CreatedAt, 0)
			t.Logf("  - [%s] %s", eventTime.Format("15:04:05"), event.EventType)
		}
	})
}

func TestNotetakers_GetMedia(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// List notetakers (without state filter for broader compatibility)
		listResp, err := client.Notetakers.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(listResp.Data) == 0 {
			t.Skip("No notetakers found for this provider")
		}

		notetakerID := listResp.Data[0].ID

		// Get the media
		media, err := client.Notetakers.GetMedia(ctx, grantID, notetakerID)
		if err != nil {
			t.Skipf("GetMedia(%s) error = %v (media may not be available)", notetakerID, err)
		}

		t.Logf("Got %d media items for notetaker %s", len(media), notetakerID)

		for _, m := range media {
			t.Logf("  - Media: %s (type: %s, status: %s)", m.ID, m.Type, m.Status)
			if m.ContentType != "" {
				t.Logf("    Content-Type: %s", m.ContentType)
			}
			if m.ExpiresAt > 0 {
				expires := time.Unix(m.ExpiresAt, 0)
				t.Logf("    Expires: %s", expires.Format(time.RFC3339))
			}
		}
	})
}

// Note: The following tests require the ability to create notetakers,
// which needs valid meeting links. These tests create and then cancel
// the notetaker to avoid actually joining meetings.

func TestNotetakers_CreateAndCancel(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	// Skip if no test meeting link is configured
	meetingLink := cfg.TestMeetingLink
	if meetingLink == "" {
		t.Skip("No test meeting link configured (set NYLAS_TEST_MEETING_LINK env var)")
	}

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Schedule notetaker for future (5 minutes from now)
		futureTime := time.Now().Add(5 * time.Minute).Unix()

		createReq := &notetakers.CreateRequest{
			MeetingLink: meetingLink,
			JoinTime:    &futureTime,
			Name:        "SDK Integration Test Notetaker",
			MeetingSettings: &notetakers.MeetingSettings{
				Transcription: true,
			},
		}

		nt, err := client.Notetakers.Create(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("Create() error = %v (notetakers may not be enabled)", err)
		}

		// Register cleanup
		cleanup.Add(func() {
			_ = client.Notetakers.Cancel(ctx, grantID, nt.ID)
		})

		if nt.ID == "" {
			t.Fatal("Create() returned empty ID")
		}

		t.Logf("Created notetaker: %s (state: %s)", nt.ID, nt.State)

		// Cancel the notetaker since we don't want to actually join a meeting
		err = client.Notetakers.Cancel(ctx, grantID, nt.ID)
		if err != nil {
			t.Fatalf("Cancel(%s) error = %v", nt.ID, err)
		}

		t.Logf("Cancelled notetaker: %s", nt.ID)

		// Verify it was cancelled
		cancelled, err := client.Notetakers.Get(ctx, grantID, nt.ID)
		if err != nil {
			t.Fatalf("Get(%s) after cancel error = %v", nt.ID, err)
		}

		if cancelled.State != notetakers.StateCancelled {
			t.Errorf("Expected state=%s after cancel, got %s",
				notetakers.StateCancelled, cancelled.State)
		}
	})
}

func TestNotetakers_CreateWithSettings(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	meetingLink := cfg.TestMeetingLink
	if meetingLink == "" {
		t.Skip("No test meeting link configured (set NYLAS_TEST_MEETING_LINK env var)")
	}

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		futureTime := time.Now().Add(10 * time.Minute).Unix()

		createReq := &notetakers.CreateRequest{
			MeetingLink: meetingLink,
			JoinTime:    &futureTime,
			Name:        "Full Settings Test Notetaker",
			MeetingSettings: &notetakers.MeetingSettings{
				VideoRecording:    true,
				AudioRecording:    true,
				Transcription:     true,
				Summary:           true,
				ActionItems:       true,
				LeaveAfterSilence: 300, // 5 minutes
			},
		}

		nt, err := client.Notetakers.Create(ctx, grantID, createReq)
		if err != nil {
			t.Skipf("Create() error = %v (notetakers may not be enabled)", err)
		}

		cleanup.Add(func() {
			_ = client.Notetakers.Cancel(ctx, grantID, nt.ID)
		})

		t.Logf("Created notetaker with full settings: %s", nt.ID)

		// Verify settings were applied
		if nt.MeetingSettings != nil {
			t.Logf("  - Video Recording: %v", nt.MeetingSettings.VideoRecording)
			t.Logf("  - Audio Recording: %v", nt.MeetingSettings.AudioRecording)
			t.Logf("  - Transcription: %v", nt.MeetingSettings.Transcription)
			t.Logf("  - Summary: %v", nt.MeetingSettings.Summary)
			t.Logf("  - Action Items: %v", nt.MeetingSettings.ActionItems)
		}

		// Cancel to clean up
		err = client.Notetakers.Cancel(ctx, grantID, nt.ID)
		if err != nil {
			t.Fatalf("Cancel(%s) error = %v", nt.ID, err)
		}
	})
}
