package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/calendars"
)

// List returns calendars for a grant.
func (s *CalendarsService) List(ctx context.Context, grantID string, opts *calendars.ListOptions) (*ListResponse[calendars.Calendar], error) {
	path := fmt.Sprintf("/v3/grants/%s/calendars", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("calendars.List: %w", err)
	}

	q := req.URL.Query()
	if opts != nil {
		setQueryParams(q, opts.Values())
	}
	// Set default limit if not provided (matches API default of 50)
	if q.Get("limit") == "" {
		q.Set("limit", "50")
	}
	req.URL.RawQuery = q.Encode()

	var data []calendars.Calendar
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("calendars.List: %w", err)
	}

	return &ListResponse[calendars.Calendar]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single calendar.
func (s *CalendarsService) Get(ctx context.Context, grantID, calendarID string) (*calendars.Calendar, error) {
	path := fmt.Sprintf("/v3/grants/%s/calendars/%s", grantID, calendarID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("calendars.Get(%s): %w", calendarID, err)
	}

	var cal calendars.Calendar
	_, err = s.client.Do(req, &cal)
	if err != nil {
		return nil, fmt.Errorf("calendars.Get(%s): %w", calendarID, err)
	}

	return &cal, nil
}

// Create creates a new calendar.
func (s *CalendarsService) Create(ctx context.Context, grantID string, create *calendars.CreateRequest) (*calendars.Calendar, error) {
	path := fmt.Sprintf("/v3/grants/%s/calendars", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("calendars.Create: %w", err)
	}

	var cal calendars.Calendar
	_, err = s.client.Do(req, &cal)
	if err != nil {
		return nil, fmt.Errorf("calendars.Create: %w", err)
	}

	return &cal, nil
}

// Update updates a calendar.
func (s *CalendarsService) Update(ctx context.Context, grantID, calendarID string, update *calendars.UpdateRequest) (*calendars.Calendar, error) {
	path := fmt.Sprintf("/v3/grants/%s/calendars/%s", grantID, calendarID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("calendars.Update(%s): %w", calendarID, err)
	}

	var cal calendars.Calendar
	_, err = s.client.Do(req, &cal)
	if err != nil {
		return nil, fmt.Errorf("calendars.Update(%s): %w", calendarID, err)
	}

	return &cal, nil
}

// Delete deletes a calendar.
func (s *CalendarsService) Delete(ctx context.Context, grantID, calendarID string) error {
	path := fmt.Sprintf("/v3/grants/%s/calendars/%s", grantID, calendarID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("calendars.Delete(%s): %w", calendarID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("calendars.Delete(%s): %w", calendarID, err)
	}

	return nil
}

// ListAll returns an iterator for all calendars.
func (s *CalendarsService) ListAll(ctx context.Context, grantID string, opts *calendars.ListOptions) *Iterator[calendars.Calendar] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]calendars.Calendar, string, error) {
		o := opts
		if o == nil {
			o = &calendars.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}

// Availability checks availability for participants.
func (s *CalendarsService) Availability(ctx context.Context, avail *calendars.AvailabilityRequest) (*calendars.AvailabilityResponse, error) {
	path := "/v3/calendars/availability"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, avail)
	if err != nil {
		return nil, fmt.Errorf("calendars.Availability: %w", err)
	}

	var result calendars.AvailabilityResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("calendars.Availability: %w", err)
	}

	return &result, nil
}

// FreeBusy returns free/busy information for emails.
func (s *CalendarsService) FreeBusy(ctx context.Context, grantID string, freeBusy *calendars.FreeBusyRequest) ([]calendars.FreeBusyResponse, error) {
	path := fmt.Sprintf("/v3/grants/%s/calendars/free-busy", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, freeBusy)
	if err != nil {
		return nil, fmt.Errorf("calendars.FreeBusy: %w", err)
	}

	var result []calendars.FreeBusyResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("calendars.FreeBusy: %w", err)
	}

	return result, nil
}
