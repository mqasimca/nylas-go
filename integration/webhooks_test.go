//go:build integration

// Webhooks Integration Tests Coverage:
//   - List, ListWithOptions, Get ✓
//   - CRUD, RotateSecret (skipped when webhook URL is unreachable)
//
// Note: Nylas validates that webhook URLs are reachable before creating webhooks.
// CRUD tests will skip gracefully if using an unreachable URL.
// Set NYLAS_TEST_WEBHOOK_URL to a reachable HTTPS endpoint to test CRUD operations.

package integration

import (
	"os"
	"testing"

	"github.com/mqasimca/nylas-go/webhooks"
)

// testWebhookURL is used for webhook CRUD testing.
// Note: Nylas validates that webhook URLs are reachable, so tests will skip
// if the URL cannot be validated. To run CRUD tests, set NYLAS_TEST_WEBHOOK_URL
// to a reachable HTTPS endpoint (e.g., from webhook.site or your own server).
var testWebhookURL = getWebhookURL()

func getWebhookURL() string {
	if url := os.Getenv("NYLAS_TEST_WEBHOOK_URL"); url != "" {
		return url
	}
	return "https://example.com/nylas-sdk-test-webhook"
}

func TestWebhooks_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	resp, err := client.Webhooks.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	t.Logf("Found %d webhooks", len(resp.Data))

	for _, webhook := range resp.Data {
		if webhook.ID == "" {
			t.Error("Webhook ID should not be empty")
		}
		t.Logf("  - %s (url: %s, status: %s, triggers: %v)",
			webhook.ID, webhook.WebhookURL, webhook.Status, webhook.TriggerTypes)
	}
}

func TestWebhooks_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	limit := 5
	resp, err := client.Webhooks.List(ctx, &webhooks.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) > limit {
		t.Errorf("List() returned %d webhooks, want <= %d", len(resp.Data), limit)
	}

	t.Logf("Found %d webhooks (limit %d)", len(resp.Data), limit)
}

func TestWebhooks_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list to get a webhook ID
	listResp, err := client.Webhooks.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(listResp.Data) == 0 {
		t.Skip("No webhooks found")
	}

	webhookID := listResp.Data[0].ID

	webhook, err := client.Webhooks.Get(ctx, webhookID)
	if err != nil {
		t.Fatalf("Get(%s) error = %v", webhookID, err)
	}

	if webhook.ID != webhookID {
		t.Errorf("Get() ID = %s, want %s", webhook.ID, webhookID)
	}

	t.Logf("Got webhook: %s (url: %s, status: %s)",
		webhook.ID, webhook.WebhookURL, webhook.Status)
}

func TestWebhooks_CRUD(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	// Create a test webhook
	createReq := &webhooks.CreateRequest{
		WebhookURL:   testWebhookURL,
		TriggerTypes: []string{"message.created"},
		Description:  "SDK Integration Test Webhook",
	}

	created, err := client.Webhooks.Create(ctx, createReq)
	if err != nil {
		// Nylas validates webhook URLs are reachable - skip if URL validation fails
		t.Skipf("Create() error = %v (Nylas may require reachable webhook URL)", err)
	}

	// Register cleanup to delete the webhook
	cleanup.Add(func() {
		_ = client.Webhooks.Delete(ctx, created.ID)
	})

	if created.ID == "" {
		t.Fatal("Create() returned empty ID")
	}
	if created.WebhookURL != testWebhookURL {
		t.Errorf("Create() WebhookURL = %s, want %s", created.WebhookURL, testWebhookURL)
	}

	t.Logf("Created webhook: %s (url: %s)", created.ID, created.WebhookURL)

	// Verify webhook secret was returned
	if created.WebhookSecret == "" {
		t.Log("Warning: WebhookSecret not returned on create")
	}

	// Get the webhook
	got, err := client.Webhooks.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get(%s) error = %v", created.ID, err)
	}

	if got.ID != created.ID {
		t.Errorf("Get() ID = %s, want %s", got.ID, created.ID)
	}

	// Update the webhook
	newDescription := "Updated SDK Integration Test Webhook"
	newTriggers := []string{"message.created", "message.updated"}
	updated, err := client.Webhooks.Update(ctx, created.ID, &webhooks.UpdateRequest{
		Description:  &newDescription,
		TriggerTypes: newTriggers,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Description != newDescription {
		t.Errorf("Update() Description = %s, want %s", updated.Description, newDescription)
	}
	if len(updated.TriggerTypes) != len(newTriggers) {
		t.Errorf("Update() TriggerTypes = %v, want %v", updated.TriggerTypes, newTriggers)
	}

	t.Logf("Updated webhook: %s -> %s", createReq.Description, updated.Description)

	// Delete the webhook
	err = client.Webhooks.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	t.Log("Deleted webhook successfully")
}

func TestWebhooks_RotateSecret(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)
	cleanup := NewCleanup(t)

	// Create a test webhook
	created, err := client.Webhooks.Create(ctx, &webhooks.CreateRequest{
		WebhookURL:   testWebhookURL,
		TriggerTypes: []string{"message.created"},
		Description:  "SDK Rotate Secret Test",
	})
	if err != nil {
		// Nylas validates webhook URLs are reachable - skip if URL validation fails
		t.Skipf("Create() error = %v (Nylas may require reachable webhook URL)", err)
	}

	cleanup.Add(func() {
		_ = client.Webhooks.Delete(ctx, created.ID)
	})

	t.Logf("Created webhook: %s", created.ID)

	// Rotate the secret
	rotated, err := client.Webhooks.RotateSecret(ctx, created.ID)
	if err != nil {
		t.Fatalf("RotateSecret(%s) error = %v", created.ID, err)
	}

	if rotated.WebhookSecret == "" {
		t.Error("RotateSecret() returned empty secret")
	}

	t.Logf("Rotated secret for webhook %s (new secret length: %d)",
		created.ID, len(rotated.WebhookSecret))
}
