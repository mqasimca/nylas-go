//go:build integration

// Connectors Integration Tests Coverage:
//   - List, Get ✓
//
// Intentionally NOT tested (safety reasons):
//   - Create: Would create new connector configuration
//   - Update: Could break existing connector settings
//   - Delete: Would remove connector and break provider access
//
// These operations are tested via unit tests with mocks.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/connectors"
)

func TestConnectors_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	resp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	t.Logf("Found %d connectors", len(resp.Data))

	for _, conn := range resp.Data {
		if conn.Provider == "" {
			t.Error("Connector provider should not be empty")
		}
		t.Logf("  - Provider: %s, Name: %s", conn.Provider, conn.Name)
	}
}

func TestConnectors_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	limit := 5
	resp, err := client.Connectors.List(ctx, &connectors.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) > limit {
		t.Errorf("List() returned %d connectors, want <= %d", len(resp.Data), limit)
	}

	t.Logf("Found %d connectors (limit %d)", len(resp.Data), limit)
}

func TestConnectors_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list connectors to find available providers
	listResp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(listResp.Data) == 0 {
		t.Skip("No connectors found")
	}

	// Get each connector by provider
	for _, conn := range listResp.Data {
		t.Run(string(conn.Provider), func(t *testing.T) {
			got, err := client.Connectors.Get(ctx, conn.Provider)
			if err != nil {
				t.Fatalf("Get(%s) error = %v", conn.Provider, err)
			}

			if got.Provider != conn.Provider {
				t.Errorf("Get() provider = %s, want %s", got.Provider, conn.Provider)
			}

			t.Logf("Got connector: %s (name: %s)", got.Provider, got.Name)

			if got.Scope != nil {
				t.Logf("  Scopes: %v", got.Scope)
			}
		})
	}
}

func TestConnectors_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	iter := client.Connectors.ListAll(ctx, nil)
	all, err := iter.Collect()
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	t.Logf("ListAll() returned %d connectors", len(all))

	for _, conn := range all {
		if conn.Provider == "" {
			t.Error("Connector provider should not be empty")
		}
	}
}

// Note: TestConnectors_Create, TestConnectors_Update, TestConnectors_Delete
// are intentionally skipped to avoid breaking application configuration.
// These operations are tested via unit tests.
