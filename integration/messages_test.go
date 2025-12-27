//go:build integration

// Package integration contains integration tests for the Nylas Go SDK.
//
// Messages Integration Tests Coverage:
//   - List, Get, ListAll, Update, ListScheduled, GetScheduled, Clean, ListWithFilters ✓
//   - Send ✓ (requires NYLAS_TEST_EMAIL env var, cleans up sent messages)
//
// Note: StopScheduled requires active scheduled messages; tested indirectly via unit tests

package integration

import (
	"os"
	"testing"

	nylas "github.com/mqasimca/nylas-go"
	"github.com/mqasimca/nylas-go/messages"
)

func TestMessages_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(5),
		})
		if err != nil {
			t.Fatalf("Messages.List failed: %v", err)
		}

		t.Logf("Found %d messages", len(resp.Data))

		for _, msg := range resp.Data {
			if msg.ID == "" {
				t.Error("Message ID should not be empty")
			}
			t.Logf("  - %s: %s", msg.ID, msg.Subject)
		}
	})
}

func TestMessages_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, list messages to get an ID
		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Messages.List failed: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No messages found to test Get")
		}

		messageID := resp.Data[0].ID

		// Now get the specific message
		msg, err := client.Messages.Get(ctx, grantID, messageID)
		if err != nil {
			t.Fatalf("Messages.Get failed: %v", err)
		}

		if msg.ID != messageID {
			t.Errorf("Message ID = %s, want %s", msg.ID, messageID)
		}

		t.Logf("Got message: %s - %s", msg.ID, msg.Subject)
	})
}

func TestMessages_ListAll_Pagination(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		iter := client.Messages.ListAll(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(2), // Small limit to test pagination
		})

		count := 0
		maxMessages := 5 // Limit how many we fetch in test

		for {
			msg, err := iter.Next()
			if err != nil {
				break // Done or error
			}
			if msg.ID == "" {
				t.Error("Message ID should not be empty")
			}
			count++
			if count >= maxMessages {
				break
			}
		}

		t.Logf("Iterated through %d messages", count)
	})
}

func TestMessages_Update(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, list messages to get an ID
		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Messages.List failed: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No messages found to test Update")
		}

		messageID := resp.Data[0].ID
		originalUnread := resp.Data[0].Unread

		// Toggle unread status
		newUnread := !originalUnread
		updated, err := client.Messages.Update(ctx, grantID, messageID, &messages.UpdateRequest{
			Unread: &newUnread,
		})
		if err != nil {
			t.Fatalf("Messages.Update failed: %v", err)
		}

		t.Logf("Updated message %s: unread %v -> %v", messageID, originalUnread, updated.Unread)

		// Restore original state
		_, err = client.Messages.Update(ctx, grantID, messageID, &messages.UpdateRequest{
			Unread: &originalUnread,
		})
		if err != nil {
			t.Logf("Warning: failed to restore original unread state: %v", err)
		}
	})
}

func TestMessages_ListScheduled(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		scheduled, err := client.Messages.ListScheduled(ctx, grantID)
		if err != nil {
			t.Fatalf("Messages.ListScheduled failed: %v", err)
		}

		t.Logf("Found %d scheduled messages", len(scheduled))

		for _, s := range scheduled {
			t.Logf("  - %s: %s", s.ScheduleID, s.Status)
		}
	})
}

func TestMessages_GetScheduled(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, list scheduled messages to get a schedule ID
		scheduled, err := client.Messages.ListScheduled(ctx, grantID)
		if err != nil {
			t.Fatalf("Messages.ListScheduled failed: %v", err)
		}

		if len(scheduled) == 0 {
			t.Skip("No scheduled messages found to test GetScheduled")
		}

		scheduleID := scheduled[0].ScheduleID

		// Now get the specific scheduled message
		msg, err := client.Messages.GetScheduled(ctx, grantID, scheduleID)
		if err != nil {
			t.Fatalf("Messages.GetScheduled failed: %v", err)
		}

		if msg.ScheduleID != scheduleID {
			t.Errorf("Schedule ID = %s, want %s", msg.ScheduleID, scheduleID)
		}

		t.Logf("Got scheduled message: %s - status: %s", msg.ScheduleID, msg.Status)
	})
}

