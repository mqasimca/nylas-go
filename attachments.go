package nylas

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mqasimca/nylas-go/attachments"
)

// Get returns attachment metadata.
func (s *AttachmentsService) Get(ctx context.Context, grantID, attachmentID, messageID string) (*attachments.Attachment, error) {
	path := fmt.Sprintf("/v3/grants/%s/attachments/%s", grantID, attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("attachments.Get(%s): %w", attachmentID, err)
	}

	q := req.URL.Query()
	q.Set("message_id", messageID)
	req.URL.RawQuery = q.Encode()

	var attachment attachments.Attachment
	_, err = s.client.Do(req, &attachment)
	if err != nil {
		return nil, fmt.Errorf("attachments.Get(%s): %w", attachmentID, err)
	}

	return &attachment, nil
}

// Download downloads an attachment and returns the response.
func (s *AttachmentsService) Download(ctx context.Context, grantID, attachmentID, messageID string) (*attachments.DownloadResponse, error) {
	path := fmt.Sprintf("/v3/grants/%s/attachments/%s/download", grantID, attachmentID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("attachments.Download(%s): %w", attachmentID, err)
	}

	q := req.URL.Query()
	q.Set("message_id", messageID)
	req.URL.RawQuery = q.Encode()

	resp, err := s.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("attachments.Download(%s): %w", attachmentID, err)
	}

	if resp.StatusCode >= 400 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("attachments.Download(%s): status %d", attachmentID, resp.StatusCode)
	}

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	return &attachments.DownloadResponse{
		Content:     resp.Body,
		ContentType: resp.Header.Get("Content-Type"),
		Filename:    resp.Header.Get("Content-Disposition"),
		Size:        size,
	}, nil
}
