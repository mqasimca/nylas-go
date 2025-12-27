// Package notetakers provides types for the Nylas Notetaker API.
//
// Nylas Notetaker is a real-time meeting bot that can be invited to
// online meetings to record and transcribe discussions. It supports
// Google Meet, Microsoft Teams, and Zoom sessions.
//
// Example usage:
//
//	// Invite a notetaker to join immediately
//	notetaker, err := client.Notetakers.Create(ctx, grantID, &notetakers.CreateRequest{
//		MeetingLink: "https://meet.google.com/abc-defg-hij",
//		Name:        "Meeting Bot",
//		MeetingSettings: &notetakers.MeetingSettings{
//			Transcription: true,
//			Summary:       true,
//			ActionItems:   true,
//		},
//	})
//
//	// Schedule a notetaker to join at a specific time
//	joinTime := time.Now().Add(1 * time.Hour).Unix()
//	scheduled, err := client.Notetakers.Create(ctx, grantID, &notetakers.CreateRequest{
//		MeetingLink: "https://zoom.us/j/123456789",
//		JoinTime:    &joinTime,
//	})
//
//	// List notetakers
//	list, err := client.Notetakers.List(ctx, grantID, nil)
//
//	// Get notetaker history
//	history, err := client.Notetakers.GetHistory(ctx, grantID, notetakerID)
package notetakers