func TestMessages_Clean(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, get a message to clean
		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Messages.List failed: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No messages found to test Clean")
		}

		messageID := resp.Data[0].ID

		// Clean the message
		cleanResp, err := client.Messages.Clean(ctx, grantID, &messages.CleanRequest{
			MessageID:    []string{messageID},
			IgnoreLinks:  nylas.Ptr(true),
			IgnoreImages: nylas.Ptr(true),
		})
		if err != nil {
			t.Fatalf("Messages.Clean failed: %v", err)
		}

		t.Logf("Cleaned %d messages", len(cleanResp))
		for _, cleaned := range cleanResp {
			t.Logf("  - %s: conversation length %d", cleaned.ID, len(cleaned.Conversation))
		}
	})
}

func TestMessages_ListWithFilters(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// Test unread filter
		t.Run("unread filter", func(t *testing.T) {
			unread := true
			resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
				Limit:  nylas.Ptr(5),
				Unread: &unread,
			})
			if err != nil {
				t.Fatalf("Messages.List with unread filter failed: %v", err)
			}
			t.Logf("Found %d unread messages", len(resp.Data))
		})

		// Test starred filter
		t.Run("starred filter", func(t *testing.T) {
			starred := true
			resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
				Limit:   nylas.Ptr(5),
				Starred: &starred,
			})
			if err != nil {
				t.Fatalf("Messages.List with starred filter failed: %v", err)
			}
			t.Logf("Found %d starred messages", len(resp.Data))
		})

		// Test has_attachment filter
		t.Run("has_attachment filter", func(t *testing.T) {
			hasAttachment := true
			resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
				Limit:         nylas.Ptr(5),
				HasAttachment: &hasAttachment,
			})
			if err != nil {
				t.Fatalf("Messages.List with has_attachment filter failed: %v", err)
			}
			t.Logf("Found %d messages with attachments", len(resp.Data))
		})

		// Test subject filter
		t.Run("subject filter", func(t *testing.T) {
			subject := "test"
			resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
				Limit:   nylas.Ptr(5),
				Subject: &subject,
			})
			if err != nil {
				t.Fatalf("Messages.List with subject filter failed: %v", err)
			}
			t.Logf("Found %d messages with 'test' in subject", len(resp.Data))
		})
	})
}

func TestMessages_Send(t *testing.T) {
	testEmail := os.Getenv("NYLAS_TEST_EMAIL")
	if testEmail == "" {
		t.Skip("NYLAS_TEST_EMAIL not set, skipping Send test")
	}

	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)
		cleanup := NewCleanup(t)

		// Send a test email
		sendReq := &messages.SendRequest{
			To: []messages.Participant{
				{Email: testEmail, Name: "SDK Test Recipient"},
			},
			Subject: "Nylas Go SDK Integration Test - Messages.Send",
			Body:    "<html><body><p>This is an automated test message from the Nylas Go SDK integration tests.</p><p>You can safely ignore or delete this message.</p></body></html>",
		}

		sent, err := client.Messages.Send(ctx, grantID, sendReq)
		if err != nil {
			t.Fatalf("Messages.Send failed: %v", err)
		}

		// Clean up: delete the sent message after test
		cleanup.Add(func() {
			if err := client.Messages.Delete(ctx, grantID, sent.ID); err != nil {
				t.Logf("Warning: failed to delete sent message %s: %v", sent.ID, err)
			} else {
				t.Logf("Cleaned up sent message: %s", sent.ID)
			}
		})

		if sent.ID == "" {
			t.Error("Sent message ID should not be empty")
		}

		t.Logf("Sent message: %s (subject: %s)", sent.ID, sent.Subject)

		// Verify the message was sent by checking it exists
		got, err := client.Messages.Get(ctx, grantID, sent.ID)
		if err != nil {
			t.Logf("Warning: could not verify sent message: %v", err)
		} else {
			if got.Subject != sendReq.Subject {
				t.Errorf("Sent message subject = %s, want %s", got.Subject, sendReq.Subject)
			}
		}
	})
}

func TestClient_RateLimits(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// Make an API call to populate rate limits
		_, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Messages.List failed: %v", err)
		}

		// Check rate limits (Nylas may not return rate limit headers)
		rate := client.RateLimits()
		t.Logf("Rate Limits - Limit: %d, Remaining: %d, Reset: %v",
			rate.Limit, rate.Remaining, rate.Reset)

		// Only validate if rate limits are returned
		if rate.Limit > 0 && rate.Remaining > rate.Limit {
			t.Errorf("Remaining (%d) should not exceed Limit (%d)", rate.Remaining, rate.Limit)
		}
	})
}
