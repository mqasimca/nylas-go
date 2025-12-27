package connectors

import (
	"encoding/json"
	"testing"
)

func TestProvider_Constants(t *testing.T) {
	tests := []struct {
		provider Provider
		want     string
	}{
		{ProviderGoogle, "google"},
		{ProviderMicrosoft, "microsoft"},
		{ProviderIMAP, "imap"},
		{ProviderVirtualCalendars, "virtual-calendars"},
	}

	for _, tt := range tests {
		t.Run(string(tt.provider), func(t *testing.T) {
			if string(tt.provider) != tt.want {
				t.Errorf("Provider = %s, want %s", tt.provider, tt.want)
			}
		})
	}
}

func TestConnector_JSON(t *testing.T) {
	jsonData := `{
		"provider": "google",
		"name": "My Google Connector",
		"scope": ["email", "calendar"],
		"settings": {"client_id": "abc123"}
	}`

	var c Connector
	if err := json.Unmarshal([]byte(jsonData), &c); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if c.Provider != ProviderGoogle {
		t.Errorf("Provider = %s, want google", c.Provider)
	}
	if c.Name != "My Google Connector" {
		t.Errorf("Name = %s, want My Google Connector", c.Name)
	}
	if len(c.Scope) != 2 {
		t.Errorf("Scope length = %d, want 2", len(c.Scope))
	}
}

func TestCreateGoogleRequest(t *testing.T) {
	req := CreateGoogleRequest("Test", "client-id", "client-secret", []string{"email", "calendar"})

	if req.Name != "Test" {
		t.Errorf("Name = %s, want Test", req.Name)
	}
	if req.Provider != ProviderGoogle {
		t.Errorf("Provider = %s, want google", req.Provider)
	}
	if req.Settings["client_id"] != "client-id" {
		t.Errorf("client_id = %v, want client-id", req.Settings["client_id"])
	}
	if req.Settings["client_secret"] != "client-secret" {
		t.Errorf("client_secret = %v, want client-secret", req.Settings["client_secret"])
	}
	if len(req.Scope) != 2 {
		t.Errorf("Scope length = %d, want 2", len(req.Scope))
	}
}

func TestCreateMicrosoftRequest(t *testing.T) {
	tests := []struct {
		name       string
		tenant     string
		wantTenant bool
	}{
		{"with tenant", "my-tenant", true},
		{"without tenant", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := CreateMicrosoftRequest("Test", "client-id", "client-secret", tt.tenant, []string{"email"})

			if req.Provider != ProviderMicrosoft {
				t.Errorf("Provider = %s, want microsoft", req.Provider)
			}
			if req.Settings["client_id"] != "client-id" {
				t.Errorf("client_id = %v, want client-id", req.Settings["client_id"])
			}

			_, hasTenant := req.Settings["tenant"]
			if hasTenant != tt.wantTenant {
				t.Errorf("has tenant = %v, want %v", hasTenant, tt.wantTenant)
			}
		})
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
