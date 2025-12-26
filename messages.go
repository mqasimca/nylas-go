package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/messages"
)

// List returns messages for a grant.
func (s *MessagesService) List(ctx context.Context, grantID string, opts *messages.ListOptions) (*ListResponse[messages.Message], error) {
	path := fmt.Sprintf("/v3/grants/%s/messages", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("messages.List: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []messages.Message
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("messages.List: %w", err)
	}

	return &ListResponse[messages.Message]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// Get returns a single message.
func (s *MessagesService) Get(ctx context.Context, grantID, messageID string) (*messages.Message, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/%s", grantID, messageID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("messages.Get(%s): %w", messageID, err)
	}

	var msg messages.Message
	_, err = s.client.Do(req, &msg)
	if err != nil {
		return nil, fmt.Errorf("messages.Get(%s): %w", messageID, err)
	}

	return &msg, nil
}

// Send sends a new message.
func (s *MessagesService) Send(ctx context.Context, grantID string, send *messages.SendRequest) (*messages.Message, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/send", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, send)
	if err != nil {
		return nil, fmt.Errorf("messages.Send: %w", err)
	}

	var msg messages.Message
	_, err = s.client.Do(req, &msg)
	if err != nil {
		return nil, fmt.Errorf("messages.Send: %w", err)
	}

	return &msg, nil
}

// Update updates a message.
func (s *MessagesService) Update(ctx context.Context, grantID, messageID string, update *messages.UpdateRequest) (*messages.Message, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/%s", grantID, messageID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, update)
	if err != nil {
		return nil, fmt.Errorf("messages.Update(%s): %w", messageID, err)
	}

	var msg messages.Message
	_, err = s.client.Do(req, &msg)
	if err != nil {
		return nil, fmt.Errorf("messages.Update(%s): %w", messageID, err)
	}

	return &msg, nil
}

// Delete deletes a message.
func (s *MessagesService) Delete(ctx context.Context, grantID, messageID string) error {
	path := fmt.Sprintf("/v3/grants/%s/messages/%s", grantID, messageID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("messages.Delete(%s): %w", messageID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("messages.Delete(%s): %w", messageID, err)
	}

	return nil
}

// ListAll returns an iterator for all messages.
func (s *MessagesService) ListAll(ctx context.Context, grantID string, opts *messages.ListOptions) *Iterator[messages.Message] {
	return NewIterator(ctx, func(ctx context.Context, pageToken string) ([]messages.Message, string, error) {
		o := opts
		if o == nil {
			o = &messages.ListOptions{}
		}
		o.PageToken = pageToken

		resp, err := s.List(ctx, grantID, o)
		if err != nil {
			return nil, "", err
		}
		return resp.Data, resp.NextCursor, nil
	})
}

// ListScheduled returns scheduled messages.
func (s *MessagesService) ListScheduled(ctx context.Context, grantID string) (messages.ScheduledMessagesList, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/schedules", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("messages.ListScheduled: %w", err)
	}

	var result messages.ScheduledMessagesList
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("messages.ListScheduled: %w", err)
	}

	return result, nil
}

// GetScheduled returns a scheduled message.
func (s *MessagesService) GetScheduled(ctx context.Context, grantID, scheduleID string) (*messages.ScheduledMessage, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/schedules/%s", grantID, scheduleID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("messages.GetScheduled(%s): %w", scheduleID, err)
	}

	var msg messages.ScheduledMessage
	_, err = s.client.Do(req, &msg)
	if err != nil {
		return nil, fmt.Errorf("messages.GetScheduled(%s): %w", scheduleID, err)
	}

	return &msg, nil
}

// StopScheduled stops a scheduled message.
func (s *MessagesService) StopScheduled(ctx context.Context, grantID, scheduleID string) error {
	path := fmt.Sprintf("/v3/grants/%s/messages/schedules/%s", grantID, scheduleID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("messages.StopScheduled(%s): %w", scheduleID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("messages.StopScheduled(%s): %w", scheduleID, err)
	}

	return nil
}

// Clean removes extra information from messages.
func (s *MessagesService) Clean(ctx context.Context, grantID string, clean *messages.CleanRequest) ([]messages.CleanResponse, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/clean", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, clean)
	if err != nil {
		return nil, fmt.Errorf("messages.Clean: %w", err)
	}

	var result []messages.CleanResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("messages.Clean: %w", err)
	}

	return result, nil
}
