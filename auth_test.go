package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/mqasimca/nylas-go/auth"
)

func TestAuthService_URLForOAuth2(t *testing.T) {
	client, _ := NewClient(WithAPIKey("test-key"))

	tests := []struct {
		name   string
		config *auth.URLForAuthenticationConfig
		want   map[string]string
	}{
		{
			name: "basic config",
			config: &auth.URLForAuthenticationConfig{
				ClientID:    "client123",
				RedirectURI: "https://example.com/callback",
			},
			want: map[string]string{
				"client_id":     "client123",
				"redirect_uri":  "https://example.com/callback",
				"response_type": "code",
			},
		},
		{
			name: "with all options",
			config: &auth.URLForAuthenticationConfig{
				ClientID:           "client123",
				RedirectURI:        "https://example.com/callback",
				Provider:           "google",
				LoginHint:          "user@example.com",
				State:              "state123",
				Scopes:             []string{"email", "calendar"},
				AccessType:         "offline",
				Prompt:             "select_provider",
				IncludeGrantScopes: true,
			},
			want: map[string]string{
				"client_id":            "client123",
				"redirect_uri":         "https://example.com/callback",
				"response_type":        "code",
				"provider":             "google",
				"login_hint":           "user@example.com",
				"state":                "state123",
				"access_type":          "offline",
				"prompt":               "select_provider",
				"include_grant_scopes": "true",
			},
		},
		{
			name:   "nil config",
			config: nil,
			want:   map[string]string{"response_type": "code"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlStr := client.Auth.URLForOAuth2(tt.config)

			u, err := url.Parse(urlStr)
			if err != nil {
				t.Fatalf("failed to parse URL: %v", err)
			}

			if !strings.Contains(u.Path, "/v3/connect/auth") {
				t.Errorf("unexpected path: %s", u.Path)
			}

			for key, wantValue := range tt.want {
				gotValue := u.Query().Get(key)
				if gotValue != wantValue {
					t.Errorf("param %s: got %q, want %q", key, gotValue, wantValue)
				}
			}
		})
	}
}

func TestAuthService_URLForOAuth2PKCE(t *testing.T) {
	client, _ := NewClient(WithAPIKey("test-key"))

	config := &auth.PKCEURLConfig{
		URLForAuthenticationConfig: auth.URLForAuthenticationConfig{
			ClientID:    "client123",
			RedirectURI: "https://example.com/callback",
		},
		CodeChallenge:       "challenge123",
		CodeChallengeMethod: "S256",
	}

	urlStr := client.Auth.URLForOAuth2PKCE(config)

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	want := map[string]string{
		"client_id":             "client123",
		"redirect_uri":          "https://example.com/callback",
		"response_type":         "code",
		"code_challenge":        "challenge123",
		"code_challenge_method": "S256",
	}

	for key, wantValue := range want {
		gotValue := u.Query().Get(key)
		if gotValue != wantValue {
			t.Errorf("param %s: got %q, want %q", key, gotValue, wantValue)
		}
	}
}

func TestAuthService_URLForAdminConsent(t *testing.T) {
	client, _ := NewClient(WithAPIKey("test-key"))

	config := &auth.AdminConsentURLConfig{
		ClientID:     "client123",
		RedirectURI:  "https://example.com/callback",
		State:        "state123",
		CredentialID: "cred123",
	}

	urlStr := client.Auth.URLForAdminConsent(config)

	u, err := url.Parse(urlStr)
	if err != nil {
		t.Fatalf("failed to parse URL: %v", err)
	}

	want := map[string]string{
		"client_id":     "client123",
		"redirect_uri":  "https://example.com/callback",
		"state":         "state123",
		"credential_id": "cred123",
		"response_type": "adminconsent",
	}

	for key, wantValue := range want {
		gotValue := u.Query().Get(key)
		if gotValue != wantValue {
			t.Errorf("param %s: got %q, want %q", key, gotValue, wantValue)
		}
	}
}

