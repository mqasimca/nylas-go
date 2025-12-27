package nylas

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mqasimca/nylas-go/auth"
	"github.com/mqasimca/nylas-go/grants"
)

// URLForOAuth2 generates an OAuth 2.0 authorization URL.
// Users should be redirected to this URL to authenticate with their email provider.
//
// Example - basic OAuth flow:
//
//	authURL := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
//	    ClientID:    "your-client-id",
//	    RedirectURI: "https://yourapp.com/callback",
//	})
//	// Redirect user to authURL
//
// Example - with provider hint and scopes:
//
//	authURL := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
//	    ClientID:    "your-client-id",
//	    RedirectURI: "https://yourapp.com/callback",
//	    Provider:    "google",
//	    LoginHint:   "user@gmail.com",
//	    Scopes:      []string{"email", "calendar"},
//	    State:       "random-state-for-csrf-protection",
//	    AccessType:  "offline", // Required for refresh tokens
//	})
func (s *AuthService) URLForOAuth2(config *auth.URLForAuthenticationConfig) string {
	builder := auth.NewURLBuilder(s.client.BaseURL + "/v3/connect/auth")

	// Always set default response_type
	responseType := "code"

	if config != nil {
		builder.Add("client_id", config.ClientID)
		builder.Add("redirect_uri", config.RedirectURI)
		if config.ResponseType != "" {
			responseType = config.ResponseType
		}
		builder.Add("provider", config.Provider)
		builder.Add("login_hint", config.LoginHint)
		builder.Add("state", config.State)
		builder.Add("access_type", config.AccessType)
		builder.Add("prompt", config.Prompt)
		builder.Add("credential_id", config.CredentialID)
		builder.AddBool("include_grant_scopes", config.IncludeGrantScopes)
		if len(config.Scopes) > 0 {
			builder.Add("scope", strings.Join(config.Scopes, " "))
		}
	}

	builder.Add("response_type", responseType)
	return builder.Build()
}

// URLForOAuth2PKCE generates an OAuth 2.0 authorization URL with PKCE support.
// Use this for public clients (mobile apps, SPAs) where the client secret cannot be securely stored.
//
// Example:
//
//	// Generate a code verifier and challenge (use crypto/rand in production)
//	codeVerifier := "random-43-to-128-character-string"
//	h := sha256.Sum256([]byte(codeVerifier))
//	codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])
//
//	authURL := client.Auth.URLForOAuth2PKCE(&auth.PKCEURLConfig{
//	    ClientID:            "your-client-id",
//	    RedirectURI:         "https://yourapp.com/callback",
//	    CodeChallenge:       codeChallenge,
//	    CodeChallengeMethod: "S256",
//	    State:               "random-state",
//	})
//	// Redirect user to authURL, then exchange code with codeVerifier
func (s *AuthService) URLForOAuth2PKCE(config *auth.PKCEURLConfig) string {
	builder := auth.NewURLBuilder(s.client.BaseURL + "/v3/connect/auth")

	// Always set default response_type
	responseType := "code"

	if config != nil {
		builder.Add("client_id", config.ClientID)
		builder.Add("redirect_uri", config.RedirectURI)
		if config.ResponseType != "" {
			responseType = config.ResponseType
		}
		builder.Add("provider", config.Provider)
		builder.Add("login_hint", config.LoginHint)
		builder.Add("state", config.State)
		builder.Add("access_type", config.AccessType)
		builder.Add("prompt", config.Prompt)
		builder.Add("credential_id", config.CredentialID)
		builder.AddBool("include_grant_scopes", config.IncludeGrantScopes)
		builder.Add("code_challenge", config.CodeChallenge)
		builder.Add("code_challenge_method", config.CodeChallengeMethod)
		if len(config.Scopes) > 0 {
			builder.Add("scope", strings.Join(config.Scopes, " "))
		}
	}

	builder.Add("response_type", responseType)
	return builder.Build()
}

// URLForAdminConsent generates an OAuth 2.0 admin consent URL for Microsoft.
// This is used for Microsoft 365 tenant-wide admin consent.
func (s *AuthService) URLForAdminConsent(config *auth.AdminConsentURLConfig) string {
	builder := auth.NewURLBuilder(s.client.BaseURL + "/v3/connect/auth")

	if config != nil {
		builder.Add("client_id", config.ClientID)
		builder.Add("redirect_uri", config.RedirectURI)
		builder.Add("state", config.State)
		builder.Add("credential_id", config.CredentialID)
		builder.Add("response_type", "adminconsent")
	}

	return builder.Build()
}

