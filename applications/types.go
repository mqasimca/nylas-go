// Package applications provides types for the Nylas Applications API.
package applications

// ApplicationDetails represents application configuration details.
type ApplicationDetails struct {
	ApplicationID        string                `json:"application_id"`
	OrganizationID       string                `json:"organization_id"`
	Region               string                `json:"region"`      // "us" or "eu"
	Environment          string                `json:"environment"` // "production" or "staging"
	Branding             *Branding             `json:"branding,omitempty"`
	HostedAuthentication *HostedAuthentication `json:"hosted_authentication,omitempty"`
	CallbackURIs         []CallbackURI         `json:"callback_uris,omitempty"`
}

// CallbackURI represents a callback/redirect URI configuration.
type CallbackURI struct {
	ID       string            `json:"id"`
	URL      string            `json:"url"`
	Platform string            `json:"platform"` // "web", "ios", "android", "js"
	Settings *CallbackSettings `json:"settings,omitempty"`
}

// CallbackSettings contains platform-specific callback URI settings.
type CallbackSettings struct {
	Origin                     string `json:"origin,omitempty"`
	BundleID                   string `json:"bundle_id,omitempty"`
	AppStoreID                 string `json:"app_store_id,omitempty"`
	TeamID                     string `json:"team_id,omitempty"`
	PackageName                string `json:"package_name,omitempty"`
	SHA1CertificateFingerprint string `json:"sha1_certificate_fingerprint,omitempty"`
}

// Branding contains branding details for the application.
type Branding struct {
	Name              string `json:"name,omitempty"`
	IconURL           string `json:"icon_url,omitempty"`
	WebsiteURL        string `json:"website_url,omitempty"`
	PrivacyPolicyURL  string `json:"privacy_policy_url,omitempty"`
	TermsOfServiceURL string `json:"terms_of_service_url,omitempty"`
}

// HostedAuthentication contains hosted authentication branding details.
type HostedAuthentication struct {
	BackgroundColor    string `json:"background_color,omitempty"`
	BackgroundImageURL string `json:"background_image_url,omitempty"`
	LogoURL            string `json:"logo_url,omitempty"`
	PrimaryColor       string `json:"primary_color,omitempty"`
	SecondaryColor     string `json:"secondary_color,omitempty"`
	Title              string `json:"title,omitempty"`
	Subtitle           string `json:"subtitle,omitempty"`
}
