// Package smartcompose provides types for the Nylas Smart Compose API.
//
// The Smart Compose API provides AI-powered message composition capabilities,
// allowing you to generate message suggestions based on prompts or create
// replies to existing messages.
//
// Example:
//
//	// Generate a new message
//	resp, err := client.SmartCompose.ComposeMessage(ctx, grantID, &smartcompose.ComposeRequest{
//	    Prompt: "Write a professional email declining a meeting invitation",
//	})
//	fmt.Println(resp.Suggestion)
//
//	// Generate a reply to an existing message
//	resp, err := client.SmartCompose.ComposeReply(ctx, grantID, messageID, &smartcompose.ComposeRequest{
//	    Prompt: "Accept the proposal and suggest a follow-up call",
//	})
package smartcompose
