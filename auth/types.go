package auth

import (
	"net/url"
	"strconv"
)

// URLForAuthenticationConfig contains parameters for generating OAuth2 authorization URLs.
type URLForAuthenticationConfig struct {
	// ClientID is your Nylas application's client ID (required).
	ClientID string

	// RedirectURI is your application's callback URI (required).
	RedirectURI string

	// Provider hints the provider to use (optional).
	// Values: "google", "microsoft", "imap", etc.
	Provider string

	// LoginHint pre-fills the user's email address (optional).
	LoginHint string

	// State is an opaque value to maintain state (optional).
	State string

	// Scopes are the OAuth scopes to request (optional).
	Scopes []string

	// AccessType determines if a refresh token is returned (optional).
	// Set to "offline" to get a refresh token.
	AccessType string

	// ResponseType is the OAuth response type (optional).
	// Default is "code".
	ResponseType string

	// Prompt controls the consent screen behavior (optional).
	// Values: "select_provider", "detect", etc.
	Prompt string

	// IncludeGrantScopes includes existing grant scopes (optional).
	IncludeGrantScopes bool

	// CredentialID is the connector credential ID (optional).
	CredentialID string
}

// Values returns the configuration as query parameters.
func (c *URLForAuthenticationConfig) Values() map[string]any {
	if c == nil {
		return nil
	}

	v := make(map[string]any)
	if c.ClientID != "" {
		v["client_id"] = c.ClientID
	}
	if c.RedirectURI != "" {
		v["redirect_uri"] = c.RedirectURI
	}
	if c.Provider != "" {
		v["provider"] = c.Provider
	}
	if c.LoginHint != "" {
		v["login_hint"] = c.LoginHint
	}
	if c.State != "" {
		v["state"] = c.State
	}
	if len(c.Scopes) > 0 {
		v["scope"] = c.Scopes
	}
	if c.AccessType != "" {
		v["access_type"] = c.AccessType
	}
	if c.ResponseType != "" {
		v["response_type"] = c.ResponseType
	}
	if c.Prompt != "" {
		v["prompt"] = c.Prompt
	}
	if c.IncludeGrantScopes {
		v["include_grant_scopes"] = true
	}
	if c.CredentialID != "" {
		v["credential_id"] = c.CredentialID
	}
	return v
}

// PKCEURLConfig extends URLForAuthenticationConfig with PKCE parameters.
type PKCEURLConfig struct {
	URLForAuthenticationConfig

	// CodeChallenge is the PKCE code challenge (required for PKCE).
	CodeChallenge string

	// CodeChallengeMethod is the PKCE method, usually "S256" (required for PKCE).
	CodeChallengeMethod string
}

// Values returns the configuration as query parameters.
func (c *PKCEURLConfig) Values() map[string]any {
	if c == nil {
		return nil
	}

	v := c.URLForAuthenticationConfig.Values()
	if v == nil {
		v = make(map[string]any)
	}
	if c.CodeChallenge != "" {
		v["code_challenge"] = c.CodeChallenge
	}
	if c.CodeChallengeMethod != "" {
		v["code_challenge_method"] = c.CodeChallengeMethod
	}
	return v
}

// AdminConsentURLConfig contains parameters for generating admin consent URLs.
type AdminConsentURLConfig struct {
	// ClientID is your Nylas application's client ID (required).
	ClientID string

	// RedirectURI is your application's callback URI (required).
	RedirectURI string

	// State is an opaque value to maintain state (optional).
	State string

	// CredentialID is the connector credential ID (required).
	CredentialID string
}

// Values returns the configuration as query parameters.
func (c *AdminConsentURLConfig) Values() map[string]any {
	if c == nil {
		return nil
	}

	v := make(map[string]any)
	if c.ClientID != "" {
		v["client_id"] = c.ClientID
	}
	if c.RedirectURI != "" {
		v["redirect_uri"] = c.RedirectURI
	}
	if c.State != "" {
		v["state"] = c.State
	}
	if c.CredentialID != "" {
		v["credential_id"] = c.CredentialID
	}
	v["response_type"] = "adminconsent"
	return v
}

// CodeExchangeRequest contains parameters for exchanging an authorization code for tokens.
type CodeExchangeRequest struct {
	// ClientID is your Nylas application's client ID (required).
	ClientID string `json:"client_id"`

	// ClientSecret is your Nylas API key (required unless using PKCE).
	ClientSecret string `json:"client_secret,omitempty"`

	// Code is the authorization code from the provider (required).
	Code string `json:"code"`

	// RedirectURI is your application's callback URI (required).
	RedirectURI string `json:"redirect_uri"`

	// GrantType should be "authorization_code" (required).
	GrantType string `json:"grant_type"`

	// CodeVerifier is the PKCE code verifier (required for PKCE).
	CodeVerifier string `json:"code_verifier,omitempty"`
}

