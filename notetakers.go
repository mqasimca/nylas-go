package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/notetakers"
)

// List returns all notetakers for a grant.
func (s *NotetakersService) List(ctx context.Context, grantID string, opts *notetakers.ListOptions) (*ListResponse[notetakers.Notetaker], error) {
	path := fmt.Sprintf("/v3/grants/%s/notetakers", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("notetakers.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []notetakers.Notetaker
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("notetakers.List: %w", err)
	}

	return &ListResponse[notetakers.Notetaker]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get retrieves a single notetaker by ID.
func (s *NotetakersService) Get(ctx context.Context, grantID, notetakerID string) (*notetakers.Notetaker, error) {
	path := fmt.Sprintf("/v3/grants/%s/notetakers/%s", grantID, notetakerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("notetakers.Get(%s): %w", notetakerID, err)
	}

	var nt notetakers.Notetaker
	_, err = s.client.Do(req, &nt)
	if err != nil {
		return nil, fmt.Errorf("notetakers.Get(%s): %w", notetakerID, err)
	}

	return &nt, nil
}

// Create invites a notetaker to join a meeting.
// If JoinTime is not specified, the notetaker will attempt to join immediately.
func (s *NotetakersService) Create(ctx context.Context, grantID string, createReq *notetakers.CreateRequest) (*notetakers.Notetaker, error) {
	path := fmt.Sprintf("/v3/grants/%s/notetakers", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, fmt.Errorf("notetakers.Create: %w", err)
	}

	var nt notetakers.Notetaker
	_, err = s.client.Do(req, &nt)
	if err != nil {
		return nil, fmt.Errorf("notetakers.Create: %w", err)
	}

	return &nt, nil
}

// Cancel cancels a scheduled notetaker before it joins the meeting.
func (s *NotetakersService) Cancel(ctx context.Context, grantID, notetakerID string) error {
	path := fmt.Sprintf("/v3/grants/%s/notetakers/%s/cancel", grantID, notetakerID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("notetakers.Cancel(%s): %w", notetakerID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("notetakers.Cancel(%s): %w", notetakerID, err)
	}

	return nil
}

// Leave removes a notetaker from an active meeting.
func (s *NotetakersService) Leave(ctx context.Context, grantID, notetakerID string) error {
	path := fmt.Sprintf("/v3/grants/%s/notetakers/%s/leave", grantID, notetakerID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("notetakers.Leave(%s): %w", notetakerID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("notetakers.Leave(%s): %w", notetakerID, err)
	}

	return nil
}

// GetHistory retrieves the event history for a notetaker.
// The history provides a timeline of everything that happened during the session.
func (s *NotetakersService) GetHistory(ctx context.Context, grantID, notetakerID string) (*notetakers.History, error) {
	path := fmt.Sprintf("/v3/grants/%s/notetakers/%s/history", grantID, notetakerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("notetakers.GetHistory(%s): %w", notetakerID, err)
	}

	var history notetakers.History
	_, err = s.client.Do(req, &history)
	if err != nil {
		return nil, fmt.Errorf("notetakers.GetHistory(%s): %w", notetakerID, err)
	}

	return &history, nil
}

// GetMedia retrieves media files (recordings, transcripts) for a notetaker.
func (s *NotetakersService) GetMedia(ctx context.Context, grantID, notetakerID string) ([]notetakers.Media, error) {
	path := fmt.Sprintf("/v3/grants/%s/notetakers/%s/media", grantID, notetakerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("notetakers.GetMedia(%s): %w", notetakerID, err)
	}

	var media []notetakers.Media
	_, err = s.client.Do(req, &media)
	if err != nil {
		return nil, fmt.Errorf("notetakers.GetMedia(%s): %w", notetakerID, err)
	}

	return media, nil
}

// ListAll returns an iterator for all notetakers.
func (s *NotetakersService) ListAll(ctx context.Context, grantID string, opts *notetakers.ListOptions) *Iterator[notetakers.Notetaker] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]notetakers.Notetaker, string, error) {
		o := opts
		if o == nil {
			o = &notetakers.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
