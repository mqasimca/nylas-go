// Package redirecturis provides types for the Nylas Redirect URIs API.
package redirecturis

// RedirectURI represents a redirect URI configuration.
type RedirectURI struct {
	ID       string            `json:"id"`
	URL      string            `json:"url"`
	Platform string            `json:"platform"` // "web", "ios", "android", "js"
	Settings *RedirectSettings `json:"settings,omitempty"`
}

// RedirectSettings contains platform-specific redirect URI settings.
type RedirectSettings struct {
	// JavaScript platform
	Origin string `json:"origin,omitempty"`

	// iOS platform
	BundleID   string `json:"bundle_id,omitempty"`
	AppStoreID string `json:"app_store_id,omitempty"`
	TeamID     string `json:"team_id,omitempty"`

	// Android platform
	PackageName                string `json:"package_name,omitempty"`
	SHA1CertificateFingerprint string `json:"sha1_certificate_fingerprint,omitempty"`
}

// CreateRequest contains the data to create a redirect URI.
type CreateRequest struct {
	URL      string            `json:"url"`
	Platform string            `json:"platform"` // "web", "ios", "android", "js"
	Settings *RedirectSettings `json:"settings,omitempty"`
}

// UpdateRequest contains the data to update a redirect URI.
type UpdateRequest struct {
	URL      *string           `json:"url,omitempty"`
	Platform *string           `json:"platform,omitempty"`
	Settings *RedirectSettings `json:"settings,omitempty"`
}

// ListOptions contains options for listing redirect URIs.
type ListOptions struct {
	Limit     *int   `json:"limit,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

// Values returns the options as a map for query parameters.
func (o *ListOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	m := make(map[string]any)
	if o.Limit != nil {
		m["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		m["page_token"] = o.PageToken
	}
	return m
}
