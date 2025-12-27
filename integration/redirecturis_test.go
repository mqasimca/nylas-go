//go:build integration

// RedirectURIs Integration Tests Coverage:
//   - List, Get ✓
//
// Intentionally NOT tested (safety reasons):
//   - Create: Would add new OAuth redirect URIs
//   - Update: Could break existing OAuth flows
//   - Delete: Would remove redirect URIs and break OAuth authentication
//
// These operations are tested via unit tests with mocks.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/redirecturis"
)

func TestRedirectURIs_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	resp, err := client.RedirectURIs.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	t.Logf("Found %d redirect URIs", len(resp.Data))

	for _, uri := range resp.Data {
		if uri.ID == "" {
			t.Error("RedirectURI ID should not be empty")
		}
		if uri.URL == "" {
			t.Error("RedirectURI URL should not be empty")
		}
		t.Logf("  - ID: %s, URL: %s, Platform: %s",
			uri.ID, uri.URL, uri.Platform)
	}
}

func TestRedirectURIs_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	limit := 5
	resp, err := client.RedirectURIs.List(ctx, &redirecturis.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(resp.Data) > limit {
		t.Errorf("List() returned %d redirect URIs, want <= %d", len(resp.Data), limit)
	}

	t.Logf("Found %d redirect URIs (limit %d)", len(resp.Data), limit)
}

func TestRedirectURIs_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list redirect URIs to find one
	listResp, err := client.RedirectURIs.List(ctx, nil)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(listResp.Data) == 0 {
		t.Skip("No redirect URIs found")
	}

	// Get each redirect URI by ID
	for _, uri := range listResp.Data {
		t.Run(uri.ID, func(t *testing.T) {
			got, err := client.RedirectURIs.Get(ctx, uri.ID)
			if err != nil {
				t.Fatalf("Get(%s) error = %v", uri.ID, err)
			}

			if got.ID != uri.ID {
				t.Errorf("Get() ID = %s, want %s", got.ID, uri.ID)
			}

			if got.URL != uri.URL {
				t.Errorf("Get() URL = %s, want %s", got.URL, uri.URL)
			}

			t.Logf("Got redirect URI: %s", got.URL)
			t.Logf("  Platform: %s", got.Platform)

			if got.Settings != nil {
				if got.Settings.Origin != "" {
					t.Logf("  Origin: %s", got.Settings.Origin)
				}
				if got.Settings.BundleID != "" {
					t.Logf("  Bundle ID: %s", got.Settings.BundleID)
				}
				if got.Settings.PackageName != "" {
					t.Logf("  Package Name: %s", got.Settings.PackageName)
				}
			}
		})
	}
}

func TestRedirectURIs_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	iter := client.RedirectURIs.ListAll(ctx, nil)
	all, err := iter.Collect()
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}

	t.Logf("ListAll() returned %d redirect URIs", len(all))

	for _, uri := range all {
		if uri.ID == "" {
			t.Error("RedirectURI ID should not be empty")
		}
		if uri.URL == "" {
			t.Error("RedirectURI URL should not be empty")
		}
	}
}

// Note: TestRedirectURIs_Create, TestRedirectURIs_Update, TestRedirectURIs_Delete
// are intentionally skipped to avoid breaking OAuth configuration.
// These operations are tested via unit tests.
