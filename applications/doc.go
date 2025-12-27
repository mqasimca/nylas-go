// Package applications provides types for the Nylas Applications API.
//
// The Applications API allows you to retrieve application configuration details
// including branding, environment, and hosted authentication settings.
//
// Example:
//
//	details, err := client.Applications.GetDetails(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("App ID: %s, Region: %s\n", details.ApplicationID, details.Region)
package applications
