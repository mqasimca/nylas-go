package redirecturis

import (
	"encoding/json"
	"testing"
)

func TestRedirectURI_JSON(t *testing.T) {
	jsonData := `{
		"id": "uri-123",
		"url": "https://example.com/callback",
		"platform": "web",
		"settings": {
			"origin": "https://example.com"
		}
	}`

	var uri RedirectURI
	if err := json.Unmarshal([]byte(jsonData), &uri); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if uri.ID != "uri-123" {
		t.Errorf("ID = %s, want uri-123", uri.ID)
	}
	if uri.URL != "https://example.com/callback" {
		t.Errorf("URL = %s, want https://example.com/callback", uri.URL)
	}
	if uri.Platform != "web" {
		t.Errorf("Platform = %s, want web", uri.Platform)
	}
	if uri.Settings == nil {
		t.Error("Settings is nil, want non-nil")
	}
	if uri.Settings.Origin != "https://example.com" {
		t.Errorf("Settings.Origin = %s, want https://example.com", uri.Settings.Origin)
	}
}

func TestRedirectURI_iOSSettings(t *testing.T) {
	jsonData := `{
		"id": "uri-ios",
		"url": "myapp://callback",
		"platform": "ios",
		"settings": {
			"bundle_id": "com.example.app",
			"app_store_id": "123456789",
			"team_id": "TEAM123"
		}
	}`

	var uri RedirectURI
	if err := json.Unmarshal([]byte(jsonData), &uri); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if uri.Platform != "ios" {
		t.Errorf("Platform = %s, want ios", uri.Platform)
	}
	if uri.Settings.BundleID != "com.example.app" {
		t.Errorf("BundleID = %s, want com.example.app", uri.Settings.BundleID)
	}
	if uri.Settings.AppStoreID != "123456789" {
		t.Errorf("AppStoreID = %s, want 123456789", uri.Settings.AppStoreID)
	}
	if uri.Settings.TeamID != "TEAM123" {
		t.Errorf("TeamID = %s, want TEAM123", uri.Settings.TeamID)
	}
}

func TestRedirectURI_AndroidSettings(t *testing.T) {
	jsonData := `{
		"id": "uri-android",
		"url": "myapp://callback",
		"platform": "android",
		"settings": {
			"package_name": "com.example.app",
			"sha1_certificate_fingerprint": "AA:BB:CC:DD:EE:FF"
		}
	}`

	var uri RedirectURI
	if err := json.Unmarshal([]byte(jsonData), &uri); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if uri.Platform != "android" {
		t.Errorf("Platform = %s, want android", uri.Platform)
	}
	if uri.Settings.PackageName != "com.example.app" {
		t.Errorf("PackageName = %s, want com.example.app", uri.Settings.PackageName)
	}
	if uri.Settings.SHA1CertificateFingerprint != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("SHA1CertificateFingerprint = %s, want AA:BB:CC:DD:EE:FF", uri.Settings.SHA1CertificateFingerprint)
	}
}

func TestCreateRequest_JSON(t *testing.T) {
	req := CreateRequest{
		URL:      "https://example.com/callback",
		Platform: "web",
		Settings: &RedirectSettings{
			Origin: "https://example.com",
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var decoded CreateRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if decoded.URL != req.URL {
		t.Errorf("URL = %s, want %s", decoded.URL, req.URL)
	}
	if decoded.Platform != req.Platform {
		t.Errorf("Platform = %s, want %s", decoded.Platform, req.Platform)
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
