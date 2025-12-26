package drafts

import (
	"time"

	"github.com/mqasimca/nylas-go/common"
)

// Participant is an alias for common.Participant.
type Participant = common.Participant

// Attachment is an alias for common.Attachment.
type Attachment = common.Attachment

// Draft represents an email draft.
type Draft struct {
	ID          string        `json:"id"`
	GrantID     string        `json:"grant_id"`
	Object      string        `json:"object,omitempty"`
	ThreadID    string        `json:"thread_id,omitempty"`
	Subject     string        `json:"subject,omitempty"`
	From        []Participant `json:"from,omitempty"`
	To          []Participant `json:"to,omitempty"`
	CC          []Participant `json:"cc,omitempty"`
	BCC         []Participant `json:"bcc,omitempty"`
	ReplyTo     []Participant `json:"reply_to,omitempty"`
	Date        int64         `json:"date,omitempty"`
	Body        string        `json:"body,omitempty"`
	Snippet     string        `json:"snippet,omitempty"`
	Starred     bool          `json:"starred,omitempty"`
	Folders     []string      `json:"folders,omitempty"`
	Labels      []string      `json:"labels,omitempty"`
	Attachments []Attachment  `json:"attachments,omitempty"`
	CreatedAt   int64         `json:"created_at,omitempty"`
}

// ListOptions specifies options for listing drafts.
type ListOptions struct {
	Limit         *int     `json:"limit,omitempty"`
	PageToken     string   `json:"page_token,omitempty"`
	Subject       *string  `json:"subject,omitempty"`
	AnyEmail      []string `json:"any_email,omitempty"`
	To            *string  `json:"to,omitempty"`
	CC            *string  `json:"cc,omitempty"`
	BCC           *string  `json:"bcc,omitempty"`
	Unread        *bool    `json:"unread,omitempty"`
	Starred       *bool    `json:"starred,omitempty"`
	ThreadID      *string  `json:"thread_id,omitempty"`
	HasAttachment *bool    `json:"has_attachment,omitempty"`
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
	if o.Subject != nil {
		v["subject"] = *o.Subject
	}
	if len(o.AnyEmail) > 0 {
		v["any_email"] = o.AnyEmail
	}
	if o.To != nil {
		v["to"] = *o.To
	}
	if o.CC != nil {
		v["cc"] = *o.CC
	}
	if o.BCC != nil {
		v["bcc"] = *o.BCC
	}
	if o.Unread != nil {
		v["unread"] = *o.Unread
	}
	if o.Starred != nil {
		v["starred"] = *o.Starred
	}
	if o.ThreadID != nil {
		v["thread_id"] = *o.ThreadID
	}
	if o.HasAttachment != nil {
		v["has_attachment"] = *o.HasAttachment
	}
	return v
}

// CreateRequest represents a request to create a draft.
type CreateRequest struct {
	Subject          string              `json:"subject,omitempty"`
	Body             string              `json:"body,omitempty"`
	From             []Participant       `json:"from,omitempty"`
	To               []Participant       `json:"to,omitempty"`
	CC               []Participant       `json:"cc,omitempty"`
	BCC              []Participant       `json:"bcc,omitempty"`
	ReplyTo          []Participant       `json:"reply_to,omitempty"`
	ReplyToMessageID string              `json:"reply_to_message_id,omitempty"`
	TrackingOptions  *TrackingOptions    `json:"tracking_options,omitempty"`
	Attachments      []AttachmentRequest `json:"attachments,omitempty"`
}

// UpdateRequest represents a request to update a draft.
type UpdateRequest struct {
	Subject          string              `json:"subject,omitempty"`
	Body             string              `json:"body,omitempty"`
	From             []Participant       `json:"from,omitempty"`
	To               []Participant       `json:"to,omitempty"`
	CC               []Participant       `json:"cc,omitempty"`
	BCC              []Participant       `json:"bcc,omitempty"`
	ReplyTo          []Participant       `json:"reply_to,omitempty"`
	ReplyToMessageID string              `json:"reply_to_message_id,omitempty"`
	Starred          *bool               `json:"starred,omitempty"`
	Attachments      []AttachmentRequest `json:"attachments,omitempty"`
}

// TrackingOptions specifies email tracking options.
type TrackingOptions struct {
	Opens         bool   `json:"opens,omitempty"`
	Links         bool   `json:"links,omitempty"`
	ThreadReplies bool   `json:"thread_replies,omitempty"`
	Label         string `json:"label,omitempty"`
}

// AttachmentRequest represents an attachment in a create/update request.
type AttachmentRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"` // base64 encoded
	ContentID   string `json:"content_id,omitempty"`
	IsInline    bool   `json:"is_inline,omitempty"`
}

// DateTime returns the draft date as time.Time.
func (d *Draft) DateTime() time.Time {
	return time.Unix(d.Date, 0)
}

// CreatedDateTime returns the created date as time.Time.
func (d *Draft) CreatedDateTime() time.Time {
	return time.Unix(d.CreatedAt, 0)
}
