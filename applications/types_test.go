package applications

import (
	"encoding/json"
	"testing"
)

func TestApplicationDetails_JSON(t *testing.T) {
	jsonData := `{
		"application_id": "app-123",
		"organization_id": "org-456",
		"region": "us",
		"environment": "production",
		"branding": {
			"name": "My App",
			"icon_url": "https://example.com/icon.png"
		}
	}`

	var details ApplicationDetails
	if err := json.Unmarshal([]byte(jsonData), &details); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if details.ApplicationID != "app-123" {
		t.Errorf("ApplicationID = %s, want app-123", details.ApplicationID)
	}
	if details.Region != "us" {
		t.Errorf("Region = %s, want us", details.Region)
	}
	if details.Branding == nil {
		t.Fatal("Branding should not be nil")
	}
	if details.Branding.Name != "My App" {
		t.Errorf("Branding.Name = %s, want My App", details.Branding.Name)
	}
}
