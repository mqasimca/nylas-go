// Package connectors provides types for the Nylas Connectors API.
package connectors

// Provider represents supported email provider types.
type Provider string

const (
	ProviderGoogle           Provider = "google"
	ProviderMicrosoft        Provider = "microsoft"
	ProviderIMAP             Provider = "imap"
	ProviderVirtualCalendars Provider = "virtual-calendars"
)

// Connector represents a connector configuration.
type Connector struct {
	Provider Provider               `json:"provider"`
	Name     string                 `json:"name,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Scope    []string               `json:"scope,omitempty"`
}

// GoogleSettings contains Google-specific connector settings.
type GoogleSettings struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TopicName    string `json:"topic_name,omitempty"`
}

// MicrosoftSettings contains Microsoft-specific connector settings.
type MicrosoftSettings struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Tenant       string `json:"tenant,omitempty"`
}

// CreateRequest contains the data to create a connector.
type CreateRequest struct {
	Name     string                 `json:"name"`
	Provider Provider               `json:"provider"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Scope    []string               `json:"scope,omitempty"`
}

// CreateGoogleRequest creates a Google connector request.
func CreateGoogleRequest(name string, clientID, clientSecret string, scope []string) *CreateRequest {
	return &CreateRequest{
		Name:     name,
		Provider: ProviderGoogle,
		Settings: map[string]interface{}{
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
		Scope: scope,
	}
}

// CreateMicrosoftRequest creates a Microsoft connector request.
func CreateMicrosoftRequest(name string, clientID, clientSecret, tenant string, scope []string) *CreateRequest {
	req := &CreateRequest{
		Name:     name,
		Provider: ProviderMicrosoft,
		Settings: map[string]interface{}{
			"client_id":     clientID,
			"client_secret": clientSecret,
		},
		Scope: scope,
	}
	if tenant != "" {
		req.Settings["tenant"] = tenant
	}
	return req
}

// UpdateRequest contains the data to update a connector.
type UpdateRequest struct {
	Name     *string                `json:"name,omitempty"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Scope    []string               `json:"scope,omitempty"`
}

// ListOptions contains options for listing connectors.
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
