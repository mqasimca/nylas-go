package nylas

import (
	"context"
	"fmt"
	"net/http"

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
