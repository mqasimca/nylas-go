package webhooks

// Webhook represents a webhook subscription in the Nylas API.
// Webhooks deliver notifications when events occur in connected accounts.
type Webhook struct {
	// ID is the unique identifier for this webhook.
	ID string `json:"id"`
	// Description is a human-readable description of the webhook.
	Description string `json:"description,omitempty"`
	// TriggerTypes lists the event types that trigger this webhook.
	// Examples: "message.created", "calendar.created", "event.updated".
	TriggerTypes []string `json:"trigger_types"`
	// WebhookURL is the HTTPS URL that receives webhook notifications.
	WebhookURL string `json:"webhook_url"`
	// WebhookSecret is used to verify webhook signatures.
	WebhookSecret string `json:"webhook_secret,omitempty"`
	// Status is the webhook status ("active", "failing", "failed", "paused").
	Status string `json:"status,omitempty"`
	// NotificationEmailAddresses receives alerts when the webhook fails.
	NotificationEmailAddresses []string `json:"notification_email_addresses,omitempty"`
	// StatusUpdatedAt is the Unix timestamp when the status last changed.
	StatusUpdatedAt int64 `json:"status_updated_at,omitempty"`
	// CreatedAt is the Unix timestamp when the webhook was created.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt is the Unix timestamp when the webhook was last modified.
	UpdatedAt int64 `json:"updated_at,omitempty"`
}

// ListOptions specifies options for listing webhooks.
type ListOptions struct {
	// Limit is the maximum number of webhooks to return.
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
}

// Values converts ListOptions to URL query parameters.
func (o *ListOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	v := make(map[string]any)
	if o.Limit != nil {
		v["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		v["page_token"] = o.PageToken
	}
	return v
}

// CreateRequest represents a request to create a webhook subscription.
type CreateRequest struct {
	// TriggerTypes lists the event types to subscribe to (required).
	TriggerTypes []string `json:"trigger_types"`
	// WebhookURL is the HTTPS URL to receive notifications (required).
	WebhookURL string `json:"webhook_url"`
	// Description is a human-readable description.
	Description string `json:"description,omitempty"`
	// NotificationEmailAddresses receives failure alerts.
	NotificationEmailAddresses []string `json:"notification_email_addresses,omitempty"`
}

// UpdateRequest represents a request to update a webhook subscription.
type UpdateRequest struct {
	// TriggerTypes replaces the event types to subscribe to.
	TriggerTypes []string `json:"trigger_types,omitempty"`
	// WebhookURL changes the notification URL.
	WebhookURL *string `json:"webhook_url,omitempty"`
	// Description updates the description.
	Description *string `json:"description,omitempty"`
	// NotificationEmailAddresses replaces the failure alert recipients.
	NotificationEmailAddresses []string `json:"notification_email_addresses,omitempty"`
}

// RotateSecretResponse contains the new webhook secret after rotation.
type RotateSecretResponse struct {
	// WebhookSecret is the newly generated secret for signature verification.
	WebhookSecret string `json:"webhook_secret"`
}

// IPAddressesResponse contains the list of Nylas IP addresses for webhook whitelisting.
type IPAddressesResponse struct {
	// IPAddresses is the list of IP addresses that Nylas uses to send webhooks.
	IPAddresses []string `json:"ip_addresses"`
	// UpdatedAt is the Unix timestamp when the IP list was last updated.
	UpdatedAt int64 `json:"updated_at,omitempty"`
}
