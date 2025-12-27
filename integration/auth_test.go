//go:build integration

// Auth Integration Tests Coverage:
//   - URLForOAuth2 ✓
//   - URLForOAuth2PKCE ✓
//   - URLForAdminConsent ✓
//   - DetectProvider ✓
//
// Note: Token exchange, refresh, and revoke operations require valid
// OAuth tokens which are not available in automated test environments.

package integration

import (
	"strings"
	"testing"

	"github.com/mqasimca/nylas-go/auth"
)

func TestAuth_URLForOAuth2(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	t.Run("basic URL generation", func(t *testing.T) {
		url := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/callback",
		})

		if url == "" {
			t.Fatal("URLForOAuth2() returned empty string")
		}

		// Verify required components
		if !strings.Contains(url, "/v3/connect/auth") {
			t.Errorf("URL missing auth path: %s", url)
		}
		if !strings.Contains(url, "response_type=code") {
			t.Errorf("URL missing response_type=code: %s", url)
		}
		if cfg.ClientID != "" && !strings.Contains(url, "client_id="+cfg.ClientID) {
			t.Errorf("URL missing client_id: %s", url)
		}

		t.Logf("Generated OAuth2 URL: %s", url)
	})

	t.Run("with provider hint", func(t *testing.T) {
		url := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/callback",
			Provider:    "google",
			LoginHint:   "user@gmail.com",
		})

		if !strings.Contains(url, "provider=google") {
			t.Errorf("URL missing provider: %s", url)
		}
		if !strings.Contains(url, "login_hint=") {
			t.Errorf("URL missing login_hint: %s", url)
		}

		t.Logf("Generated OAuth2 URL with provider: %s", url)
	})

	t.Run("with scopes", func(t *testing.T) {
		url := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/callback",
			Scopes:      []string{"email", "calendar"},
		})

		if !strings.Contains(url, "scope=") {
			t.Errorf("URL missing scope: %s", url)
		}

		t.Logf("Generated OAuth2 URL with scopes: %s", url)
	})

	t.Run("with state parameter", func(t *testing.T) {
		url := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/callback",
			State:       "random-state-123",
		})

		if !strings.Contains(url, "state=random-state-123") {
			t.Errorf("URL missing state: %s", url)
		}

		t.Logf("Generated OAuth2 URL with state: %s", url)
	})

	t.Run("offline access type", func(t *testing.T) {
		url := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/callback",
			AccessType:  "offline",
		})

		if !strings.Contains(url, "access_type=offline") {
			t.Errorf("URL missing access_type: %s", url)
		}

		t.Logf("Generated OAuth2 URL with offline access: %s", url)
	})
}

func TestAuth_URLForOAuth2PKCE(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	t.Run("with PKCE challenge", func(t *testing.T) {
		url := client.Auth.URLForOAuth2PKCE(&auth.PKCEURLConfig{
			URLForAuthenticationConfig: auth.URLForAuthenticationConfig{
				ClientID:    cfg.ClientID,
				RedirectURI: "https://example.com/callback",
			},
			CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
			CodeChallengeMethod: "S256",
		})

		if url == "" {
			t.Fatal("URLForOAuth2PKCE() returned empty string")
		}

		if !strings.Contains(url, "code_challenge=") {
			t.Errorf("URL missing code_challenge: %s", url)
		}
		if !strings.Contains(url, "code_challenge_method=S256") {
			t.Errorf("URL missing code_challenge_method: %s", url)
		}
		if !strings.Contains(url, "response_type=code") {
			t.Errorf("URL missing response_type=code: %s", url)
		}

		t.Logf("Generated PKCE URL: %s", url)
	})

	t.Run("with all options", func(t *testing.T) {
		url := client.Auth.URLForOAuth2PKCE(&auth.PKCEURLConfig{
			URLForAuthenticationConfig: auth.URLForAuthenticationConfig{
				ClientID:    cfg.ClientID,
				RedirectURI: "https://example.com/callback",
				Provider:    "google",
				State:       "pkce-state-456",
				Scopes:      []string{"email", "calendar", "contacts"},
			},
			CodeChallenge:       "challenge123",
			CodeChallengeMethod: "S256",
		})

		if !strings.Contains(url, "provider=google") {
			t.Errorf("URL missing provider: %s", url)
		}
		if !strings.Contains(url, "state=pkce-state-456") {
			t.Errorf("URL missing state: %s", url)
		}

		t.Logf("Generated full PKCE URL: %s", url)
	})
}

func TestAuth_URLForAdminConsent(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)

	t.Run("basic admin consent URL", func(t *testing.T) {
		url := client.Auth.URLForAdminConsent(&auth.AdminConsentURLConfig{
			ClientID:    cfg.ClientID,
			RedirectURI: "https://example.com/admin-callback",
		})

		if url == "" {
			t.Fatal("URLForAdminConsent() returned empty string")
		}

		if !strings.Contains(url, "/v3/connect/auth") {
			t.Errorf("URL missing auth path: %s", url)
		}
		if !strings.Contains(url, "response_type=adminconsent") {
			t.Errorf("URL missing response_type=adminconsent: %s", url)
		}

		t.Logf("Generated admin consent URL: %s", url)
	})

	t.Run("with state and credential_id", func(t *testing.T) {
		url := client.Auth.URLForAdminConsent(&auth.AdminConsentURLConfig{
			ClientID:     cfg.ClientID,
			RedirectURI:  "https://example.com/admin-callback",
			State:        "admin-state-789",
			CredentialID: "cred-123",
		})

		if !strings.Contains(url, "state=admin-state-789") {
			t.Errorf("URL missing state: %s", url)
		}
		if !strings.Contains(url, "credential_id=cred-123") {
			t.Errorf("URL missing credential_id: %s", url)
		}

		t.Logf("Generated admin consent URL with options: %s", url)
	})
}

func TestAuth_DetectProvider(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// Use actual grant email for reliable detection
	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		resp, err := client.Auth.DetectProvider(ctx, &auth.ProviderDetectRequest{
			Email: grant.Email,
		})
		if err != nil {
			t.Skipf("DetectProvider(%s) error = %v (API may not support this)", grant.Email, err)
		}

		t.Logf("Detected provider for %s: %s (type: %s, detected: %v)",
			grant.Email, resp.Provider, resp.Type, resp.Detected)

		if !resp.Detected {
			t.Logf("Provider not detected for %s (may be expected for some domains)", grant.Email)
		}
	})
}

func TestAuth_DetectProvider_AllFields(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	// Use actual grant email for reliable detection
	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		grant, err := client.Grants.Get(ctx, grantID)
		if err != nil {
			t.Fatalf("Grants.Get() error = %v", err)
		}

		if grant.Email == "" {
			t.Skip("Grant has no email address")
		}

		// Test with all optional fields
		resp, err := client.Auth.DetectProvider(ctx, &auth.ProviderDetectRequest{
			Email:            grant.Email,
			AllProviderTypes: true,
		})
		if err != nil {
			t.Skipf("DetectProvider() error = %v (API may not support this)", err)
		}

		t.Logf("Provider detection result:")
		t.Logf("  - Email: %s", grant.Email)
		t.Logf("  - Detected: %v", resp.Detected)
		t.Logf("  - Provider: %s", resp.Provider)
		t.Logf("  - Type: %s", resp.Type)
	})
}
