//go:build integration

// SmartCompose Integration Tests Coverage:
//   - ComposeMessage ✓
//   - ComposeReply ✓
//
// These are read-only AI generation endpoints that don't modify any data.

package integration

import (
	"testing"

	"github.com/mqasimca/nylas-go/smartcompose"
)

func TestSmartCompose_ComposeMessage(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.SmartCompose.ComposeMessage(ctx, grantID, &smartcompose.ComposeRequest{
			Prompt: "Write a professional email introducing myself as a software developer looking for new opportunities.",
		})
		if err != nil {
			// SmartCompose may not be enabled for all accounts
			t.Logf("ComposeMessage() error = %v (SmartCompose may not be enabled)", err)
			t.Skip("SmartCompose not available for this grant")
		}

		if resp.Suggestion == "" {
			t.Error("Suggestion should not be empty")
		}

		t.Logf("Generated suggestion (%d chars):", len(resp.Suggestion))
		// Truncate for logging
		if len(resp.Suggestion) > 200 {
			t.Logf("  %s...", resp.Suggestion[:200])
		} else {
			t.Logf("  %s", resp.Suggestion)
		}
	})
}

func TestSmartCompose_ComposeMessage_ShortPrompt(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		resp, err := client.SmartCompose.ComposeMessage(ctx, grantID, &smartcompose.ComposeRequest{
			Prompt: "Thank you for your time",
		})
		if err != nil {
			t.Logf("ComposeMessage() error = %v (SmartCompose may not be enabled)", err)
			t.Skip("SmartCompose not available for this grant")
		}

		if resp.Suggestion == "" {
			t.Error("Suggestion should not be empty")
		}

		t.Logf("Generated short suggestion: %s", resp.Suggestion)
	})
}

func TestSmartCompose_ComposeReply(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First, list messages to find one to reply to
		msgResp, err := client.Messages.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Messages.List() error = %v", err)
		}

		if len(msgResp.Data) == 0 {
			t.Skip("No messages found to reply to")
		}

		// Use the first message
		messageID := msgResp.Data[0].ID
		t.Logf("Composing reply to message: %s", messageID)

		resp, err := client.SmartCompose.ComposeReply(ctx, grantID, messageID, &smartcompose.ComposeRequest{
			Prompt: "Write a polite response acknowledging receipt and promising to review",
		})
		if err != nil {
			t.Logf("ComposeReply() error = %v (SmartCompose may not be enabled)", err)
			t.Skip("SmartCompose not available for this grant")
		}

		if resp.Suggestion == "" {
			t.Error("Suggestion should not be empty")
		}

		t.Logf("Generated reply suggestion (%d chars):", len(resp.Suggestion))
		// Truncate for logging
		if len(resp.Suggestion) > 200 {
			t.Logf("  %s...", resp.Suggestion[:200])
		} else {
			t.Logf("  %s", resp.Suggestion)
		}
	})
}

func TestSmartCompose_ComposeReply_Decline(t *testing.T) {
	cfg := LoadConfig(t)
	client := NewTestClient(t, cfg)
	ctx := NewTestContext(t)

	RunForEachProvider(t, cfg, func(t *testing.T, grantID string) {
		// First, list messages to find one to reply to
		msgResp, err := client.Messages.List(ctx, grantID, nil)
		if err != nil {
			t.Fatalf("Messages.List() error = %v", err)
		}

		if len(msgResp.Data) == 0 {
			t.Skip("No messages found to reply to")
		}

		messageID := msgResp.Data[0].ID

		resp, err := client.SmartCompose.ComposeReply(ctx, grantID, messageID, &smartcompose.ComposeRequest{
			Prompt: "Politely decline the request",
		})
		if err != nil {
			t.Logf("ComposeReply() error = %v (SmartCompose may not be enabled)", err)
			t.Skip("SmartCompose not available for this grant")
		}

		if resp.Suggestion == "" {
			t.Error("Suggestion should not be empty")
		}

		t.Logf("Generated decline reply: %s", resp.Suggestion)
	})
}
