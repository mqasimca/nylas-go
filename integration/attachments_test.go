//go:build integration

// Attachments Integration Tests Coverage:
//   - Get, Download ✓
//
// All AttachmentsService methods are fully tested.

package integration

import (
	"io"
	"testing"

	"github.com/mqasimca/nylas-go/messages"
)

func TestAttachments_Get(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Find a message with attachments
		hasAttachment := true
		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			HasAttachment: &hasAttachment,
			Limit:         intPtr(5),
		})
		if err != nil {
			t.Fatalf("Messages.List() error = %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No messages with attachments found")
		}

		// Find a message with actual attachment metadata
		var messageID, attachmentID string
		for _, msg := range resp.Data {
			if len(msg.Attachments) > 0 {
				messageID = msg.ID
				attachmentID = msg.Attachments[0].ID
				break
			}
		}

		if attachmentID == "" {
			t.Skip("No attachment metadata found in messages")
		}

		t.Logf("Testing with message %s, attachment %s", messageID, attachmentID)

		// Get attachment metadata
		attachment, err := client.Attachments.Get(ctx, grantID, attachmentID, messageID)
		if err != nil {
			t.Fatalf("Get(%s) error = %v", attachmentID, err)
		}

		if attachment.ID != attachmentID {
			t.Errorf("Get() ID = %s, want %s", attachment.ID, attachmentID)
		}

		t.Logf("Got attachment: %s (type: %s, size: %d bytes)",
			attachment.Filename, attachment.ContentType, attachment.Size)
	})
}

func TestAttachments_Download(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// Find a message with attachments
		hasAttachment := true
		resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
			HasAttachment: &hasAttachment,
			Limit:         intPtr(5),
		})
		if err != nil {
			t.Fatalf("Messages.List() error = %v", err)
		}

		if len(resp.Data) == 0 {
			t.Skip("No messages with attachments found")
		}

		// Find a message with actual attachment metadata (prefer small files)
		var messageID, attachmentID string
		var attachmentSize int
		for _, msg := range resp.Data {
			for _, att := range msg.Attachments {
				// Skip very large attachments to keep tests fast
				if att.Size > 0 && att.Size < 1024*1024 { // < 1MB
					messageID = msg.ID
					attachmentID = att.ID
					attachmentSize = att.Size
					break
				}
			}
			if attachmentID != "" {
				break
			}
		}

		if attachmentID == "" {
			// Fall back to first attachment regardless of size
			for _, msg := range resp.Data {
				if len(msg.Attachments) > 0 {
					messageID = msg.ID
					attachmentID = msg.Attachments[0].ID
					attachmentSize = msg.Attachments[0].Size
					break
				}
			}
		}

		if attachmentID == "" {
			t.Skip("No attachment metadata found in messages")
		}

		t.Logf("Downloading attachment %s from message %s (size: %d bytes)",
			attachmentID, messageID, attachmentSize)

		// Download attachment
		download, err := client.Attachments.Download(ctx, grantID, attachmentID, messageID)
		if err != nil {
			t.Fatalf("Download(%s) error = %v", attachmentID, err)
		}
		defer download.Content.Close()

		// Read content to verify download works
		data, err := io.ReadAll(download.Content)
		if err != nil {
			t.Fatalf("ReadAll() error = %v", err)
		}

		if len(data) == 0 {
			t.Error("Download() returned empty content")
		}

		t.Logf("Downloaded %d bytes (content-type: %s)", len(data), download.ContentType)
	})
}
