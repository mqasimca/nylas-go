package auth

import (
	"strings"
	"testing"
)

func TestURLForAuthenticationConfig_Values(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *URLForAuthenticationConfig
		if c.Values() != nil {
			t.Error("expected nil for nil config")
		}
	})

	t.Run("empty config", func(t *testing.T) {
		c := &URLForAuthenticationConfig{}
		v := c.Values()
		if len(v) != 0 {
			t.Errorf("expected empty map, got %d entries", len(v))
		}
	})

	t.Run("with all options", func(t *testing.T) {
		c := &URLForAuthenticationConfig{
			ClientID:           "client123",
			RedirectURI:        "https://example.com/callback",
			Provider:           "google",
			LoginHint:          "user@example.com",
			State:              "state123",
			Scopes:             []string{"email", "calendar"},
			AccessType:         "offline",
			ResponseType:       "code",
			Prompt:             "select_provider",
			IncludeGrantScopes: true,
			CredentialID:       "cred123",
		}
		v := c.Values()

		if v["client_id"] != "client123" {
			t.Errorf("expected client_id=client123, got %v", v["client_id"])
		}
		if v["redirect_uri"] != "https://example.com/callback" {
			t.Errorf("expected redirect_uri, got %v", v["redirect_uri"])
		}
		if v["provider"] != "google" {
			t.Errorf("expected provider=google, got %v", v["provider"])
		}
		if v["login_hint"] != "user@example.com" {
			t.Errorf("expected login_hint, got %v", v["login_hint"])
		}
		if v["state"] != "state123" {
			t.Errorf("expected state=state123, got %v", v["state"])
		}
		if v["access_type"] != "offline" {
			t.Errorf("expected access_type=offline, got %v", v["access_type"])
		}
		if v["response_type"] != "code" {
			t.Errorf("expected response_type=code, got %v", v["response_type"])
		}
		if v["prompt"] != "select_provider" {
			t.Errorf("expected prompt=select_provider, got %v", v["prompt"])
		}
		if v["include_grant_scopes"] != true {
			t.Errorf("expected include_grant_scopes=true, got %v", v["include_grant_scopes"])
		}
		if v["credential_id"] != "cred123" {
			t.Errorf("expected credential_id=cred123, got %v", v["credential_id"])
		}

		// Check scopes
		scopes, ok := v["scope"].([]string)
		if !ok {
			t.Errorf("expected scope to be []string, got %T", v["scope"])
		} else if len(scopes) != 2 {
			t.Errorf("expected 2 scopes, got %d", len(scopes))
		}
	})

	t.Run("partial options", func(t *testing.T) {
		c := &URLForAuthenticationConfig{
			ClientID:    "client123",
			RedirectURI: "https://example.com/callback",
		}
		v := c.Values()

		if len(v) != 2 {
			t.Errorf("expected 2 entries, got %d", len(v))
		}
	})
}

func TestPKCEURLConfig_Values(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *PKCEURLConfig
		if c.Values() != nil {
			t.Error("expected nil for nil config")
		}
	})

	t.Run("with PKCE params", func(t *testing.T) {
		c := &PKCEURLConfig{
			URLForAuthenticationConfig: URLForAuthenticationConfig{
				ClientID:    "client123",
				RedirectURI: "https://example.com/callback",
			},
			CodeChallenge:       "challenge_abc123",
			CodeChallengeMethod: "S256",
		}
		v := c.Values()

		if v["client_id"] != "client123" {
			t.Errorf("expected client_id=client123, got %v", v["client_id"])
		}
		if v["code_challenge"] != "challenge_abc123" {
			t.Errorf("expected code_challenge, got %v", v["code_challenge"])
		}
		if v["code_challenge_method"] != "S256" {
			t.Errorf("expected code_challenge_method=S256, got %v", v["code_challenge_method"])
		}
	})

	t.Run("empty base config", func(t *testing.T) {
		c := &PKCEURLConfig{
			CodeChallenge:       "challenge",
			CodeChallengeMethod: "S256",
		}
		v := c.Values()

		if v["code_challenge"] != "challenge" {
			t.Errorf("expected code_challenge, got %v", v["code_challenge"])
		}
	})
}

