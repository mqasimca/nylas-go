package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/threads"
)

// List returns threads for a grant.
func (s *ThreadsService) List(ctx context.Context, grantID string, opts *threads.ListOptions) (*ListResponse[threads.Thread], error) {
	path := fmt.Sprintf("/v3/grants/%s/threads", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("threads.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []threads.Thread
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("threads.List: %w", err)
	}

	return &ListResponse[threads.Thread]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single thread.
func (s *ThreadsService) Get(ctx context.Context, grantID, threadID string) (*threads.Thread, error) {
	path := fmt.Sprintf("/v3/grants/%s/threads/%s", grantID, threadID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("threads.Get(%s): %w", threadID, err)
	}

	var thread threads.Thread
	_, err = s.client.Do(req, &thread)
	if err != nil {
		return nil, fmt.Errorf("threads.Get(%s): %w", threadID, err)
	}

	return &thread, nil
}

// Update updates a thread.
func (s *ThreadsService) Update(ctx context.Context, grantID, threadID string, update *threads.UpdateRequest) (*threads.Thread, error) {
	path := fmt.Sprintf("/v3/grants/%s/threads/%s", grantID, threadID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("threads.Update(%s): %w", threadID, err)
	}

	var thread threads.Thread
	_, err = s.client.Do(req, &thread)
	if err != nil {
		return nil, fmt.Errorf("threads.Update(%s): %w", threadID, err)
	}

	return &thread, nil
}

// Delete deletes a thread.
func (s *ThreadsService) Delete(ctx context.Context, grantID, threadID string) error {
	path := fmt.Sprintf("/v3/grants/%s/threads/%s", grantID, threadID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("threads.Delete(%s): %w", threadID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("threads.Delete(%s): %w", threadID, err)
	}

	return nil
}

// ListAll returns an iterator for all threads.
func (s *ThreadsService) ListAll(ctx context.Context, grantID string, opts *threads.ListOptions) *Iterator[threads.Thread] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]threads.Thread, string, error) {
		o := opts
		if o == nil {
			o = &threads.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
