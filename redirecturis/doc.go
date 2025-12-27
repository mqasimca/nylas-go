// Package redirecturis provides types for the Nylas Redirect URIs API.
//
// The Redirect URIs API allows you to manage OAuth redirect URIs for your application,
// including platform-specific settings for web, iOS, Android, and JavaScript apps.
//
// Example:
//
//	// Create a redirect URI
//	uri, err := client.RedirectURIs.Create(ctx, &redirecturis.CreateRequest{
//	    URL:      "https://myapp.com/callback",
//	    Platform: "web",
//	})
//
//	// List all redirect URIs
//	resp, err := client.RedirectURIs.List(ctx, nil)
package redirecturis