func TestAdminConsentURLConfig_Values(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		var c *AdminConsentURLConfig
		if c.Values() != nil {
			t.Error("expected nil for nil config")
		}
	})

	t.Run("with all options", func(t *testing.T) {
		c := &AdminConsentURLConfig{
			ClientID:     "client123",
			RedirectURI:  "https://example.com/callback",
			State:        "state123",
			CredentialID: "cred123",
		}
		v := c.Values()

		if v["client_id"] != "client123" {
			t.Errorf("expected client_id=client123, got %v", v["client_id"])
		}
		if v["redirect_uri"] != "https://example.com/callback" {
			t.Errorf("expected redirect_uri, got %v", v["redirect_uri"])
		}
		if v["state"] != "state123" {
			t.Errorf("expected state=state123, got %v", v["state"])
		}
		if v["credential_id"] != "cred123" {
			t.Errorf("expected credential_id=cred123, got %v", v["credential_id"])
		}
		if v["response_type"] != "adminconsent" {
			t.Errorf("expected response_type=adminconsent, got %v", v["response_type"])
		}
	})

	t.Run("empty config sets response_type", func(t *testing.T) {
		c := &AdminConsentURLConfig{}
		v := c.Values()

		if v["response_type"] != "adminconsent" {
			t.Errorf("expected response_type=adminconsent, got %v", v["response_type"])
		}
	})
}

func TestURLBuilder(t *testing.T) {
	t.Run("empty params", func(t *testing.T) {
		b := NewURLBuilder("https://api.nylas.com/v3/connect/auth")
		url := b.Build()

		if url != "https://api.nylas.com/v3/connect/auth" {
			t.Errorf("unexpected URL: %s", url)
		}
	})

	t.Run("with params", func(t *testing.T) {
		b := NewURLBuilder("https://api.nylas.com/v3/connect/auth")
		b.Add("client_id", "client123")
		b.Add("redirect_uri", "https://example.com/callback")
		url := b.Build()

		if !strings.HasPrefix(url, "https://api.nylas.com/v3/connect/auth?") {
			t.Errorf("unexpected URL prefix: %s", url)
		}
		if !strings.Contains(url, "client_id=client123") {
			t.Errorf("expected client_id in URL: %s", url)
		}
		if !strings.Contains(url, "redirect_uri=") {
			t.Errorf("expected redirect_uri in URL: %s", url)
		}
	})

	t.Run("Add skips empty values", func(t *testing.T) {
		b := NewURLBuilder("https://example.com")
		b.Add("key1", "value1")
		b.Add("key2", "")
		b.Add("key3", "value3")
		url := b.Build()

		if strings.Contains(url, "key2") {
			t.Errorf("expected key2 to be skipped: %s", url)
		}
		if !strings.Contains(url, "key1=value1") {
			t.Errorf("expected key1 in URL: %s", url)
		}
		if !strings.Contains(url, "key3=value3") {
			t.Errorf("expected key3 in URL: %s", url)
		}
	})

	t.Run("AddBool", func(t *testing.T) {
		b := NewURLBuilder("https://example.com")
		b.AddBool("include", true)
		b.AddBool("exclude", false)
		url := b.Build()

		if !strings.Contains(url, "include=true") {
			t.Errorf("expected include=true in URL: %s", url)
		}
		if strings.Contains(url, "exclude") {
			t.Errorf("expected exclude to be skipped: %s", url)
		}
	})

	t.Run("AddSlice", func(t *testing.T) {
		b := NewURLBuilder("https://example.com")
		b.AddSlice("scope", []string{"email", "calendar"})
		url := b.Build()

		if !strings.Contains(url, "scope=email") {
			t.Errorf("expected scope=email in URL: %s", url)
		}
		if !strings.Contains(url, "scope=calendar") {
			t.Errorf("expected scope=calendar in URL: %s", url)
		}
	})

	t.Run("chaining", func(t *testing.T) {
		url := NewURLBuilder("https://example.com").
			Add("a", "1").
			Add("b", "2").
			AddBool("c", true).
			Build()

		if !strings.Contains(url, "a=1") || !strings.Contains(url, "b=2") || !strings.Contains(url, "c=true") {
			t.Errorf("chaining failed: %s", url)
		}
	})
}
