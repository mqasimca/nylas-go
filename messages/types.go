package messages

import (
	"time"

	"github.com/mqasimca/nylas-go/common"
)

// Participant is an alias for common.Participant.
type Participant = common.Participant

// Attachment is an alias for common.Attachment.
type Attachment = common.Attachment

// AttachmentRequest is an alias for common.AttachmentRequest.
type AttachmentRequest = common.AttachmentRequest

// TrackingOptions is an alias for common.TrackingOptions.
type TrackingOptions = common.TrackingOptions

// Message represents an email message from the Nylas API.
type Message struct {
	// ID is the unique identifier for this message.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this message belongs to.
	GrantID string `json:"grant_id"`
	// ThreadID is the ID of the thread this message belongs to.
	ThreadID string `json:"thread_id,omitempty"`
	// Subject is the subject line of the message.
	Subject string `json:"subject,omitempty"`
	// From contains the sender(s) of the message.
	From []Participant `json:"from,omitempty"`
	// To contains the primary recipients of the message.
	To []Participant `json:"to,omitempty"`
	// CC contains the carbon copy recipients.
	CC []Participant `json:"cc,omitempty"`
	// BCC contains the blind carbon copy recipients.
	BCC []Participant `json:"bcc,omitempty"`
	// ReplyTo contains the addresses to use when replying.
	ReplyTo []Participant `json:"reply_to,omitempty"`
	// Date is the Unix timestamp when the message was sent.
	Date int64 `json:"date,omitempty"`
	// Body is the full HTML or plain text content of the message.
	Body string `json:"body,omitempty"`
	// Snippet is a short preview of the message body (first ~100 chars).
	Snippet string `json:"snippet,omitempty"`
	// Unread indicates whether the message has been read.
	Unread bool `json:"unread,omitempty"`
	// Starred indicates whether the message is starred/flagged.
	Starred bool `json:"starred,omitempty"`
	// Folders contains the folder IDs where this message is stored.
	Folders []string `json:"folders,omitempty"`
	// Labels contains Gmail label IDs (Gmail only).
	Labels []string `json:"labels,omitempty"`
	// Attachments contains metadata about files attached to this message.
	Attachments []Attachment `json:"attachments,omitempty"`
	// CreatedAt is the Unix timestamp when the message was created in Nylas.
	CreatedAt int64 `json:"created_at,omitempty"`
	// Object is the object type, always "message".
	Object string `json:"object,omitempty"`
}

// ListOptions specifies options for listing messages.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of messages to return (default 50, max 200).
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination; use NextCursor from previous response.
	PageToken string `json:"page_token,omitempty"`
	// Subject filters messages by subject line (substring match).
	Subject *string `json:"subject,omitempty"`
	// AnyEmail filters messages where any participant matches these emails.
	AnyEmail []string `json:"any_email,omitempty"`
	// From filters messages by sender email address.
	From *string `json:"from,omitempty"`
	// To filters messages by recipient email address.
	To *string `json:"to,omitempty"`
	// CC filters messages by CC recipient email address.
	CC *string `json:"cc,omitempty"`
	// BCC filters messages by BCC recipient email address.
	BCC *string `json:"bcc,omitempty"`
	// In filters messages by folder ID.
	In *string `json:"in,omitempty"`
	// Unread filters by read status: true for unread, false for read.
	Unread *bool `json:"unread,omitempty"`
	// Starred filters by starred status.
	Starred *bool `json:"starred,omitempty"`
	// ThreadID filters messages by thread ID.
	ThreadID *string `json:"thread_id,omitempty"`
	// ReceivedAfter filters messages received after this Unix timestamp.
	ReceivedAfter *int64 `json:"received_after,omitempty"`
	// ReceivedBefore filters messages received before this Unix timestamp.
	ReceivedBefore *int64 `json:"received_before,omitempty"`
	// HasAttachment filters messages with or without attachments.
	HasAttachment *bool `json:"has_attachment,omitempty"`
	// Fields specifies which fields to return (comma-separated).
	Fields *string `json:"fields,omitempty"`
	// SearchQueryNative is a provider-specific search query (Gmail, Microsoft, etc.).
	SearchQueryNative *string `json:"search_query_native,omitempty"`
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
	// To is the list of primary recipients (required).
	To []Participant `json:"to"`
	// From overrides the default sender address.
	From []Participant `json:"from,omitempty"`
	// CC is the list of carbon copy recipients.
	CC []Participant `json:"cc,omitempty"`
	// BCC is the list of blind carbon copy recipients.
	BCC []Participant `json:"bcc,omitempty"`
	// ReplyTo sets custom reply-to addresses.
	ReplyTo []Participant `json:"reply_to,omitempty"`
	// Subject is the subject line of the message.
	Subject string `json:"subject,omitempty"`
	// Body is the HTML or plain text content of the message.
	Body string `json:"body,omitempty"`
	// ReplyToMessageID is the ID of the message being replied to.
	ReplyToMessageID string `json:"reply_to_message_id,omitempty"`
	// TrackingOptions enables open/link tracking.
	TrackingOptions *TrackingOptions `json:"tracking_options,omitempty"`
	// SendAt schedules the message for future delivery (Unix timestamp).
	SendAt *int64 `json:"send_at,omitempty"`
	// UseReplyTo determines whether to use reply_to addresses.
	UseReplyTo *bool `json:"use_draft,omitempty"`
	// Attachments is the list of files to attach.
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
}

