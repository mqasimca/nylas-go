// Package connectors provides types for the Nylas Connectors API.
//
// The Connectors API allows you to manage email provider connectors,
// which define how your application connects to providers like Google,
// Microsoft, and IMAP.
//
// Example:
//
//	// Create a Google connector
//	connector, err := client.Connectors.Create(ctx, connectors.CreateGoogleRequest(
//	    "my-google-connector",
//	    "google-client-id",
//	    "google-client-secret",
//	    []string{"email", "calendar"},
//	))
//
//	// List all connectors
//	resp, err := client.Connectors.List(ctx, nil)
package connectors