// ExchangeCodeForToken exchanges an authorization code for access and refresh tokens.
// This is the final step of the OAuth 2.0 authorization code flow.
//
// Example:
//
//	// After user is redirected back with ?code=xxx
//	tokens, err := client.Auth.ExchangeCodeForToken(ctx, &auth.CodeExchangeRequest{
//	    ClientID:    "your-client-id",
//	    RedirectURI: "https://yourapp.com/callback",
//	    Code:        code, // From query parameter
//	})
//	if err != nil {
//	    return err
//	}
//
//	// Store tokens.GrantID to access user's data
//	// tokens.AccessToken expires in 1 hour
//	// tokens.RefreshToken can be used to get new access tokens
func (s *AuthService) ExchangeCodeForToken(ctx context.Context, req *auth.CodeExchangeRequest) (*auth.TokenExchangeResponse, error) {
	path := "/v3/connect/token"

	// Auto-inject client secret and grant_type
	if req.ClientSecret == "" {
		req.ClientSecret = s.client.APIKey
	}
	req.GrantType = "authorization_code"

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("auth.ExchangeCodeForToken: %w", err)
	}

	var resp auth.TokenExchangeResponse
	if err := s.client.DoRaw(httpReq, &resp); err != nil {
		return nil, fmt.Errorf("auth.ExchangeCodeForToken: %w", err)
	}

	return &resp, nil
}

// RefreshAccessToken refreshes an expired access token using a refresh token.
// Access tokens expire after one hour; use this method to get a new one.
//
// Example:
//
//	tokens, err := client.Auth.RefreshAccessToken(ctx, &auth.RefreshTokenRequest{
//	    ClientID:     "your-client-id",
//	    RefreshToken: storedRefreshToken,
//	})
//	if err != nil {
//	    // Handle error - user may need to re-authenticate
//	    return err
//	}
//	// Update stored tokens with new values
func (s *AuthService) RefreshAccessToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.TokenExchangeResponse, error) {
	path := "/v3/connect/token"

	// Auto-inject client secret and grant_type
	if req.ClientSecret == "" {
		req.ClientSecret = s.client.APIKey
	}
	req.GrantType = "refresh_token"

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("auth.RefreshAccessToken: %w", err)
	}

	var resp auth.TokenExchangeResponse
	if err := s.client.DoRaw(httpReq, &resp); err != nil {
		return nil, fmt.Errorf("auth.RefreshAccessToken: %w", err)
	}

	return &resp, nil
}

// CustomAuthentication performs custom/native authentication to create a grant.
// Use this for server-side authentication where you already have provider tokens.
func (s *AuthService) CustomAuthentication(ctx context.Context, req *auth.CustomAuthRequest) (*grants.Grant, error) {
	path := "/v3/connect/custom"

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("auth.CustomAuthentication: %w", err)
	}

	var grant grants.Grant
	_, err = s.client.Do(httpReq, &grant)
	if err != nil {
		return nil, fmt.Errorf("auth.CustomAuthentication: %w", err)
	}

	return &grant, nil
}

// IDTokenInfo validates an ID token and returns information about it.
func (s *AuthService) IDTokenInfo(ctx context.Context, idToken string) (*auth.TokenInfoResponse, error) {
	path := "/v3/connect/tokeninfo"

	httpReq, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("auth.IDTokenInfo: %w", err)
	}

	q := httpReq.URL.Query()
	q.Set("id_token", idToken)
	httpReq.URL.RawQuery = q.Encode()

	var resp auth.TokenInfoResponse
	_, err = s.client.Do(httpReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("auth.IDTokenInfo: %w", err)
	}

	return &resp, nil
}

// ValidateAccessToken validates an access token and returns information about it.
// Deprecated: Use AccessTokenInfo instead.
func (s *AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (*auth.TokenInfoResponse, error) {
	return s.AccessTokenInfo(ctx, accessToken)
}

// AccessTokenInfo retrieves information about an access token.
func (s *AuthService) AccessTokenInfo(ctx context.Context, accessToken string) (*auth.TokenInfoResponse, error) {
	path := "/v3/connect/tokeninfo"

	httpReq, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("auth.AccessTokenInfo: %w", err)
	}

	q := httpReq.URL.Query()
	q.Set("access_token", accessToken)
	httpReq.URL.RawQuery = q.Encode()

	var resp auth.TokenInfoResponse
	_, err = s.client.Do(httpReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("auth.AccessTokenInfo: %w", err)
	}

	return &resp, nil
}

// Revoke revokes an access token or refresh token.
// Returns true if the token was successfully revoked.
func (s *AuthService) Revoke(ctx context.Context, token string) error {
	path := "/v3/connect/revoke"

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("auth.Revoke: %w", err)
	}

	q := httpReq.URL.Query()
	q.Set("token", token)
	httpReq.URL.RawQuery = q.Encode()

	if err := s.client.DoRaw(httpReq, nil); err != nil {
		return fmt.Errorf("auth.Revoke: %w", err)
	}

	return nil
}

// DetectProvider detects the email provider for an email address.
func (s *AuthService) DetectProvider(ctx context.Context, req *auth.ProviderDetectRequest) (*auth.ProviderDetectResponse, error) {
	path := "/v3/providers/detect"

	httpReq, err := s.client.NewRequest(ctx, http.MethodPost, path, req)
	if err != nil {
		return nil, fmt.Errorf("auth.DetectProvider: %w", err)
	}

	var resp auth.ProviderDetectResponse
	_, err = s.client.Do(httpReq, &resp)
	if err != nil {
		return nil, fmt.Errorf("auth.DetectProvider: %w", err)
	}

	return &resp, nil
}