func TestAuthService_ExchangeCodeForToken(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			response: `{
				"access_token": "access123",
				"refresh_token": "refresh123",
				"grant_id": "grant123",
				"email": "user@example.com",
				"token_type": "Bearer",
				"expires_in": 3600
			}`,
			wantErr: false,
		},
		{
			name:       "invalid code",
			statusCode: http.StatusBadRequest,
			response:   `{"error": "invalid_grant", "error_description": "Invalid code"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/v3/connect/token") {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				var req auth.CodeExchangeRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Errorf("failed to decode request: %v", err)
				}

				// Verify grant_type is auto-set
				if req.GrantType != "authorization_code" {
					t.Errorf("GrantType = %s, want authorization_code (should be auto-set)", req.GrantType)
				}

				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(
				WithAPIKey("test-key"),
				WithBaseURL(srv.URL),
			)

			resp, err := client.Auth.ExchangeCodeForToken(context.Background(), &auth.CodeExchangeRequest{
				ClientID:    "client123",
				Code:        "code123",
				RedirectURI: "https://example.com/callback",
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp.AccessToken != "access123" {
					t.Errorf("AccessToken = %s, want access123", resp.AccessToken)
				}
				if resp.GrantID != "grant123" {
					t.Errorf("GrantID = %s, want grant123", resp.GrantID)
				}
			}
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}

		var req auth.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		// Verify grant_type is auto-set
		if req.GrantType != "refresh_token" {
			t.Errorf("GrantType = %s, want refresh_token (should be auto-set)", req.GrantType)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"access_token": "new_access123",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
	)

	resp, err := client.Auth.RefreshAccessToken(context.Background(), &auth.RefreshTokenRequest{
		ClientID:     "client123",
		RefreshToken: "refresh123",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.AccessToken != "new_access123" {
		t.Errorf("AccessToken = %s, want new_access123", resp.AccessToken)
	}
}

func TestAuthService_CustomAuthentication(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v3/connect/custom") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req123",
			"data": {
				"id": "grant123",
				"provider": "google",
				"email": "user@example.com",
				"grant_status": "valid"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
	)

	grant, err := client.Auth.CustomAuthentication(context.Background(), &auth.CustomAuthRequest{
		Provider: "google",
		Settings: map[string]any{
			"refresh_token": "provider_refresh_token",
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if grant.ID != "grant123" {
		t.Errorf("ID = %s, want grant123", grant.ID)
	}
	if grant.Provider != "google" {
		t.Errorf("Provider = %s, want google", grant.Provider)
	}
}

func TestAuthService_IDTokenInfo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v3/connect/tokeninfo") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		idToken := r.URL.Query().Get("id_token")
		if idToken != "test_id_token" {
			t.Errorf("id_token = %s, want test_id_token", idToken)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req123",
			"data": {
				"iss": "https://nylas.com",
				"sub": "grant123",
				"email": "user@example.com",
				"email_verified": true
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
	)

	info, err := client.Auth.IDTokenInfo(context.Background(), "test_id_token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Sub != "grant123" {
		t.Errorf("Sub = %s, want grant123", info.Sub)
	}
	if info.Email != "user@example.com" {
		t.Errorf("Email = %s, want user@example.com", info.Email)
	}
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		accessToken := r.URL.Query().Get("access_token")
		if accessToken != "test_access_token" {
			t.Errorf("access_token = %s, want test_access_token", accessToken)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req123",
			"data": {
				"iss": "https://nylas.com",
				"sub": "grant123",
				"email": "user@example.com"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
	)

	info, err := client.Auth.ValidateAccessToken(context.Background(), "test_access_token")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Sub != "grant123" {
		t.Errorf("Sub = %s, want grant123", info.Sub)
	}
}

func TestAuthService_Revoke(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "invalid token",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/v3/connect/revoke") {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				token := r.URL.Query().Get("token")
				if token != "test_token" {
					t.Errorf("token = %s, want test_token", token)
				}

				w.WriteHeader(tt.statusCode)
				if tt.statusCode != http.StatusOK {
					_, _ = w.Write([]byte(`{"error": "invalid_token"}`))
				}
			}))
			defer srv.Close()

			client, _ := NewClient(
				WithAPIKey("test-key"),
				WithBaseURL(srv.URL),
			)

			err := client.Auth.Revoke(context.Background(), "test_token")

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_DetectProvider(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/v3/providers/detect") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		var req auth.ProviderDetectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Email != "user@gmail.com" {
			t.Errorf("Email = %s, want user@gmail.com", req.Email)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"request_id": "req123",
			"data": {
				"email_address": "user@gmail.com",
				"detected": true,
				"provider": "google",
				"type": "oauth"
			}
		}`))
	}))
	defer srv.Close()

	client, _ := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(srv.URL),
	)

	resp, err := client.Auth.DetectProvider(context.Background(), &auth.ProviderDetectRequest{
		Email: "user@gmail.com",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Provider != "google" {
		t.Errorf("Provider = %s, want google", resp.Provider)
	}
	if !resp.Detected {
		t.Errorf("Detected = false, want true")
	}
}

func TestAuthService_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error": "unauthorized", "message": "Invalid API key"}`,
		},
		{
			name:       "server error",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal_error", "message": "Server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer srv.Close()

			client, _ := NewClient(
				WithAPIKey("test-key"),
				WithBaseURL(srv.URL),
				WithMaxRetries(0),
			)

			_, err := client.Auth.ExchangeCodeForToken(context.Background(), &auth.CodeExchangeRequest{
				Code: "code123",
			})

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
