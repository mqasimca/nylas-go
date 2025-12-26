//go:build integration

// Threads Integration Tests Coverage:
//   - List, Get, ListAll, Update, ListWithFilters ✓
//
// Intentionally NOT tested (safety reasons):
//   - Delete: Permanently deletes entire thread and all messages from user's mailbox

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go"
	"github.com/mqasimca/nylas-go/threads"
)

func TestThreads_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		resp, err := client.Threads.List(ctx, grantID, &threads.ListOptions{
			Limit: nylas.Ptr(5),
		})
		if err != nil {
			t.Fatalf("Threads.List failed: %v", err)
		}

		t.Logf("Found %d threads", len(resp.Data))

		for _, thread := range resp.Data {
			if thread.ID == "" {
				t.Error("Thread ID should not be empty")
			}
			t.Logf("  - %s: %s (%d messages)", thread.ID, thread.Subject, thread.MessageCount())
		}
	})
}

func TestThreads_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, list threads to get an ID
		resp, err := client.Threads.List(ctx, grantID, &threads.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Threads.List failed: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No threads found to test Get")
		}

		threadID := resp.Data[0].ID

		// Now get the specific thread
		thread, err := client.Threads.Get(ctx, grantID, threadID)
		if err != nil {
			t.Fatalf("Threads.Get failed: %v", err)
		}

		if thread.ID != threadID {
			t.Errorf("Thread ID = %s, want %s", thread.ID, threadID)
		}

		t.Logf("Got thread: %s - %s", thread.ID, thread.Subject)
		t.Logf("  Messages: %d, Drafts: %d", thread.MessageCount(), thread.DraftCount())
		t.Logf("  Participants: %d", len(thread.Participants))
	})
}

func TestThreads_ListAll_Pagination(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		iter := client.Threads.ListAll(ctx, grantID, &threads.ListOptions{
			Limit: nylas.Ptr(2), // Small limit to test pagination
		})

		count := 0
		maxThreads := 5 // Limit how many we fetch in test

		for {
			thread, err := iter.Next()
			if err != nil {
				break // Done or error
			}
			if thread.ID == "" {
				t.Error("Thread ID should not be empty")
			}
			count++
			if count >= maxThreads {
				break
			}
		}

		t.Logf("Iterated through %d threads", count)
	})
}

func TestThreads_Update(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// First, list threads to get an ID
		resp, err := client.Threads.List(ctx, grantID, &threads.ListOptions{
			Limit: nylas.Ptr(1),
		})
		if err != nil {
			t.Fatalf("Threads.List failed: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No threads found to test Update")
		}

		threadID := resp.Data[0].ID
		originalUnread := resp.Data[0].Unread

		// Toggle unread status
		newUnread := !originalUnread
		updated, err := client.Threads.Update(ctx, grantID, threadID, &threads.UpdateRequest{
			Unread: &newUnread,
		})
		if err != nil {
			t.Fatalf("Threads.Update failed: %v", err)
		}

		t.Logf("Updated thread %s: unread %v -> %v", threadID, originalUnread, updated.Unread)

		// Restore original state
		_, err = client.Threads.Update(ctx, grantID, threadID, &threads.UpdateRequest{
			Unread: &originalUnread,
		})
		if err != nil {
			t.Logf("Warning: failed to restore original unread state: %v", err)
		}
	})
}

func TestThreads_ListWithFilters(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// Test unread filter
		t.Run("unread filter", func(t *testing.T) {
			unread := true
			resp, err := client.Threads.List(ctx, grantID, &threads.ListOptions{
				Limit:  nylas.Ptr(5),
				Unread: &unread,
			})
			if err != nil {
				t.Fatalf("Threads.List with unread filter failed: %v", err)
			}
			t.Logf("Found %d unread threads", len(resp.Data))
		})

		// Test starred filter
		t.Run("starred filter", func(t *testing.T) {
			starred := true
			resp, err := client.Threads.List(ctx, grantID, &threads.ListOptions{
				Limit:   nylas.Ptr(5),
				Starred: &starred,
			})
			if err != nil {
				t.Fatalf("Threads.List with starred filter failed: %v", err)
			}
			t.Logf("Found %d starred threads", len(resp.Data))
		})
	})
}
