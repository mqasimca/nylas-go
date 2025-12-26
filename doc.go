// Package nylas provides a Go client for the Nylas API v3.
//
// The client supports all major Nylas API resources including Messages, Threads,
// Drafts, Calendars, Events, Contacts, Folders, Attachments, Grants, and Webhooks.
//
// # Creating a Client
//
// Create a client with your API key:
//
//	client, err := nylas.NewClient(
//	    nylas.WithAPIKey("your-api-key"),
//	)
//
// For EU region:
//
//	client, err := nylas.NewClient(
//	    nylas.WithAPIKey("your-api-key"),
//	    nylas.WithRegion(nylas.RegionEU),
//	)
//
// # Using Services
//
// Access Nylas resources through service methods:
//
//	// List messages
//	resp, err := client.Messages.List(ctx, grantID, nil)
//
//	// Get a single message
//	msg, err := client.Messages.Get(ctx, grantID, messageID)
//
//	// Send a message
//	msg, err := client.Messages.Send(ctx, grantID, &messages.SendRequest{
//	    To:      []messages.Participant{{Email: "recipient@example.com"}},
//	    Subject: "Hello",
//	    Body:    "World",
//	})
//
// # Pagination
//
// Use iterators for paginated results:
//
//	iter := client.Messages.ListAll(ctx, grantID, nil)
//	for {
//	    msg, err := iter.Next()
//	    if errors.Is(err, nylas.ErrDone) {
//	        break
//	    }
//	    if err != nil {
//	        return err
//	    }
//	    process(msg)
//	}
//
// Or collect all results at once:
//
//	all, err := iter.Collect()
//
// # Error Handling
//
// Check for specific error conditions:
//
//	if errors.Is(err, nylas.ErrNotFound) {
//	    // Resource not found
//	}
//	if errors.Is(err, nylas.ErrRateLimited) {
//	    // Rate limited, retry later
//	}
//
// # Rate Limiting
//
// The client automatically handles rate limiting with exponential backoff.
// You can also check current rate limit status:
//
//	limits := client.RateLimits()
//	fmt.Printf("Remaining: %d, Reset: %v\n", limits.Remaining, limits.Reset)
package nylas
