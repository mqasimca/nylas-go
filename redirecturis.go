package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/redirecturis"
)

// List returns all redirect URIs for the application.
func (s *RedirectURIsService) List(ctx context.Context, opts *redirecturis.ListOptions) (*ListResponse[redirecturis.RedirectURI], error) {
	path := "/v3/applications/redirect-uris"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []redirecturis.RedirectURI
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.List: %w", err)
	}

	return &ListResponse[redirecturis.RedirectURI]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single redirect URI by ID.
func (s *RedirectURIsService) Get(ctx context.Context, redirectURIID string) (*redirecturis.RedirectURI, error) {
	path := fmt.Sprintf("/v3/applications/redirect-uris/%s", redirectURIID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Get(%s): %w", redirectURIID, err)
	}

	var uri redirecturis.RedirectURI
	_, err = s.client.Do(req, &uri)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Get(%s): %w", redirectURIID, err)
	}

	return &uri, nil
}

// Create creates a new redirect URI.
func (s *RedirectURIsService) Create(ctx context.Context, create *redirecturis.CreateRequest) (*redirecturis.RedirectURI, error) {
	path := "/v3/applications/redirect-uris"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Create: %w", err)
	}

	var uri redirecturis.RedirectURI
	_, err = s.client.Do(req, &uri)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Create: %w", err)
	}

	return &uri, nil
}

// Update updates a redirect URI.
func (s *RedirectURIsService) Update(ctx context.Context, redirectURIID string, update *redirecturis.UpdateRequest) (*redirecturis.RedirectURI, error) {
	path := fmt.Sprintf("/v3/applications/redirect-uris/%s", redirectURIID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Update(%s): %w", redirectURIID, err)
	}

	var uri redirecturis.RedirectURI
	_, err = s.client.Do(req, &uri)
	if err != nil {
		return nil, fmt.Errorf("redirecturis.Update(%s): %w", redirectURIID, err)
	}

	return &uri, nil
}

// Delete deletes a redirect URI.
func (s *RedirectURIsService) Delete(ctx context.Context, redirectURIID string) error {
	path := fmt.Sprintf("/v3/applications/redirect-uris/%s", redirectURIID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("redirecturis.Delete(%s): %w", redirectURIID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("redirecturis.Delete(%s): %w", redirectURIID, err)
	}

	return nil
}

// ListAll returns an iterator for all redirect URIs.
func (s *RedirectURIsService) ListAll(ctx context.Context, opts *redirecturis.ListOptions) *Iterator[redirecturis.RedirectURI] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]redirecturis.RedirectURI, string, error) {
		o := opts
		if o == nil {
			o = &redirecturis.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
