// Package common provides shared types used across Nylas API resources.
package common

// Participant represents an email participant (from, to, cc, bcc).
type Participant struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

// Attachment represents an email attachment.
type Attachment struct {
	ID          string `json:"id,omitempty"`
	Filename    string `json:"filename,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
	IsInline    bool   `json:"is_inline,omitempty"`
}

// AttachmentRequest represents an attachment to include when sending a message or draft.
type AttachmentRequest struct {
	// Filename is the name of the file as it will appear to recipients.
	Filename string `json:"filename"`
	// ContentType is the MIME type (e.g., "application/pdf", "image/png").
	ContentType string `json:"content_type"`
	// Content is the base64-encoded file data.
	Content string `json:"content"`
	// ContentID is used for inline attachments referenced in HTML body.
	ContentID string `json:"content_id,omitempty"`
	// IsInline indicates whether this is an inline attachment (e.g., embedded image).
	IsInline bool `json:"is_inline,omitempty"`
}

// TrackingOptions specifies email tracking options for sent messages.
type TrackingOptions struct {
	// Opens enables tracking when recipients open the email.
	Opens bool `json:"opens,omitempty"`
	// Links enables tracking when recipients click links.
	Links bool `json:"links,omitempty"`
	// ThreadReplies enables tracking when recipients reply.
	ThreadReplies bool `json:"thread_replies,omitempty"`
	// Label is a custom label for organizing tracked messages.
	Label string `json:"label,omitempty"`
}
