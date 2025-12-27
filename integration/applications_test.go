//go:build integration

// Applications Integration Tests Coverage:
//   - GetDetails ✓
//
// This endpoint returns application-level configuration.
// No Create/Update/Delete operations exist for applications.

package integration

import (
	"testing"
)

func TestApplications_GetDetails(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	details, err := client.Applications.GetDetails(ctx)
	if err != nil {
		t.Fatalf("GetDetails() error = %v", err)
	}

	// Verify required fields
	if details.ApplicationID == "" {
		t.Error("ApplicationID should not be empty")
	}

	if details.OrganizationID == "" {
		t.Error("OrganizationID should not be empty")
	}

	if details.Region == "" {
		t.Error("Region should not be empty")
	}

	// Region must be "us" or "eu"
	if details.Region != "us" && details.Region != "eu" {
		t.Errorf("Region = %s, want 'us' or 'eu'", details.Region)
	}

	t.Logf("Application Details:")
	t.Logf("  ID: %s", details.ApplicationID)
	t.Logf("  Organization: %s", details.OrganizationID)
	t.Logf("  Region: %s", details.Region)
	t.Logf("  Environment: %s", details.Environment)

	if details.Branding != nil {
		t.Logf("  Branding Name: %s", details.Branding.Name)
	}
}
