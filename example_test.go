package nylas_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/mqasimca/nylas-go"
	"github.com/mqasimca/nylas-go/calendars"
	"github.com/mqasimca/nylas-go/events"
	"github.com/mqasimca/nylas-go/messages"
)

func ExampleNewClient() {
	// Create a client with your API key
	client, err := nylas.NewClient(
		nylas.WithAPIKey("your-api-key"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use the client to access Nylas services
	_ = client.Messages
	_ = client.Calendars
	_ = client.Events
}

func ExampleNewClient_withOptions() {
	// Create a client with multiple options
	client, err := nylas.NewClient(
		nylas.WithAPIKey("your-api-key"),
		nylas.WithRegion(nylas.RegionEU), // Use EU region
		nylas.WithMaxRetries(3),          // Retry failed requests
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}

func ExampleMessagesService_List() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()
	grantID := "your-grant-id"

	// List messages with filtering
	resp, err := client.Messages.List(ctx, grantID, &messages.ListOptions{
		Limit:  nylas.Ptr(10),
		Unread: nylas.Ptr(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range resp.Data {
		fmt.Printf("Subject: %s\n", msg.Subject)
	}
}

func ExampleMessagesService_Send() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()
	grantID := "your-grant-id"

	// Send an email
	msg, err := client.Messages.Send(ctx, grantID, &messages.SendRequest{
		To:      []messages.Participant{{Email: "recipient@example.com"}},
		Subject: "Hello from Nylas",
		Body:    "<h1>Hello!</h1><p>This is a test email.</p>",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sent message ID: %s\n", msg.ID)
}

func ExampleIterator() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()
	grantID := "your-grant-id"

	// Iterate through all messages
	iter := client.Messages.ListAll(ctx, grantID, nil)
	for {
		msg, err := iter.Next()
		if errors.Is(err, nylas.ErrDone) {
			break // No more messages
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Message: %s\n", msg.Subject)
	}
}

func ExampleIterator_Collect() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()
	grantID := "your-grant-id"

	// Collect all messages at once
	iter := client.Messages.ListAll(ctx, grantID, &messages.ListOptions{
		Limit: nylas.Ptr(100),
	})

	allMessages, err := iter.Collect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d messages\n", len(allMessages))
}

func ExampleEventsService_Create() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()
	grantID := "your-grant-id"
	calendarID := "primary"

	// Create a calendar event
	event, err := client.Events.Create(ctx, grantID, calendarID, &events.CreateRequest{
		Title:       "Team Meeting",
		Description: "Weekly sync with the team",
		When: events.When{
			StartTime: nylas.Ptr(int64(1704067200)), // Unix timestamp
			EndTime:   nylas.Ptr(int64(1704070800)),
		},
		Participants: []events.Participant{
			{Email: "colleague@example.com", Name: "Colleague"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created event: %s\n", event.ID)
}

func ExampleCalendarsService_Availability() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()

	// Check availability for multiple participants
	availability, err := client.Calendars.Availability(ctx, &calendars.AvailabilityRequest{
		StartTime:       1704067200,
		EndTime:         1704153600,
		DurationMinutes: 30,
		Participants: []calendars.AvailabilityParticipant{
			{Email: "user1@example.com"},
			{Email: "user2@example.com"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, slot := range availability.TimeSlots {
		fmt.Printf("Available: %d - %d\n", slot.StartTime, slot.EndTime)
	}
}

func Example_errorHandling() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))
	ctx := context.Background()

	_, err := client.Messages.Get(ctx, "grant-id", "invalid-message-id")
	if err != nil {
		// Check for specific error types
		if errors.Is(err, nylas.ErrNotFound) {
			fmt.Println("Message not found")
			return
		}
		if errors.Is(err, nylas.ErrUnauthorized) {
			fmt.Println("Invalid API key")
			return
		}
		if errors.Is(err, nylas.ErrRateLimited) {
			fmt.Println("Rate limited - try again later")
			return
		}

		// Get detailed error information
		var apiErr *nylas.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API Error: %s (status: %d, request_id: %s)\n",
				apiErr.Message, apiErr.StatusCode, apiErr.RequestID)
		}
	}
}

func Example_rateLimits() {
	client, _ := nylas.NewClient(nylas.WithAPIKey("your-api-key"))

	// Check current rate limit status
	limits := client.RateLimits()
	fmt.Printf("Remaining: %d, Reset: %v\n", limits.Remaining, limits.Reset)
}
