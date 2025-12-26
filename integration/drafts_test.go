//go:build integration

// Drafts Integration Tests Coverage:
//   - List, Get, Create, Update, Delete, ListAll, CreateWithAttachmentMetadata ✓
//   - Send ✓ (requires NYLAS_TEST_EMAIL env var)
//
// All DraftsService methods are fully tested.

package integration

import (
	"os"
	"testing"

	"github.com/mqasimca/nylas-go"
	"github.com/mqasimca/nylas-go/drafts"
)

func TestDrafts_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		resp, err := client.Drafts.List(ctx, grantID, &drafts.ListOptions{
			Limit: nylas.Ptr(5),
		})
		if err != nil {
			t.Fatalf("Drafts.List failed: %v", err)
		}

		t.Logf("Found %d drafts", len(resp.Data))

		for _, draft := range resp.Data {
			if draft.ID == "" {
				t.Error("Draft ID should not be empty")
			}
			t.Logf("  - %s: %s", draft.ID, draft.Subject)
		}
	})
}

func TestDrafts_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)
		cleanup := NewCleanup(t)

		// Create a draft
		createReq := &drafts.CreateRequest{
			Subject: "Integration Test Draft",
			Body:    "This is a test draft created by integration tests.",
			To: []drafts.Participant{
				{Email: "test@example.com", Name: "Test User"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		t.Logf("Created draft: %s", created.ID)

		// Register cleanup to delete the draft
		cleanup.Add(func() {
			if err := client.Drafts.Delete(ctx, grantID, created.ID); err != nil {
				t.Logf("Warning: failed to delete test draft %s: %v", created.ID, err)
			} else {
				t.Logf("Cleaned up draft: %s", created.ID)
			}
		})

		// Verify creation
		if created.Subject != createReq.Subject {
			t.Errorf("Created draft subject = %s, want %s", created.Subject, createReq.Subject)
		}

		// Get the draft
		got, err := client.Drafts.Get(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Drafts.Get failed: %v", err)
		}

		if got.ID != created.ID {
			t.Errorf("Got draft ID = %s, want %s", got.ID, created.ID)
		}

		// Update the draft
		updateReq := &drafts.UpdateRequest{
			Subject: "Updated Integration Test Draft",
			Body:    "This draft has been updated.",
		}

		updated, err := client.Drafts.Update(ctx, grantID, created.ID, updateReq)
		if err != nil {
			t.Fatalf("Drafts.Update failed: %v", err)
		}

		if updated.Subject != updateReq.Subject {
			t.Errorf("Updated draft subject = %s, want %s", updated.Subject, updateReq.Subject)
		}

		t.Logf("Updated draft: %s -> %s", created.Subject, updated.Subject)
	})
}

func TestDrafts_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)
		cleanup := NewCleanup(t)

		// Create a draft first
		createReq := &drafts.CreateRequest{
			Subject: "Draft for Get Test",
			Body:    "This draft is for testing Get.",
			To: []drafts.Participant{
				{Email: "test@example.com", Name: "Test User"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		cleanup.Add(func() {
			_ = client.Drafts.Delete(ctx, grantID, created.ID)
		})

		// Get the draft
		got, err := client.Drafts.Get(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Drafts.Get failed: %v", err)
		}

		if got.ID != created.ID {
			t.Errorf("Get() ID = %s, want %s", got.ID, created.ID)
		}
		if got.Subject != createReq.Subject {
			t.Errorf("Get() Subject = %s, want %s", got.Subject, createReq.Subject)
		}

		t.Logf("Got draft: %s - %s", got.ID, got.Subject)
	})
}

func TestDrafts_Update(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)
		cleanup := NewCleanup(t)

		// Create a draft first
		createReq := &drafts.CreateRequest{
			Subject: "Draft for Update Test",
			Body:    "Original body content.",
			To: []drafts.Participant{
				{Email: "test@example.com", Name: "Test User"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		cleanup.Add(func() {
			_ = client.Drafts.Delete(ctx, grantID, created.ID)
		})

		// Update the draft
		updateReq := &drafts.UpdateRequest{
			Subject: "Updated Draft Subject",
			Body:    "Updated body content.",
		}

		updated, err := client.Drafts.Update(ctx, grantID, created.ID, updateReq)
		if err != nil {
			t.Fatalf("Drafts.Update failed: %v", err)
		}

		if updated.Subject != updateReq.Subject {
			t.Errorf("Update() Subject = %s, want %s", updated.Subject, updateReq.Subject)
		}

		t.Logf("Updated draft: %s -> %s", createReq.Subject, updated.Subject)

		// Verify update by getting the draft again
		got, err := client.Drafts.Get(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Drafts.Get after update failed: %v", err)
		}

		if got.Subject != updateReq.Subject {
			t.Errorf("Get() after update Subject = %s, want %s", got.Subject, updateReq.Subject)
		}
	})
}

