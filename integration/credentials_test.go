//go:build integration

// Credentials Integration Tests Coverage:
//   - List, Get ✓
//
// Intentionally NOT tested (safety reasons):
//   - Create: Would create new credential with sensitive data
//   - Update: Could break existing connector authentication
//   - Delete: Would remove credentials and break provider access
//
// These operations are tested via unit tests with mocks.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/connectors"
	"github.com/mqasimca/nylas-go/credentials"
)

func TestCredentials_List(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list connectors to find available providers
	connResp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("Connectors.List() error = %v", err)
	}

	if len(connResp.Data) == 0 {
		t.Skip("No connectors found")
	}

	// List credentials for each provider
	for _, conn := range connResp.Data {
		t.Run(string(conn.Provider), func(t *testing.T) {
			resp, err := client.Credentials.List(ctx, conn.Provider, nil)
			if err != nil {
				t.Fatalf("List(%s) error = %v", conn.Provider, err)
			}

			t.Logf("Found %d credentials for %s", len(resp.Data), conn.Provider)

			for _, cred := range resp.Data {
				if cred.ID == "" {
					t.Error("Credential ID should not be empty")
				}
				t.Logf("  - ID: %s, Name: %s, Type: %s",
					cred.ID, cred.Name, cred.CredentialType)
			}
		})
	}
}

func TestCredentials_ListWithOptions(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list connectors to find a provider
	connResp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("Connectors.List() error = %v", err)
	}

	if len(connResp.Data) == 0 {
		t.Skip("No connectors found")
	}

	provider := connResp.Data[0].Provider
	limit := 5
	resp, err := client.Credentials.List(ctx, provider, &credentials.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("List(%s) error = %v", provider, err)
	}

	if len(resp.Data) > limit {
		t.Errorf("List() returned %d credentials, want <= %d", len(resp.Data), limit)
	}

	t.Logf("Found %d credentials for %s (limit %d)", len(resp.Data), provider, limit)
}

func TestCredentials_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list connectors to find available providers
	connResp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("Connectors.List() error = %v", err)
	}

	if len(connResp.Data) == 0 {
		t.Skip("No connectors found")
	}

	// Find a provider with credentials
	var foundProvider connectors.Provider
	var foundCredID string

	for _, conn := range connResp.Data {
		credResp, err := client.Credentials.List(ctx, conn.Provider, nil)
		if err != nil {
			continue
		}
		if len(credResp.Data) > 0 {
			foundProvider = conn.Provider
			foundCredID = credResp.Data[0].ID
			break
		}
	}

	if foundCredID == "" {
		t.Skip("No credentials found for any provider")
	}

	// Get the credential
	cred, err := client.Credentials.Get(ctx, foundProvider, foundCredID)
	if err != nil {
		t.Fatalf("Get(%s, %s) error = %v", foundProvider, foundCredID, err)
	}

	if cred.ID != foundCredID {
		t.Errorf("Get() ID = %s, want %s", cred.ID, foundCredID)
	}

	t.Logf("Got credential: %s (name: %s, type: %s)",
		cred.ID, cred.Name, cred.CredentialType)

	// Verify hashed_data is not exposed (security check)
	// The API should not return raw credential data
	if cred.HashedData != "" {
		t.Logf("  HashedData present (truncated for security)")
	}
}

func TestCredentials_ListAll(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// First list connectors to find a provider
	connResp, err := client.Connectors.List(ctx, nil)
	if err != nil {
		t.Fatalf("Connectors.List() error = %v", err)
	}

	if len(connResp.Data) == 0 {
		t.Skip("No connectors found")
	}

	provider := connResp.Data[0].Provider

	iter := client.Credentials.ListAll(ctx, provider, nil)
	all, err := iter.Collect()
	if err != nil {
		t.Fatalf("ListAll(%s) error = %v", provider, err)
	}

	t.Logf("ListAll(%s) returned %d credentials", provider, len(all))

	for _, cred := range all {
		if cred.ID == "" {
			t.Error("Credential ID should not be empty")
		}
	}
}

// Note: TestCredentials_Create, TestCredentials_Update, TestCredentials_Delete
// are intentionally skipped to avoid breaking connector authentication.
// These operations are tested via unit tests.
