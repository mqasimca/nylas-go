package nylas

import (
	"net/http"
	"strings"
	"time"
)

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.APIKey = key }
}

// WithBaseURL sets a custom base URL for the API. Useful for testing or proxying.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.BaseURL = strings.TrimSuffix(baseURL, "/") }
}

// WithRegion sets the API region (US or EU).
func WithRegion(region Region) Option {
	return func(c *Client) {
		switch region {
		case RegionEU:
			c.BaseURL = "https://api.eu.nylas.com"
		default:
			c.BaseURL = "https://api.us.nylas.com"
		}
	}
}

// WithHTTPClient sets a custom HTTP client for making requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.HTTPClient = hc }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{}
		}
		c.HTTPClient.Timeout = d
	}
}

// WithMaxRetries sets the maximum number of retry attempts for failed requests.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.MaxRetries = n }
}

// WithRetryWait sets the base wait time between retries (uses exponential backoff).
func WithRetryWait(d time.Duration) Option {
	return func(c *Client) { c.RetryWait = d }
}

// Region represents a Nylas API region.
type Region string

// Available Nylas API regions.
const (
	RegionUS Region = "us" // United States (default)
	RegionEU Region = "eu" // European Union
)
