package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/events"
)

// List returns events for a grant.
func (s *EventsService) List(ctx context.Context, grantID string, opts *events.ListOptions) (*ListResponse[events.Event], error) {
	path := fmt.Sprintf("/v3/grants/%s/events", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("events.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []events.Event
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("events.List: %w", err)
	}

	return &ListResponse[events.Event]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single event.
func (s *EventsService) Get(ctx context.Context, grantID, eventID string, calendarID string) (*events.Event, error) {
	path := fmt.Sprintf("/v3/grants/%s/events/%s", grantID, eventID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("events.Get(%s): %w", eventID, err)
	}

	if calendarID != "" {
		q := req.URL.Query()
		q.Set("calendar_id", calendarID)
		req.URL.RawQuery = q.Encode()
	}

	var event events.Event
	_, err = s.client.Do(req, &event)
	if err != nil {
		return nil, fmt.Errorf("events.Get(%s): %w", eventID, err)
	}

	return &event, nil
}

// Create creates a new event.
func (s *EventsService) Create(ctx context.Context, grantID, calendarID string, create *events.CreateRequest) (*events.Event, error) {
	path := fmt.Sprintf("/v3/grants/%s/events", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, create)
	if err != nil {
		return nil, fmt.Errorf("events.Create: %w", err)
	}

	q := req.URL.Query()
	q.Set("calendar_id", calendarID)
	req.URL.RawQuery = q.Encode()

	var event events.Event
	_, err = s.client.Do(req, &event)
	if err != nil {
		return nil, fmt.Errorf("events.Create: %w", err)
	}

	return &event, nil
}

// Update updates an event.
func (s *EventsService) Update(ctx context.Context, grantID, eventID, calendarID string, update *events.UpdateRequest) (*events.Event, error) {
	path := fmt.Sprintf("/v3/grants/%s/events/%s", grantID, eventID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("events.Update(%s): %w", eventID, err)
	}

	q := req.URL.Query()
	q.Set("calendar_id", calendarID)
	req.URL.RawQuery = q.Encode()

	var event events.Event
	_, err = s.client.Do(req, &event)
	if err != nil {
		return nil, fmt.Errorf("events.Update(%s): %w", eventID, err)
	}

	return &event, nil
}

// Delete deletes an event.
func (s *EventsService) Delete(ctx context.Context, grantID, eventID, calendarID string) error {
	path := fmt.Sprintf("/v3/grants/%s/events/%s", grantID, eventID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("events.Delete(%s): %w", eventID, err)
	}

	q := req.URL.Query()
	q.Set("calendar_id", calendarID)
	req.URL.RawQuery = q.Encode()

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("events.Delete(%s): %w", eventID, err)
	}

	return nil
}

// ListAll returns an iterator for all events.
func (s *EventsService) ListAll(ctx context.Context, grantID string, opts *events.ListOptions) *Iterator[events.Event] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]events.Event, string, error) {
		o := opts
		if o == nil {
			o = &events.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}

// SendRSVP sends an RSVP response for an event.
func (s *EventsService) SendRSVP(ctx context.Context, grantID, eventID, calendarID string, rsvp *events.RSVPRequest) error {
	path := fmt.Sprintf("/v3/grants/%s/events/%s/send-rsvp", grantID, eventID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rsvp)
	if err != nil {
		return fmt.Errorf("events.SendRSVP(%s): %w", eventID, err)
	}

	q := req.URL.Query()
	q.Set("calendar_id", calendarID)
	req.URL.RawQuery = q.Encode()

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("events.SendRSVP(%s): %w", eventID, err)
	}

	return nil
}
