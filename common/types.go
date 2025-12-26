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
