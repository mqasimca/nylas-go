package messages

import (
	"time"

	"github.com/mqasimca/nylas-go/common"
)

// Participant is an alias for common.Participant.
type Participant = common.Participant

// Attachment is an alias for common.Attachment.
type Attachment = common.Attachment

// Message represents an email message.
type Message struct {
	ID          string        `json:"id"`
	GrantID     string        `json:"grant_id"`
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
	Unread      bool          `json:"unread,omitempty"`
	Starred     bool          `json:"starred,omitempty"`
	Folders     []string      `json:"folders,omitempty"`
	Labels      []string      `json:"labels,omitempty"`
	Attachments []Attachment  `json:"attachments,omitempty"`
	CreatedAt   int64         `json:"created_at,omitempty"`
	Object      string        `json:"object,omitempty"`
}

// ListOptions specifies options for listing messages.
type ListOptions struct {
	Limit             *int     `json:"limit,omitempty"`
	PageToken         string   `json:"page_token,omitempty"`
	Subject           *string  `json:"subject,omitempty"`
	AnyEmail          []string `json:"any_email,omitempty"`
	From              *string  `json:"from,omitempty"`
	To                *string  `json:"to,omitempty"`
	CC                *string  `json:"cc,omitempty"`
	BCC               *string  `json:"bcc,omitempty"`
	In                *string  `json:"in,omitempty"`
	Unread            *bool    `json:"unread,omitempty"`
	Starred           *bool    `json:"starred,omitempty"`
	ThreadID          *string  `json:"thread_id,omitempty"`
	ReceivedAfter     *int64   `json:"received_after,omitempty"`
	ReceivedBefore    *int64   `json:"received_before,omitempty"`
	HasAttachment     *bool    `json:"has_attachment,omitempty"`
	Fields            *string  `json:"fields,omitempty"`
	SearchQueryNative *string  `json:"search_query_native,omitempty"`
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
	if o.From != nil {
		v["from"] = *o.From
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
	if o.In != nil {
		v["in"] = *o.In
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
	if o.ReceivedAfter != nil {
		v["received_after"] = *o.ReceivedAfter
	}
	if o.ReceivedBefore != nil {
		v["received_before"] = *o.ReceivedBefore
	}
	if o.HasAttachment != nil {
		v["has_attachment"] = *o.HasAttachment
	}
	if o.Fields != nil {
		v["fields"] = *o.Fields
	}
	if o.SearchQueryNative != nil {
		v["search_query_native"] = *o.SearchQueryNative
	}
	return v
}

// SendRequest represents a request to send a message.
type SendRequest struct {
	To               []Participant       `json:"to"`
	From             []Participant       `json:"from,omitempty"`
	CC               []Participant       `json:"cc,omitempty"`
	BCC              []Participant       `json:"bcc,omitempty"`
	ReplyTo          []Participant       `json:"reply_to,omitempty"`
	Subject          string              `json:"subject,omitempty"`
	Body             string              `json:"body,omitempty"`
	ReplyToMessageID string              `json:"reply_to_message_id,omitempty"`
	TrackingOptions  *TrackingOptions    `json:"tracking_options,omitempty"`
	SendAt           *int64              `json:"send_at,omitempty"`
	UseReplyTo       *bool               `json:"use_draft,omitempty"`
	Attachments      []AttachmentRequest `json:"attachments,omitempty"`
}

// TrackingOptions specifies email tracking options.
type TrackingOptions struct {
	Opens         bool   `json:"opens,omitempty"`
	Links         bool   `json:"links,omitempty"`
	ThreadReplies bool   `json:"thread_replies,omitempty"`
	Label         string `json:"label,omitempty"`
}

// AttachmentRequest represents an attachment in a send request.
type AttachmentRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"` // base64 encoded
	ContentID   string `json:"content_id,omitempty"`
	IsInline    bool   `json:"is_inline,omitempty"`
}

// UpdateRequest represents a request to update a message.
type UpdateRequest struct {
	Unread  *bool    `json:"unread,omitempty"`
	Starred *bool    `json:"starred,omitempty"`
	Folders []string `json:"folders,omitempty"`
}

// ScheduledMessage represents a scheduled message.
type ScheduledMessage struct {
	ScheduleID string `json:"schedule_id"`
	Status     string `json:"status"`
	CloseTime  int64  `json:"close_time,omitempty"`
}

// ScheduledMessagesList represents a list of scheduled messages.
// Note: Nylas API returns an array directly, not wrapped in an object.
type ScheduledMessagesList []ScheduledMessage

// CleanRequest represents a request to clean messages.
type CleanRequest struct {
	MessageID        []string `json:"message_id"`
	IgnoreLinks      *bool    `json:"ignore_links,omitempty"`
	IgnoreImages     *bool    `json:"ignore_images,omitempty"`
	IgnoreTables     *bool    `json:"ignore_tables,omitempty"`
	ImagesAsMarkdown *bool    `json:"images_as_markdown,omitempty"`
	RemoveConclusion *bool    `json:"remove_conclusion_phrases,omitempty"`
}

// CleanResponse represents a cleaned message.
type CleanResponse struct {
	ID           string `json:"id"`
	GrantID      string `json:"grant_id"`
	Conversation string `json:"conversation,omitempty"`
	MessageID    string `json:"message_id,omitempty"`
}

// DateTime helper for formatting.
func (m *Message) DateTime() time.Time {
	return time.Unix(m.Date, 0)
}

// CreatedDateTime helper for formatting.
func (m *Message) CreatedDateTime() time.Time {
	return time.Unix(m.CreatedAt, 0)
}
