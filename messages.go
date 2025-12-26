package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/messages"
)

// List returns messages for a grant with optional filtering.
//
// The grantID is the ID of the connected account (obtained via OAuth).
// Use opts to filter messages by sender, recipient, subject, date range, etc.
// Returns a paginated response; use NextCursor for pagination or ListAll for iteration.
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

// Get returns a single message by ID.
//
// The grantID is the ID of the connected account.
// The messageID is the unique identifier of the message to retrieve.
// Returns ErrNotFound if the message does not exist.
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

// Send sends a new email message immediately.
//
// The grantID is the ID of the connected account to send from.
// At minimum, the To field and either Subject or Body must be provided.
// To schedule a message for later delivery, set SendAt to a future Unix timestamp.
// Returns the sent message with its assigned ID.
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

// Update modifies a message's metadata (read status, starred, folders).
//
// The grantID is the ID of the connected account.
// The messageID is the unique identifier of the message to update.
// Only the fields specified in the UpdateRequest are modified.
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

// Delete permanently removes a message.
//
// The grantID is the ID of the connected account.
// The messageID is the unique identifier of the message to delete.
// This action cannot be undone. Consider moving to trash instead.
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

// ListAll returns an iterator that automatically paginates through all messages.
//
// Use Next() to retrieve messages one at a time, or Collect() to get all at once.
// The iterator handles pagination automatically using the NextCursor from each response.
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

// ListScheduled returns all messages scheduled for future delivery.
//
// Scheduled messages were created with SendAt set to a future timestamp.
// Use StopScheduled to cancel a scheduled message before it's sent.
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

// GetScheduled returns details about a specific scheduled message.
//
// The scheduleID is obtained from ListScheduled or the response when scheduling a message.
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

// StopScheduled cancels a scheduled message before it's sent.
//
// The message will not be sent and cannot be recovered after cancellation.
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

// Clean extracts clean conversation text from messages.
//
// This removes quoted text, signatures, and other noise to get the core message content.
// Useful for AI processing, summarization, or displaying clean conversation threads.
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
