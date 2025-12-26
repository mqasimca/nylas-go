//go:build integration

// Grants Integration Tests Coverage:
//   - List, ListWithOptions, Get, ListByProvider, ListByEmail ✓
//
// Intentionally NOT tested (safety reasons):
//   - Update: Could break test account authentication
//   - Delete: Would revoke grant and break all other tests
//
// These operations are tested via unit tests with mocks.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/grants"
)

func TestGrants_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// Grants.List doesn't need a provider-specific grant ID
	resp, err := client.Grants.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) == 0 {
		t.Skip("No grants found")
	}

	t.Logf("Found %d grants", len(resp.Data))

	for _, grant := range resp.Data {
		if grant.ID == "" {
			t.Error("Grant ID should not be empty")
		}
		t.Logf("  - %s (provider: %s, email: %s, status: %s)",
			grant.ID, grant.Provider, grant.Email, grant.GrantStatus)
	}
}

func TestGrants_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	limit := 5
	resp, err := client.Grants.List(ctx, &grants.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) > limit {
		t.Errorf("List() returned %d grants, want <= %d", len(resp.Data), limit)
	}

	t.Logf("Found %d grants (limit %d)", len(resp.Data), limit)
}

func TestGrants_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", grantID, err)
		}

		if grant.ID != grantID {
			t.Errorf("Get() ID = %s, want %s", grant.ID, grantID)
		}

		t.Logf("Got grant: %s (provider: %s, email: %s, status: %s)",
			grant.ID, grant.Provider, grant.Email, grant.GrantStatus)

		// Verify required fields
		if grant.Provider == "" {
			t.Error("Grant Provider should not be empty")
		}
		if grant.GrantStatus == "" {
			t.Error("Grant Status should not be empty")
		}
	})
}

func TestGrants_ListByProvider(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First get all grants to find available providers
	allResp, err := client.Grants.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(allResp.Data) == 0 {
		t.Skip("No grants found")
	}

	// Get the provider of the first grant
	provider := allResp.Data[0].Provider

	// Filter by that provider
	resp, err := client.Grants.List(ctx, &grants.ListOptions{
		Provider: &provider,
	})
	if err != nil {
		t.Fatalf("List(provider=%s) error = %v", provider, err)
	}

	t.Logf("Found %d grants for provider %s", len(resp.Data), provider)

	// Verify all returned grants have the correct provider
	for _, grant := range resp.Data {
		if grant.Provider != provider {
			t.Errorf("Grant provider = %s, want %s", grant.Provider, provider)
		}
	}
}

func TestGrants_ListByEmail(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First get the grant to find its email
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", grantID, err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email")
		}

		// Filter by that email
		resp, err := client.Grants.List(ctx, &grants.ListOptions{
			Email: &grant.Email,
		})
		if err != nil {
			t.Fatalf("List(email=%s) error = %v", grant.Email, err)
		}

		if len(resp.Data) == 0 {
			t.Errorf("List(email=%s) returned no grants", grant.Email)
		}

		t.Logf("Found %d grants for email %s", len(resp.Data), grant.Email)

		// Verify the grant we queried is in the results
		found := false
		for _, g := range resp.Data {
			if g.ID == grantID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Grant %s not found in email filter results", grantID)
		}
	})
}

// Note: TestGrants_Update and TestGrants_Delete are intentionally skipped
// to avoid breaking test accounts. These operations are tested via unit tests.