// TokenExchangeResponse contains the response from token exchange operations.
type TokenExchangeResponse struct {
	// AccessToken is the OAuth access token.
	AccessToken string `json:"access_token,omitempty"`

	// RefreshToken is the OAuth refresh token (if access_type=offline).
	RefreshToken string `json:"refresh_token,omitempty"`

	// GrantID is the unique identifier for the grant.
	GrantID string `json:"grant_id,omitempty"`

	// Email is the user's email address.
	Email string `json:"email,omitempty"`

	// TokenType is the token type, usually "Bearer".
	TokenType string `json:"token_type,omitempty"`

	// ExpiresIn is the token lifetime in seconds.
	ExpiresIn int `json:"expires_in,omitempty"`

	// IDToken is the OpenID Connect ID token.
	IDToken string `json:"id_token,omitempty"`

	// Scope is the authorized scopes.
	Scope string `json:"scope,omitempty"`

	// Provider is the email provider.
	Provider string `json:"provider,omitempty"`
}

// RefreshTokenRequest contains parameters for refreshing an access token.
type RefreshTokenRequest struct {
	// ClientID is your Nylas application's client ID (required).
	ClientID string `json:"client_id"`

	// ClientSecret is your Nylas API key (required).
	ClientSecret string `json:"client_secret,omitempty"`

	// RefreshToken is the refresh token to use (required).
	RefreshToken string `json:"refresh_token"`

	// GrantType should be "refresh_token" (required).
	GrantType string `json:"grant_type"`
}

// CustomAuthRequest contains parameters for custom/native authentication.
type CustomAuthRequest struct {
	// Provider is the email provider (required).
	// Values: "google", "microsoft", "imap", "virtual-calendar", etc.
	Provider string `json:"provider"`

	// Settings contains provider-specific authentication settings (required).
	Settings map[string]any `json:"settings"`

	// Scopes are the OAuth scopes to request (optional).
	Scopes []string `json:"scope,omitempty"`

	// State is an opaque value to maintain state (optional).
	State string `json:"state,omitempty"`
}

// TokenInfoResponse contains information about a token.
type TokenInfoResponse struct {
	// Iss is the token issuer.
	Iss string `json:"iss,omitempty"`

	// Aud is the token audience.
	Aud string `json:"aud,omitempty"`

	// Sub is the subject (grant ID).
	Sub string `json:"sub,omitempty"`

	// Email is the user's email address.
	Email string `json:"email,omitempty"`

	// EmailVerified indicates if the email is verified.
	EmailVerified bool `json:"email_verified,omitempty"`

	// AtHash is the access token hash.
	AtHash string `json:"at_hash,omitempty"`

	// Iat is the token issue time (Unix timestamp).
	Iat int64 `json:"iat,omitempty"`

	// Exp is the token expiration time (Unix timestamp).
	Exp int64 `json:"exp,omitempty"`
}

// RevokeRequest contains parameters for revoking a token.
type RevokeRequest struct {
	// Token is the token to revoke (required).
	Token string
}

// ProviderDetectRequest contains parameters for detecting a provider.
type ProviderDetectRequest struct {
	// Email is the email address to detect the provider for (required).
	Email string `json:"email"`

	// AllProviderTypes includes all provider types in detection (optional).
	AllProviderTypes bool `json:"all_provider_types,omitempty"`
}

// ProviderDetectResponse contains the provider detection result.
type ProviderDetectResponse struct {
	// EmailAddress is the email address that was checked.
	EmailAddress string `json:"email_address,omitempty"`

	// Detected indicates if a provider was detected.
	Detected bool `json:"detected,omitempty"`

	// Provider is the detected provider type.
	Provider string `json:"provider,omitempty"`

	// Type is the provider type classification.
	Type string `json:"type,omitempty"`
}

// URLBuilder helps construct OAuth URLs with query parameters.
type URLBuilder struct {
	baseURL string
	params  url.Values
}

// NewURLBuilder creates a new URL builder.
func NewURLBuilder(baseURL string) *URLBuilder {
	return &URLBuilder{
		baseURL: baseURL,
		params:  make(url.Values),
	}
}

// Add adds a parameter to the URL.
func (b *URLBuilder) Add(key, value string) *URLBuilder {
	if value != "" {
		b.params.Set(key, value)
	}
	return b
}

// AddBool adds a boolean parameter to the URL.
func (b *URLBuilder) AddBool(key string, value bool) *URLBuilder {
	if value {
		b.params.Set(key, strconv.FormatBool(value))
	}
	return b
}

// AddSlice adds a slice parameter to the URL as comma-separated values.
func (b *URLBuilder) AddSlice(key string, values []string) *URLBuilder {
	for _, v := range values {
		b.params.Add(key, v)
	}
	return b
}

// Build returns the final URL string.
func (b *URLBuilder) Build() string {
	if len(b.params) == 0 {
		return b.baseURL
	}
	return b.baseURL + "?" + b.params.Encode()
}
