package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// Default configuration values for the Nylas API client.
const (
	defaultBaseURL   = "https://api.us.nylas.com"
	defaultTimeout   = 90 * time.Second
	defaultRetries   = 2
	defaultRetryWait = 500 * time.Millisecond
)

// Client is the Nylas API client. It is safe for concurrent use by multiple goroutines.
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client

	MaxRetries int
	RetryWait  time.Duration

	Messages     *MessagesService
	Threads      *ThreadsService
	Drafts       *DraftsService
	Calendars    *CalendarsService
	Events       *EventsService
	Contacts     *ContactsService
	Folders      *FoldersService
	Attachments  *AttachmentsService
	Grants       *GrantsService
	Webhooks     *WebhooksService
	Auth         *AuthService
	Scheduler    *SchedulerService
	Notetakers   *NotetakersService
	Applications *ApplicationsService
	RedirectURIs *RedirectURIsService
	Connectors   *ConnectorsService
	Credentials  *CredentialsService
	SmartCompose *SmartComposeService

	rateMu     sync.Mutex
	rateLimits Rate

	common service
}

type service struct {
	client *Client
}

// NewClient creates a new Nylas API client. An API key is required; use WithAPIKey to provide it.
//
// Example:
//
//	client, err := nylas.NewClient(
//	    nylas.WithAPIKey("your-api-key"),
//	    nylas.WithRegion(nylas.RegionEU),
//	)
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		MaxRetries: defaultRetries,
		RetryWait:  defaultRetryWait,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.APIKey == "" {
		return nil, ErrMissingAPIKey
	}

	c.common.client = c
	c.Messages = (*MessagesService)(&c.common)
	c.Threads = (*ThreadsService)(&c.common)
	c.Drafts = (*DraftsService)(&c.common)
	c.Calendars = (*CalendarsService)(&c.common)
	c.Events = (*EventsService)(&c.common)
	c.Contacts = (*ContactsService)(&c.common)
	c.Folders = (*FoldersService)(&c.common)
	c.Attachments = (*AttachmentsService)(&c.common)
	c.Grants = (*GrantsService)(&c.common)
	c.Webhooks = (*WebhooksService)(&c.common)
	c.Auth = (*AuthService)(&c.common)
	c.Scheduler = (*SchedulerService)(&c.common)
	c.Notetakers = (*NotetakersService)(&c.common)
	c.Applications = (*ApplicationsService)(&c.common)
	c.RedirectURIs = (*RedirectURIsService)(&c.common)
	c.Connectors = (*ConnectorsService)(&c.common)
	c.Credentials = (*CredentialsService)(&c.common)
	c.SmartCompose = (*SmartComposeService)(&c.common)

	return c, nil
}

// NewRequest creates an HTTP request for the Nylas API with proper headers and authentication.
func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// doWithRetry executes an HTTP request with exponential backoff retry on 5xx and 429 errors.
func (c *Client) doWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			if attempt < c.MaxRetries {
				time.Sleep(c.RetryWait * time.Duration(1<<attempt))
				continue
			}
			return nil, err
		}

		// Don't retry on success or client errors (except 429)
		if resp.StatusCode < 500 && resp.StatusCode != 429 {
			return resp, nil
		}

		_ = resp.Body.Close()

		if attempt < c.MaxRetries {
			wait := c.RetryWait * time.Duration(1<<attempt)

			// Use Retry-After header if present (for 429)
			if resp.StatusCode == 429 {
				if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
					if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil {
						wait = time.Duration(seconds) * time.Second
					}
				}
			}

			time.Sleep(wait)
		}
	}

	return resp, nil
}

// Do executes an HTTP request and decodes the JSON response into v.
// The response is expected to be wrapped in a standard Nylas response envelope with data and request_id fields.
func (c *Client) Do(req *http.Request, v any) (*Response[any], error) {
	resp, err := c.doWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	c.updateRateLimits(resp)

	if resp.StatusCode >= 400 {
		return nil, parseError(resp)
	}

	if v == nil {
		return &Response[any]{RequestID: resp.Header.Get("X-Request-Id")}, nil
	}

	var result struct {
		Data      json.RawMessage `json:"data"`
		RequestID string          `json:"request_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(result.Data, v); err != nil {
		return nil, err
	}

	return &Response[any]{
		Data:      v,
		RequestID: result.RequestID,
	}, nil
}

// DoRaw executes a request and decodes the response directly (not wrapped in data/request_id).
// Use this for endpoints that return arrays or objects directly without the standard wrapper.
func (c *Client) DoRaw(req *http.Request, v any) error {
	resp, err := c.doWithRetry(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	c.updateRateLimits(resp)

	if resp.StatusCode >= 400 {
		return parseError(resp)
	}

	if v == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// DoList executes a request and decodes a list response with pagination.
func (c *Client) DoList(req *http.Request, v any) (string, string, error) {
	resp, err := c.doWithRetry(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()

	c.updateRateLimits(resp)

	if resp.StatusCode >= 400 {
		return "", "", parseError(resp)
	}

	var result struct {
		Data       json.RawMessage `json:"data"`
		RequestID  string          `json:"request_id"`
		NextCursor string          `json:"next_cursor"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	if err := json.Unmarshal(result.Data, v); err != nil {
		return "", "", err
	}

	return result.NextCursor, result.RequestID, nil
}

func (c *Client) updateRateLimits(resp *http.Response) {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()
	c.rateLimits = parseRateLimits(resp)
}

// RateLimits returns the current rate limit information from the last API response.
func (c *Client) RateLimits() Rate {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()
	return c.rateLimits
}

func parseError(resp *http.Response) error {
	var apiErr APIError
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		return fmt.Errorf("nylas: request failed with status %d", resp.StatusCode)
	}
	apiErr.StatusCode = resp.StatusCode
	apiErr.RequestID = resp.Header.Get("X-Request-Id")
	return &apiErr
}

// Service type aliases for accessing Nylas API resources.
type (
	// MessagesService handles operations on email messages.
	MessagesService service
	// ThreadsService handles operations on email threads.
	ThreadsService service
	// DraftsService handles operations on email drafts.
	DraftsService service
	// CalendarsService handles operations on calendars.
	CalendarsService service
	// EventsService handles operations on calendar events.
	EventsService service
	// ContactsService handles operations on contacts.
	ContactsService service
	// FoldersService handles operations on email folders/labels.
	FoldersService service
	// AttachmentsService handles operations on email attachments.
	AttachmentsService service
	// GrantsService handles operations on connected accounts (grants).
	GrantsService service
	// WebhooksService handles operations on webhook subscriptions.
	WebhooksService service
	// AuthService handles authentication operations.
	AuthService service
	// SchedulerService handles scheduling operations.
	SchedulerService service
	// NotetakersService handles notetaker operations.
	NotetakersService service
	// ApplicationsService handles application configuration operations.
	ApplicationsService service
	// RedirectURIsService handles redirect URI operations.
	RedirectURIsService service
	// ConnectorsService handles connector operations.
	ConnectorsService service
	// CredentialsService handles credential operations.
	CredentialsService service
	// SmartComposeService handles AI-powered message composition.
	SmartComposeService service
)
