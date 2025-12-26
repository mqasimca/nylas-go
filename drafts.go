package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/drafts"
)

// List returns drafts for a grant.
func (s *DraftsService) List(ctx context.Context, grantID string, opts *drafts.ListOptions) (*ListResponse[drafts.Draft], error) {
	path := fmt.Sprintf("/v3/grants/%s/drafts", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("drafts.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []drafts.Draft
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("drafts.List: %w", err)
	}

	return &ListResponse[drafts.Draft]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single draft.
func (s *DraftsService) Get(ctx context.Context, grantID, draftID string) (*drafts.Draft, error) {
	path := fmt.Sprintf("/v3/grants/%s/drafts/%s", grantID, draftID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("drafts.Get(%s): %w", draftID, err)
	}

	var draft drafts.Draft
	_, err = s.client.Do(req, &draft)
	if err != nil {
		return nil, fmt.Errorf("drafts.Get(%s): %w", draftID, err)
	}

	return &draft, nil
}

// Create creates a new draft.
func (s *DraftsService) Create(ctx context.Context, grantID string, create *drafts.CreateRequest) (*drafts.Draft, error) {
	path := fmt.Sprintf("/v3/grants/%s/drafts", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("drafts.Create: %w", err)
	}

	var draft drafts.Draft
	_, err = s.client.Do(req, &draft)
	if err != nil {
		return nil, fmt.Errorf("drafts.Create: %w", err)
	}

	return &draft, nil
}

// Update updates a draft.
func (s *DraftsService) Update(ctx context.Context, grantID, draftID string, update *drafts.UpdateRequest) (*drafts.Draft, error) {
	path := fmt.Sprintf("/v3/grants/%s/drafts/%s", grantID, draftID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("drafts.Update(%s): %w", draftID, err)
	}

	var draft drafts.Draft
	_, err = s.client.Do(req, &draft)
	if err != nil {
		return nil, fmt.Errorf("drafts.Update(%s): %w", draftID, err)
	}

	return &draft, nil
}

// Delete deletes a draft.
func (s *DraftsService) Delete(ctx context.Context, grantID, draftID string) error {
	path := fmt.Sprintf("/v3/grants/%s/drafts/%s", grantID, draftID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("drafts.Delete(%s): %w", draftID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("drafts.Delete(%s): %w", draftID, err)
	}

	return nil
}

// Send sends a draft as an email message.
func (s *DraftsService) Send(ctx context.Context, grantID, draftID string) (*drafts.Draft, error) {
	path := fmt.Sprintf("/v3/grants/%s/drafts/%s", grantID, draftID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("drafts.Send(%s): %w", draftID, err)
	}

	var draft drafts.Draft
	_, err = s.client.Do(req, &draft)
	if err != nil {
		return nil, fmt.Errorf("drafts.Send(%s): %w", draftID, err)
	}

	return &draft, nil
}

// ListAll returns an iterator for all drafts.
func (s *DraftsService) ListAll(ctx context.Context, grantID string, opts *drafts.ListOptions) *Iterator[drafts.Draft] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]drafts.Draft, string, error) {
		o := opts
		if o == nil {
			o = &drafts.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
