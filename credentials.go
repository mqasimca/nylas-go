package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/connectors"
	"github.com/mqasimca/nylas-go/credentials"
)

// List returns all credentials for a provider.
func (s *CredentialsService) List(ctx context.Context, provider connectors.Provider, opts *credentials.ListOptions) (*ListResponse[credentials.Credential], error) {
	path := fmt.Sprintf("/v3/connectors/%s/creds", provider)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("credentials.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []credentials.Credential
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("credentials.List: %w", err)
	}

	return &ListResponse[credentials.Credential]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single credential by ID.
func (s *CredentialsService) Get(ctx context.Context, provider connectors.Provider, credentialID string) (*credentials.Credential, error) {
	path := fmt.Sprintf("/v3/connectors/%s/creds/%s", provider, credentialID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("credentials.Get(%s): %w", credentialID, err)
	}

	var cred credentials.Credential
	_, err = s.client.Do(req, &cred)
	if err != nil {
		return nil, fmt.Errorf("credentials.Get(%s): %w", credentialID, err)
	}

	return &cred, nil
}

// Create creates a new credential for a provider.
func (s *CredentialsService) Create(ctx context.Context, provider connectors.Provider, create *credentials.CreateRequest) (*credentials.Credential, error) {
	path := fmt.Sprintf("/v3/connectors/%s/creds", provider)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("credentials.Create: %w", err)
	}

	var cred credentials.Credential
	_, err = s.client.Do(req, &cred)
	if err != nil {
		return nil, fmt.Errorf("credentials.Create: %w", err)
	}

	return &cred, nil
}

// Update updates a credential.
func (s *CredentialsService) Update(ctx context.Context, provider connectors.Provider, credentialID string, update *credentials.UpdateRequest) (*credentials.Credential, error) {
	path := fmt.Sprintf("/v3/connectors/%s/creds/%s", provider, credentialID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("credentials.Update(%s): %w", credentialID, err)
	}

	var cred credentials.Credential
	_, err = s.client.Do(req, &cred)
	if err != nil {
		return nil, fmt.Errorf("credentials.Update(%s): %w", credentialID, err)
	}

	return &cred, nil
}

// Delete deletes a credential.
func (s *CredentialsService) Delete(ctx context.Context, provider connectors.Provider, credentialID string) error {
	path := fmt.Sprintf("/v3/connectors/%s/creds/%s", provider, credentialID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("credentials.Delete(%s): %w", credentialID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("credentials.Delete(%s): %w", credentialID, err)
	}

	return nil
}

// ListAll returns an iterator for all credentials for a provider.
func (s *CredentialsService) ListAll(ctx context.Context, provider connectors.Provider, opts *credentials.ListOptions) *Iterator[credentials.Credential] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]credentials.Credential, string, error) {
		o := opts
		if o == nil {
			o = &credentials.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, provider, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
