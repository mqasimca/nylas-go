package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/smartcompose"
)

// ComposeMessage generates a message suggestion based on a prompt.
func (s *SmartComposeService) ComposeMessage(ctx context.Context, grantID string, compose *smartcompose.ComposeRequest) (*smartcompose.ComposeResponse, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/smart-compose", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, compose)
	if err != nil {
		return nil, fmt.Errorf("smartcompose.ComposeMessage: %w", err)
	}

	var result smartcompose.ComposeResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("smartcompose.ComposeMessage: %w", err)
	}

	return &result, nil
}

// ComposeReply generates a reply suggestion for an existing message.
func (s *SmartComposeService) ComposeReply(ctx context.Context, grantID, messageID string, compose *smartcompose.ComposeRequest) (*smartcompose.ComposeResponse, error) {
	path := fmt.Sprintf("/v3/grants/%s/messages/%s/smart-compose", grantID, messageID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, compose)
	if err != nil {
		return nil, fmt.Errorf("smartcompose.ComposeReply(%s): %w", messageID, err)
	}

	var result smartcompose.ComposeResponse
	_, err = s.client.Do(req, &result)
	if err != nil {
		return nil, fmt.Errorf("smartcompose.ComposeReply(%s): %w", messageID, err)
	}

	return &result, nil
}
