package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/events"
)

// List returns events for a grant.
//
// Use ListOptions to filter by calendar, time range, or other criteria.
//
// Example:
//
//	// List events for the next 7 days
//	now := time.Now()
//	resp, err := client.Events.List(ctx, grantID, &events.ListOptions{
//	    CalendarID: "primary",
//	    Start:      nylas.Ptr(now.Unix()),
//	    End:        nylas.Ptr(now.Add(7 * 24 * time.Hour).Unix()),
//	    Limit:      nylas.Ptr(50),
//	})
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
//
// Example - simple event:
//
//	event, err := client.Events.Create(ctx, grantID, "primary", &events.CreateRequest{
//	    Title: "Team Meeting",
//	    When: events.When{
//	        StartTime: nylas.Ptr(time.Now().Add(24 * time.Hour).Unix()),
//	        EndTime:   nylas.Ptr(time.Now().Add(25 * time.Hour).Unix()),
//	    },
//	})
//
// Example - event with Google Meet and participants:
//
//	event, err := client.Events.Create(ctx, grantID, "primary", &events.CreateRequest{
//	    Title:       "Project Sync",
//	    Description: "Weekly project status update",
//	    When: events.When{
//	        StartTime: nylas.Ptr(startTime.Unix()),
//	        EndTime:   nylas.Ptr(endTime.Unix()),
//	    },
//	    Participants: []events.Participant{
//	        {Email: "alice@example.com", Name: "Alice"},
//	        {Email: "bob@example.com", Name: "Bob"},
//	    },
//	    Conferencing: &events.Conferencing{
//	        Provider: "Google Meet",
//	        Autocreate: &events.AutocreateConfig{},
//	    },
//	})
//
// Example - recurring event (every weekday):
//
//	event, err := client.Events.Create(ctx, grantID, "primary", &events.CreateRequest{
//	    Title: "Daily Standup",
//	    When: events.When{
//	        StartTime: nylas.Ptr(startTime.Unix()),
//	        EndTime:   nylas.Ptr(endTime.Unix()),
//	    },
//	    Recurrence: &events.Recurrence{
//	        RRule: []string{"RRULE:FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR"},
//	    },
//	})
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
//
// Only fields that are set in the UpdateRequest will be modified.
//
// Example - reschedule an event:
//
//	event, err := client.Events.Update(ctx, grantID, eventID, calendarID, &events.UpdateRequest{
//	    When: &events.When{
//	        StartTime: nylas.Ptr(newStartTime.Unix()),
//	        EndTime:   nylas.Ptr(newEndTime.Unix()),
//	    },
//	})
//
// Example - update title and add a participant:
//
//	event, err := client.Events.Update(ctx, grantID, eventID, calendarID, &events.UpdateRequest{
//	    Title: nylas.Ptr("Updated Meeting Title"),
//	    Participants: []events.Participant{
//	        {Email: "newattendee@example.com", Name: "New Attendee"},
//	    },
//	})
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
//
// Use this to accept, decline, or tentatively accept an event invitation.
//
// Example:
//
//	// Accept an invitation
//	err := client.Events.SendRSVP(ctx, grantID, eventID, calendarID, &events.RSVPRequest{
//	    Status: "yes",
//	})
//
//	// Decline with a message
//	err := client.Events.SendRSVP(ctx, grantID, eventID, calendarID, &events.RSVPRequest{
//	    Status:  "no",
//	    Comment: "I have a conflict at this time.",
//	})
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

// Import returns events from a calendar including recurring event instances
// with their parent events and any overrides. This is useful for working with
// recurring events as it returns the complete recurrence information.
func (s *EventsService) Import(ctx context.Context, grantID string, opts *events.ImportOptions) (*ListResponse[events.Event], error) {
	path := fmt.Sprintf("/v3/grants/%s/events/import", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("events.Import: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []events.Event
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("events.Import: %w", err)
	}

	return &ListResponse[events.Event]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// ImportAll returns an iterator for importing all events from a calendar.
func (s *EventsService) ImportAll(ctx context.Context, grantID string, opts *events.ImportOptions) *Iterator[events.Event] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]events.Event, string, error) {
		o := opts
		if o == nil {
			o = &events.ImportOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.Import(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}