// UpdateRequest represents a request to update a message's metadata.
type UpdateRequest struct {
	// Unread sets the read/unread status.
	Unread *bool `json:"unread,omitempty"`
	// Starred sets the starred/flagged status.
	Starred *bool `json:"starred,omitempty"`
	// Folders moves the message to these folder IDs.
	Folders []string `json:"folders,omitempty"`
}

// ScheduledMessage represents a message scheduled for future delivery.
type ScheduledMessage struct {
	// ScheduleID is the unique identifier for this scheduled message.
	ScheduleID string `json:"schedule_id"`
	// Status is the current status (e.g., "scheduled", "sent", "cancelled").
	Status string `json:"status"`
	// CloseTime is the Unix timestamp when the message will be sent.
	CloseTime int64 `json:"close_time,omitempty"`
}

// ScheduledMessagesList represents a list of scheduled messages.
// Note: Nylas API returns an array directly, not wrapped in an object.
type ScheduledMessagesList []ScheduledMessage

// CleanRequest represents a request to extract clean conversation text from messages.
type CleanRequest struct {
	// MessageID is the list of message IDs to clean.
	MessageID []string `json:"message_id"`
	// IgnoreLinks removes hyperlinks from the output.
	IgnoreLinks *bool `json:"ignore_links,omitempty"`
	// IgnoreImages removes image references from the output.
	IgnoreImages *bool `json:"ignore_images,omitempty"`
	// IgnoreTables removes table formatting from the output.
	IgnoreTables *bool `json:"ignore_tables,omitempty"`
	// ImagesAsMarkdown converts images to markdown syntax.
	ImagesAsMarkdown *bool `json:"images_as_markdown,omitempty"`
	// RemoveConclusion removes signature and closing phrases.
	RemoveConclusion *bool `json:"remove_conclusion_phrases,omitempty"`
}

// CleanResponse represents the cleaned text extracted from a message.
type CleanResponse struct {
	// ID is the Nylas ID of the cleaned message.
	ID string `json:"id"`
	// GrantID is the grant ID the message belongs to.
	GrantID string `json:"grant_id"`
	// Conversation is the extracted clean conversation text.
	Conversation string `json:"conversation,omitempty"`
	// MessageID is the original message ID.
	MessageID string `json:"message_id,omitempty"`
}

// DateTime helper for formatting.
func (m *Message) DateTime() time.Time {
	return time.Unix(m.Date, 0)
}

// CreatedDateTime helper for formatting.
func (m *Message) CreatedDateTime() time.Time {
	return time.Unix(m.CreatedAt, 0)
}
