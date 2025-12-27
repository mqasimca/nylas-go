package nylas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mqasimca/nylas-go/webhooks"
)

// List returns all webhooks.
func (s *WebhooksService) List(ctx context.Context, opts *webhooks.ListOptions) (*ListResponse[webhooks.Webhook], error) {
	path := "/v3/webhooks"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("webhooks.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []webhooks.Webhook
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("webhooks.List: %w", err)
	}

	return &ListResponse[webhooks.Webhook]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single webhook.
func (s *WebhooksService) Get(ctx context.Context, webhookID string) (*webhooks.Webhook, error) {
	path := fmt.Sprintf("/v3/webhooks/%s", webhookID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Get(%s): %w", webhookID, err)
	}

	var webhook webhooks.Webhook
	_, err = s.client.Do(req, &webhook)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Get(%s): %w", webhookID, err)
	}

	return &webhook, nil
}

// Create creates a new webhook.
//
// The webhook URL must be publicly accessible and respond to challenge requests.
// Nylas will send a GET request with a ?challenge= parameter that must be echoed back.
//
// Example:
//
//	webhook, err := client.Webhooks.Create(ctx, &webhooks.CreateRequest{
//	    WebhookURL:   "https://yourapp.com/webhooks/nylas",
//	    TriggerTypes: []string{"message.created", "message.updated"},
//	    Description:  "Email notifications",
//	    NotificationEmailAddresses: []string{"alerts@yourcompany.com"},
//	})
//	if err != nil {
//	    return err
//	}
//	// Store webhook.WebhookSecret securely for signature verification
func (s *WebhooksService) Create(ctx context.Context, create *webhooks.CreateRequest) (*webhooks.Webhook, error) {
	path := "/v3/webhooks"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Create: %w", err)
	}

	var webhook webhooks.Webhook
	_, err = s.client.Do(req, &webhook)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Create: %w", err)
	}

	return &webhook, nil
}

// Update updates a webhook.
//
// Example:
//
//	webhook, err := client.Webhooks.Update(ctx, webhookID, &webhooks.UpdateRequest{
//	    TriggerTypes: []string{"message.created", "calendar.created"},
//	    Description:  "Updated webhook description",
//	})
func (s *WebhooksService) Update(ctx context.Context, webhookID string, update *webhooks.UpdateRequest) (*webhooks.Webhook, error) {
	path := fmt.Sprintf("/v3/webhooks/%s", webhookID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Update(%s): %w", webhookID, err)
	}

	var webhook webhooks.Webhook
	_, err = s.client.Do(req, &webhook)
	if err != nil {
		return nil, fmt.Errorf("webhooks.Update(%s): %w", webhookID, err)
	}

	return &webhook, nil
}

// Delete deletes a webhook.
func (s *WebhooksService) Delete(ctx context.Context, webhookID string) error {
	path := fmt.Sprintf("/v3/webhooks/%s", webhookID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("webhooks.Delete(%s): %w", webhookID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("webhooks.Delete(%s): %w", webhookID, err)
	}

	return nil
}

// RotateSecret rotates a webhook's secret.
//
// Use this periodically or after a suspected compromise to generate a new webhook secret.
// After rotation, update your webhook handler to use the new secret for signature verification.
//
// Example:
//
//	result, err := client.Webhooks.RotateSecret(ctx, webhookID)
//	if err != nil {
//	    return err
//	}
//	// Update your stored secret with result.WebhookSecret
func (s *WebhooksService) RotateSecret(ctx context.Context, webhookID string) (*webhooks.RotateSecretResponse, error) {
	path := fmt.Sprintf("/v3/webhooks/rotate-secret/%s", webhookID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("webhooks.RotateSecret(%s): %w", webhookID, err)
	}

	var result webhooks.RotateSecretResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("webhooks.RotateSecret(%s): %w", webhookID, err)
	}

	return &result, nil
}

// ListAll returns an iterator for all webhooks.
func (s *WebhooksService) ListAll(ctx context.Context, opts *webhooks.ListOptions) *Iterator[webhooks.Webhook] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]webhooks.Webhook, string, error) {
		o := opts
		if o == nil {
			o = &webhooks.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}

// GetIPAddresses returns the list of Nylas IP addresses used for sending webhooks.
// Use this to whitelist Nylas IPs in your firewall configuration.
// Note: This endpoint is available for paid customers only.
//
// Example:
//
//	ips, err := client.Webhooks.GetIPAddresses(ctx)
//	if err != nil {
//	    return err
//	}
//	for _, ip := range ips.IPAddresses {
//	    fmt.Println("Whitelist:", ip)
//	}
func (s *WebhooksService) GetIPAddresses(ctx context.Context) (*webhooks.IPAddressesResponse, error) {
	path := "/v3/webhooks/ip-addresses"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("webhooks.GetIPAddresses: %w", err)
	}

	var result webhooks.IPAddressesResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("webhooks.GetIPAddresses: %w", err)
	}

	return &result, nil
}

// ExtractChallengeParameter extracts the challenge parameter from a webhook validation URL.
// When Nylas validates a webhook endpoint, it sends a GET request with a challenge query parameter.
// Your endpoint must return this challenge value in the response body within 10 seconds.
//
// Example URL: https://your-webhook.com/endpoint?challenge=abc123
// Returns: "abc123", nil
func ExtractChallengeParameter(webhookURL string) (string, error) {
	parsed, err := url.Parse(webhookURL)
	if err != nil {
		return "", fmt.Errorf("invalid webhook URL: %w", err)
	}

	challenge := parsed.Query().Get("challenge")
	if challenge == "" {
		return "", fmt.Errorf("no challenge parameter found in URL")
	}

	return challenge, nil
}
