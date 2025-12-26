package folders

// Folder represents an email folder (IMAP) or label (Gmail) in the Nylas API.
type Folder struct {
	// ID is the unique identifier for this folder.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this folder belongs to.
	GrantID string `json:"grant_id"`
	// Name is the folder display name.
	Name string `json:"name"`
	// Object is the object type, always "folder".
	Object string `json:"object,omitempty"`
	// ParentID is the ID of the parent folder (for nested folders).
	ParentID string `json:"parent_id,omitempty"`
	// BackgroundColor is the background color in hex format (Gmail labels).
	BackgroundColor string `json:"background_color,omitempty"`
	// TextColor is the text color in hex format (Gmail labels).
	TextColor string `json:"text_color,omitempty"`
	// SystemFolder indicates if this is a provider-created system folder.
	SystemFolder bool `json:"system_folder,omitempty"`
	// ChildCount is the number of child folders.
	ChildCount *int `json:"child_count,omitempty"`
	// UnreadCount is the number of unread messages in this folder.
	UnreadCount *int `json:"unread_count,omitempty"`
	// TotalCount is the total number of messages in this folder.
	TotalCount *int `json:"total_count,omitempty"`
	// Attributes contains IMAP folder attributes (e.g., "\\Sent", "\\Trash").
	Attributes []string `json:"attributes,omitempty"`
}

// ListOptions specifies options for listing folders.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of folders to return.
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
	// ParentID filters to child folders of a specific parent.
	ParentID *string `json:"parent_id,omitempty"`
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
	if o.ParentID != nil {
		v["parent_id"] = *o.ParentID
	}
	return v
}

// CreateRequest represents a request to create a folder or label.
type CreateRequest struct {
	// Name is the folder display name (required).
	Name string `json:"name"`
	// ParentID is the ID of the parent folder for nested folders.
	ParentID string `json:"parent_id,omitempty"`
	// BackgroundColor is the background color in hex format (Gmail labels).
	BackgroundColor string `json:"background_color,omitempty"`
	// TextColor is the text color in hex format (Gmail labels).
	TextColor string `json:"text_color,omitempty"`
}

// UpdateRequest represents a request to update a folder or label.
type UpdateRequest struct {
	// Name is the folder display name.
	Name *string `json:"name,omitempty"`
	// ParentID moves the folder to a new parent.
	ParentID *string `json:"parent_id,omitempty"`
	// BackgroundColor is the background color in hex format.
	BackgroundColor *string `json:"background_color,omitempty"`
	// TextColor is the text color in hex format.
	TextColor *string `json:"text_color,omitempty"`
}
