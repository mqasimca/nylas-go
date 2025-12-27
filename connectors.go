package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/connectors"
)

// List returns all connectors.
func (s *ConnectorsService) List(ctx context.Context, opts *connectors.ListOptions) (*ListResponse[connectors.Connector], error) {
	path := "/v3/connectors"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("connectors.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []connectors.Connector
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("connectors.List: %w", err)
	}

	return &ListResponse[connectors.Connector]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single connector by provider.
func (s *ConnectorsService) Get(ctx context.Context, provider connectors.Provider) (*connectors.Connector, error) {
	path := fmt.Sprintf("/v3/connectors/%s", provider)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("connectors.Get(%s): %w", provider, err)
	}

	var connector connectors.Connector
	_, err = s.client.Do(req, &connector)
	if err != nil {
		return nil, fmt.Errorf("connectors.Get(%s): %w", provider, err)
	}

	return &connector, nil
}

// Create creates a new connector.
func (s *ConnectorsService) Create(ctx context.Context, create *connectors.CreateRequest) (*connectors.Connector, error) {
	path := "/v3/connectors"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("connectors.Create: %w", err)
	}

	var connector connectors.Connector
	_, err = s.client.Do(req, &connector)
	if err != nil {
		return nil, fmt.Errorf("connectors.Create: %w", err)
	}

	return &connector, nil
}

// Update updates a connector.
func (s *ConnectorsService) Update(ctx context.Context, provider connectors.Provider, update *connectors.UpdateRequest) (*connectors.Connector, error) {
	path := fmt.Sprintf("/v3/connectors/%s", provider)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("connectors.Update(%s): %w", provider, err)
	}

	var connector connectors.Connector
	_, err = s.client.Do(req, &connector)
	if err != nil {
		return nil, fmt.Errorf("connectors.Update(%s): %w", provider, err)
	}

	return &connector, nil
}

// Delete deletes a connector.
func (s *ConnectorsService) Delete(ctx context.Context, provider connectors.Provider) error {
	path := fmt.Sprintf("/v3/connectors/%s", provider)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("connectors.Delete(%s): %w", provider, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("connectors.Delete(%s): %w", provider, err)
	}

	return nil
}

// ListAll returns an iterator for all connectors.
func (s *ConnectorsService) ListAll(ctx context.Context, opts *connectors.ListOptions) *Iterator[connectors.Connector] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]connectors.Connector, string, error) {
		o := opts
		if o == nil {
			o = &connectors.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
