//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	nylas "github.com/mqasimca/nylas-go"
)

// Provider represents a configured email provider for testing.
type Provider struct {
	Name    string
	GrantID string
}

// TestConfig holds integration test configuration.
type TestConfig struct {
	APIKey    string
	ClientID  string
	Providers []Provider
}

// providerEnvVars maps provider names to their environment variable names.
var providerEnvVars = map[string]string{
	"Google":    "NYLAS_GOOGLE_GRANT_ID",
	"Microsoft": "NYLAS_MICROSOFT_GRANT_ID",
	"ICloud":    "NYLAS_ICLOUD_GRANT_ID",
	"Yahoo":     "NYLAS_YAHOO_GRANT_ID",
	"IMAP":      "NYLAS_IMAP_GRANT_ID",
	"EWS":       "NYLAS_EWS_GRANT_ID",
}

// LoadConfig loads configuration from environment variables.
func LoadConfig(t *testing.T) *TestConfig {
	t.Helper()

	apiKey := os.Getenv("NYLAS_API_KEY")
	if apiKey == "" {
		t.Skip("NYLAS_API_KEY not set, skipping integration tests")
	}

	// Load all configured providers
	var providers []Provider
	for name, envVar := range providerEnvVars {
		if grantID := os.Getenv(envVar); grantID != "" {
			providers = append(providers, Provider{Name: name, GrantID: grantID})
		}
	}

	if len(providers) == 0 {
		t.Skip("No provider grant IDs configured (set NYLAS_GOOGLE_GRANT_ID, NYLAS_MICROSOFT_GRANT_ID, etc.)")
	}

	return &TestConfig{
		APIKey:    apiKey,
		ClientID:  os.Getenv("NYLAS_CLIENT_ID"),
		Providers: providers,
	}
}

// RunForEachProvider runs a test function for each configured provider.
func RunForEachProvider(t *testing.T, cfg *TestConfig, testFn func(t *testing.T, grantID string)) {
	t.Helper()
	for _, provider := range cfg.Providers {
		t.Run(provider.Name, func(t *testing.T) {
			testFn(t, provider.GrantID)
		})
	}
}

// NewTestClient creates a new Nylas client for integration tests.
// Configured with higher retries and longer waits to handle rate limits.
func NewTestClient(t *testing.T, cfg *TestConfig) *nylas.Client {
	t.Helper()

	client, err := nylas.NewClient(
		nylas.WithAPIKey(cfg.APIKey),
		nylas.WithTimeout(60*time.Second),
		nylas.WithMaxRetries(5),            // More retries for rate limits
		nylas.WithRetryWait(2*time.Second), // Longer base wait for 429s
	)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

// NewTestContext returns a context with timeout for integration tests.
func NewTestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// Cleanup is a helper to register cleanup functions.
type Cleanup struct {
	t     *testing.T
	funcs []func()
}

// NewCleanup creates a new cleanup helper.
func NewCleanup(t *testing.T) *Cleanup {
	c := &Cleanup{t: t}
	t.Cleanup(func() {
		for i := len(c.funcs) - 1; i >= 0; i-- {
			c.funcs[i]()
		}
	})
	return c
}

// Add registers a cleanup function to run after the test.
func (c *Cleanup) Add(fn func()) {
	c.funcs = append(c.funcs, fn)
}