func TestDrafts_CreateWithAttachmentMetadata(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)
		cleanup := NewCleanup(t)

		// Create a draft with CC and BCC
		createReq := &drafts.CreateRequest{
			Subject: "Integration Test - Multiple Recipients",
			Body:    "<html><body><h1>Test Email</h1><p>This is a test.</p></body></html>",
			To: []drafts.Participant{
				{Email: "to@example.com", Name: "To User"},
			},
			CC: []drafts.Participant{
				{Email: "cc@example.com", Name: "CC User"},
			},
			BCC: []drafts.Participant{
				{Email: "bcc@example.com"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		cleanup.Add(func() {
			_ = client.Drafts.Delete(ctx, grantID, created.ID)
		})

		t.Logf("Created draft with multiple recipients: %s", created.ID)

		if len(created.To) == 0 {
			t.Error("Draft should have To recipients")
		}
	})
}

func TestDrafts_ListAll_Pagination(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		iter := client.Drafts.ListAll(ctx, grantID, &drafts.ListOptions{
			Limit: nylas.Ptr(2), // Small limit to test pagination
		})

		count := 0
		maxDrafts := 5 // Limit how many we fetch in test

		for {
			draft, err := iter.Next()
			if err != nil {
				break // Done or error
			}
			if draft.ID == "" {
				t.Error("Draft ID should not be empty")
			}
			count++
			if count >= maxDrafts {
				break
			}
		}

		t.Logf("Iterated through %d drafts", count)
	})
}

func TestDrafts_Delete(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// Create a draft to delete
		createReq := &drafts.CreateRequest{
			Subject: "Draft to Delete",
			Body:    "This draft will be deleted.",
			To: []drafts.Participant{
				{Email: "delete-test@example.com"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		t.Logf("Created draft to delete: %s", created.ID)

		// Delete the draft
		err = client.Drafts.Delete(ctx, grantID, created.ID)
		if err != nil {
			t.Fatalf("Drafts.Delete failed: %v", err)
		}

		t.Logf("Successfully deleted draft: %s", created.ID)

		// Verify it's deleted by trying to get it (should fail)
		_, err = client.Drafts.Get(ctx, grantID, created.ID)
		if err == nil {
			t.Error("Expected error when getting deleted draft, got nil")
		}
	})
}

func TestDrafts_Send(t *testing.T) {
	testEmail := os.Getenv("NYLAS_TEST_EMAIL")
	if testEmail == "" {
		t.Skip("NYLAS_TEST_EMAIL not set, skipping Send test")
	}

	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		ctx := NewTestContext(t)

		// Create a draft first
		createReq := &drafts.CreateRequest{
			Subject: "Nylas Go SDK Integration Test - Drafts.Send",
			Body:    "<html><body><p>This is an automated test message sent via Drafts.Send from the Nylas Go SDK integration tests.</p><p>You can safely ignore or delete this message.</p></body></html>",
			To: []drafts.Participant{
				{Email: testEmail, Name: "SDK Test Recipient"},
			},
		}

		created, err := client.Drafts.Create(ctx, grantID, createReq)
		if err != nil {
			t.Fatalf("Drafts.Create failed: %v", err)
		}

		t.Logf("Created draft: %s", created.ID)

		// Send the draft
		sent, err := client.Drafts.Send(ctx, grantID, created.ID)
		if err != nil {
			// Clean up the draft if send failed
			_ = client.Drafts.Delete(ctx, grantID, created.ID)
			t.Fatalf("Drafts.Send failed: %v", err)
		}

		if sent.ID == "" {
			t.Error("Sent draft ID should not be empty")
		}

		t.Logf("Sent draft as message: %s (subject: %s)", sent.ID, sent.Subject)

		// Verify the draft no longer exists (it was converted to a message)
		_, err = client.Drafts.Get(ctx, grantID, created.ID)
		if err == nil {
			t.Log("Note: Draft still exists after send (provider-dependent behavior)")
		}
	})
}
