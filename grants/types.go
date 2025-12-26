package grants

// Grant represents a Nylas grant (connected email/calendar account).
// A grant is created when a user authenticates via OAuth.
type Grant struct {
	// ID is the unique identifier for this grant.
	ID string `json:"id"`
	// Provider is the email provider ("google", "microsoft", "icloud", "yahoo", "ews", "imap").
	Provider string `json:"provider"`
	// GrantStatus is the current status ("valid", "invalid", "pending").
	GrantStatus string `json:"grant_status,omitempty"`
	// Email is the email address of the connected account.
	Email string `json:"email,omitempty"`
	// Scope is the list of OAuth scopes granted.
	Scope []string `json:"scope,omitempty"`
	// UserAgent is the user agent from the OAuth flow.
	UserAgent string `json:"user_agent,omitempty"`
	// IP is the IP address from the OAuth flow.
	IP string `json:"ip,omitempty"`
	// State is the OAuth state parameter value.
	State string `json:"state,omitempty"`
	// CreatedAt is the Unix timestamp when the grant was created.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt is the Unix timestamp when the grant was last modified.
	UpdatedAt int64 `json:"updated_at,omitempty"`
	// Settings contains provider-specific settings.
	Settings map[string]any `json:"settings,omitempty"`
}

// ListOptions specifies options for listing grants.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of grants to return.
	Limit *int `json:"limit,omitempty"`
	// Offset skips this many results (for pagination).
	Offset *int `json:"offset,omitempty"`
	// SortBy specifies the field to sort by.
	SortBy *string `json:"sort_by,omitempty"`
	// OrderBy specifies the sort order ("asc" or "desc").
	OrderBy *string `json:"order_by,omitempty"`
	// Since filters grants created after this Unix timestamp.
	Since *int64 `json:"since,omitempty"`
	// Before filters grants created before this Unix timestamp.
	Before *int64 `json:"before,omitempty"`
	// Email filters grants by email address.
	Email *string `json:"email,omitempty"`
	// GrantStatus filters grants by status.
	GrantStatus *string `json:"grant_status,omitempty"`
	// IP filters grants by IP address.
	IP *string `json:"ip,omitempty"`
	// Provider filters grants by provider.
	Provider *string `json:"provider,omitempty"`
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
	if o.Offset != nil {
		v["offset"] = *o.Offset
	}
	if o.SortBy != nil {
		v["sort_by"] = *o.SortBy
	}
	if o.OrderBy != nil {
		v["order_by"] = *o.OrderBy
	}
	if o.Since != nil {
		v["since"] = *o.Since
	}
	if o.Before != nil {
		v["before"] = *o.Before
	}
	if o.Email != nil {
		v["email"] = *o.Email
	}
	if o.GrantStatus != nil {
		v["grant_status"] = *o.GrantStatus
	}
	if o.IP != nil {
		v["ip"] = *o.IP
	}
	if o.Provider != nil {
		v["provider"] = *o.Provider
	}
	return v
}

// UpdateRequest represents a request to update a grant's settings.
type UpdateRequest struct {
	// Settings contains provider-specific settings to update.
	Settings map[string]any `json:"settings,omitempty"`
	// Scope updates the OAuth scopes for the grant.
	Scope []string `json:"scope,omitempty"`
}
