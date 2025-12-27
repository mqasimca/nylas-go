// Package credentials provides types for the Nylas Credentials API.
//
// The Credentials API allows you to manage connector credentials,
// which store authentication information for connecting to email providers
// via service accounts or admin consent.
//
// Example:
//
//	// Create a Microsoft admin consent credential
//	cred, err := client.Credentials.Create(ctx, connectors.ProviderMicrosoft,
//	    credentials.CreateMicrosoftRequest(
//	        "my-ms-credential",
//	        "client-id",
//	        "client-secret",
//	    ))
//
//	// List credentials for a provider
//	resp, err := client.Credentials.List(ctx, connectors.ProviderGoogle, nil)
package credentials
