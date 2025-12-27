package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/grants"
)

// List returns all grants.
func (s *GrantsService) List(ctx context.Context, opts *grants.ListOptions) (*ListResponse[grants.Grant], error) {
	path := "/v3/grants"

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("grants.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []grants.Grant
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("grants.List: %w", err)
	}

	return &ListResponse[grants.Grant]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single grant.
func (s *GrantsService) Get(ctx context.Context, grantID string) (*grants.Grant, error) {
	path := fmt.Sprintf("/v3/grants/%s", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("grants.Get(%s): %w", grantID, err)
	}

	var grant grants.Grant
	_, err = s.client.Do(req, &grant)
	if err != nil {
		return nil, fmt.Errorf("grants.Get(%s): %w", grantID, err)
	}

	return &grant, nil
}

// Update updates a grant.
func (s *GrantsService) Update(ctx context.Context, grantID string, update *grants.UpdateRequest) (*grants.Grant, error) {
	path := fmt.Sprintf("/v3/grants/%s", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, update)
	if err != nil {
		return nil, fmt.Errorf("grants.Update(%s): %w", grantID, err)
	}

	var grant grants.Grant
	_, err = s.client.Do(req, &grant)
	if err != nil {
		return nil, fmt.Errorf("grants.Update(%s): %w", grantID, err)
	}

	return &grant, nil
}

// Delete deletes a grant.
func (s *GrantsService) Delete(ctx context.Context, grantID string) error {
	path := fmt.Sprintf("/v3/grants/%s", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("grants.Delete(%s): %w", grantID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("grants.Delete(%s): %w", grantID, err)
	}

	return nil
}

// ListAll returns an iterator for all grants using offset-based pagination.
func (s *GrantsService) ListAll(ctx context.Context, opts *grants.ListOptions) *Iterator[grants.Grant] {
	offset := 0
	limit := 50
	if opts != nil && opts.Limit != nil {
		limit = *opts.Limit
	}

	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]grants.Grant, string, error) {
		o := &grants.ListOptions{}
		if opts != nil {
			*o = *opts
		}
		o.Offset = &offset
		o.Limit = &limit

		resp, err := s.List(ctx, o)
		if err != nil {
			return nil, "", err
		}

		// If we got fewer results than the limit, there are no more pages
		if len(resp.Data) < limit {
			return resp.Data, "", nil
		}

		// Prepare for next page
		offset += len(resp.Data)
		return resp.Data, "next", nil
	})
}
