package credentials

import (
	"encoding/json"
	"testing"
)

func TestCredentialType_Constants(t *testing.T) {
	tests := []struct {
		credType CredentialType
		want     string
	}{
		{CredentialTypeAdminConsent, "adminconsent"},
		{CredentialTypeServiceAccount, "serviceaccount"},
		{CredentialTypeConnector, "connector"},
	}

	for _, tt := range tests {
		t.Run(string(tt.credType), func(t *testing.T) {
			if string(tt.credType) != tt.want {
				t.Errorf("CredentialType = %s, want %s", tt.credType, tt.want)
			}
		})
	}
}

func TestCredential_JSON(t *testing.T) {
	jsonData := `{
		"id": "cred-123",
		"name": "My Credential",
		"credential_type": "adminconsent",
		"created_at": 1704067200,
		"updated_at": 1704153600
	}`

	var c Credential
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.ID != "cred-123" {
		t.Errorf("ID = %s, want cred-123", c.ID)
	}
	if c.Name != "My Credential" {
		t.Errorf("Name = %s, want My Credential", c.Name)
	}
	if c.CredentialType != CredentialTypeAdminConsent {
		t.Errorf("CredentialType = %s, want adminconsent", c.CredentialType)
	}
	if c.CreatedAt != 1704067200 {
		t.Errorf("CreatedAt = %d, want 1704067200", c.CreatedAt)
	}
}

func TestCreateMicrosoftRequest(t *testing.T) {
	req := CreateMicrosoftRequest("Test Cred", "client-id", "client-secret")

	if req.Name != "Test Cred" {
		t.Errorf("Name = %s, want Test Cred", req.Name)
	}
	if req.CredentialType != CredentialTypeAdminConsent {
		t.Errorf("CredentialType = %s, want adminconsent", req.CredentialType)
	}
	if req.CredentialData["client_id"] != "client-id" {
		t.Errorf("client_id = %v, want client-id", req.CredentialData["client_id"])
	}
	if req.CredentialData["client_secret"] != "client-secret" {
		t.Errorf("client_secret = %v, want client-secret", req.CredentialData["client_secret"])
	}
}

func TestCreateGoogleRequest(t *testing.T) {
	req := CreateGoogleRequest("Test Cred", "key-id", "private-key", "test@example.com")

	if req.Name != "Test Cred" {
		t.Errorf("Name = %s, want Test Cred", req.Name)
	}
	if req.CredentialType != CredentialTypeServiceAccount {
		t.Errorf("CredentialType = %s, want serviceaccount", req.CredentialType)
	}
	if req.CredentialData["private_key_id"] != "key-id" {
		t.Errorf("private_key_id = %v, want key-id", req.CredentialData["private_key_id"])
	}
	if req.CredentialData["private_key"] != "private-key" {
		t.Errorf("private_key = %v, want private-key", req.CredentialData["private_key"])
	}
	if req.CredentialData["client_email"] != "test@example.com" {
		t.Errorf("client_email = %v, want test@example.com", req.CredentialData["client_email"])
	}
}

func TestListOptions_Values(t *testing.T) {
	tests := []struct {
		name      string
		opts      *ListOptions
		wantNil   bool
		wantLimit bool
		wantToken bool
	}{
		{"nil options", nil, true, false, false},
		{"empty options", &ListOptions{}, false, false, false},
		{"with limit", &ListOptions{Limit: intPtr(10)}, false, true, false},
		{"with page token", &ListOptions{PageToken: "abc"}, false, false, true},
		{"with both", &ListOptions{Limit: intPtr(20), PageToken: "xyz"}, false, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals := tt.opts.Values()

			if tt.wantNil {
				if vals != nil {
					t.Errorf("Values() = %v, want nil", vals)
				}
				return
			}

			_, hasLimit := vals["limit"]
			if hasLimit != tt.wantLimit {
				t.Errorf("has limit = %v, want %v", hasLimit, tt.wantLimit)
			}

			_, hasToken := vals["page_token"]
			if hasToken != tt.wantToken {
				t.Errorf("has page_token = %v, want %v", hasToken, tt.wantToken)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
