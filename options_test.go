package nylas

import (
	"net/http"
	"testing"
	"time"
)

func TestWithAPIKey(t *testing.T) {
	client, _ := NewClient(WithAPIKey("my-api-key"))
	if client.APIKey != "my-api-key" {
		t.Errorf("WithAPIKey() = %v, want my-api-key", client.APIKey)
	}
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
	}{
		{"without trailing slash", "https://custom.api.com", "https://custom.api.com"},
		{"with trailing slash", "https://custom.api.com/", "https://custom.api.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient(WithAPIKey("key"), WithBaseURL(tt.baseURL))
			if client.BaseURL != tt.want {
				t.Errorf("WithBaseURL() = %v, want %v", client.BaseURL, tt.want)
			}
		})
	}
}

func TestWithRegion(t *testing.T) {
	tests := []struct {
		name   string
		region Region
		want   string
	}{
		{"US region", RegionUS, "https://api.us.nylas.com"},
		{"EU region", RegionEU, "https://api.eu.nylas.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := NewClient(WithAPIKey("key"), WithRegion(tt.region))
			if client.BaseURL != tt.want {
				t.Errorf("WithRegion() = %v, want %v", client.BaseURL, tt.want)
			}
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Second}
	client, _ := NewClient(WithAPIKey("key"), WithHTTPClient(customClient))
	if client.HTTPClient != customClient {
		t.Error("WithHTTPClient() did not set custom client")
	}
}

func TestWithTimeout(t *testing.T) {
	client, _ := NewClient(WithAPIKey("key"), WithTimeout(30*time.Second))
	if client.HTTPClient.Timeout != 30*time.Second {
		t.Errorf("WithTimeout() = %v, want 30s", client.HTTPClient.Timeout)
	}
}

func TestWithMaxRetries(t *testing.T) {
	client, _ := NewClient(WithAPIKey("key"), WithMaxRetries(5))
	if client.MaxRetries != 5 {
		t.Errorf("WithMaxRetries() = %v, want 5", client.MaxRetries)
	}
}

func TestWithRetryWait(t *testing.T) {
	client, _ := NewClient(WithAPIKey("key"), WithRetryWait(time.Second))
	if client.RetryWait != time.Second {
		t.Errorf("WithRetryWait() = %v, want 1s", client.RetryWait)
	}
}
