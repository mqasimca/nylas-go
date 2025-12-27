// Package auth provides types and utilities for Nylas authentication operations.
//
// This package supports the complete Nylas OAuth 2.0 flow including:
//   - OAuth URL generation (with optional PKCE support)
//   - Authorization code exchange for tokens
//   - Access token refresh
//   - Token validation and introspection
//   - Token revocation
//   - Custom/native authentication
//   - Provider detection
//
// # OAuth 2.0 Flow
//
// The standard OAuth flow consists of:
//  1. Generate an authorization URL using URLForOAuth2
//  2. Redirect the user to that URL
//  3. User authenticates with their email provider
//  4. User is redirected back to your redirect_uri with a code
//  5. Exchange the code for tokens using ExchangeCodeForToken
//
// Example:
//
//	// Step 1: Generate OAuth URL
//	authURL := client.Auth.URLForOAuth2(&auth.URLForAuthenticationConfig{
//	    ClientID:    "your-client-id",
//	    RedirectURI: "https://yourapp.com/callback",
//	    AccessType:  "offline", // Request refresh token
//	})
//
//	// Step 5: Exchange code for tokens
//	tokens, err := client.Auth.ExchangeCodeForToken(ctx, &auth.CodeExchangeRequest{
//	    ClientID:     "your-client-id",
//	    ClientSecret: "your-api-key",
//	    Code:         codeFromCallback,
//	    RedirectURI:  "https://yourapp.com/callback",
//	    GrantType:    "authorization_code",
//	})
//
// # PKCE Flow
//
// For public clients (like mobile or SPA), use PKCE:
//
//	// Generate code verifier and challenge
//	verifier := generateCodeVerifier()
//	challenge := sha256Base64URL(verifier)
//
//	// Generate OAuth URL with PKCE
//	authURL := client.Auth.URLForOAuth2PKCE(&auth.PKCEURLConfig{
//	    URLForAuthenticationConfig: auth.URLForAuthenticationConfig{
//	        ClientID:    "your-client-id",
//	        RedirectURI: "https://yourapp.com/callback",
//	    },
//	    CodeChallenge:       challenge,
//	    CodeChallengeMethod: "S256",
//	})
//
//	// Exchange with code verifier
//	tokens, err := client.Auth.ExchangeCodeForToken(ctx, &auth.CodeExchangeRequest{
//	    ClientID:     "your-client-id",
//	    Code:         codeFromCallback,
//	    RedirectURI:  "https://yourapp.com/callback",
//	    GrantType:    "authorization_code",
//	    CodeVerifier: verifier,
//	})
//
// # Token Refresh
//
// Access tokens expire after one hour. Use refresh tokens to get new access tokens:
//
//	newTokens, err := client.Auth.RefreshAccessToken(ctx, &auth.RefreshTokenRequest{
//	    ClientID:     "your-client-id",
//	    ClientSecret: "your-api-key",
//	    RefreshToken: tokens.RefreshToken,
//	    GrantType:    "refresh_token",
//	})
//
// # Custom Authentication
//
// For server-side or native authentication flows:
//
//	grant, err := client.Auth.CustomAuthentication(ctx, &auth.CustomAuthRequest{
//	    Provider: "google",
//	    Settings: map[string]any{
//	        "refresh_token": "provider-refresh-token",
//	    },
//	})
package auth
