package nylas

import (
	"net/http"
	"strings"
	"time"
)

type Option func(*Client)

func WithAPIKey(key string) Option {
	return func(c *Client) { c.APIKey = key }
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.BaseURL = strings.TrimSuffix(baseURL, "/") }
}

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

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.HTTPClient = hc }
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if c.HTTPClient == nil {
			c.HTTPClient = &http.Client{}
		}
		c.HTTPClient.Timeout = d
	}
}

func WithMaxRetries(n int) Option {
	return func(c *Client) { c.MaxRetries = n }
}

func WithRetryWait(d time.Duration) Option {
	return func(c *Client) { c.RetryWait = d }
}

type Region string

const (
	RegionUS Region = "us"
	RegionEU Region = "eu"
)
