package drafts

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

// Draft represents an unsent email draft.
type Draft struct {
	// ID is the unique identifier for this draft.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this draft belongs to.
	GrantID string `json:"grant_id"`
	// Object is the object type, always "draft".
	Object string `json:"object,omitempty"`
	// ThreadID is the ID of the thread this draft belongs to (for replies).
	ThreadID string `json:"thread_id,omitempty"`
	// Subject is the subject line of the draft.
	Subject string `json:"subject,omitempty"`
	// From contains the sender address(es).
	From []Participant `json:"from,omitempty"`
	// To contains the primary recipients.
	To []Participant `json:"to,omitempty"`
	// CC contains the carbon copy recipients.
	CC []Participant `json:"cc,omitempty"`
	// BCC contains the blind carbon copy recipients.
	BCC []Participant `json:"bcc,omitempty"`
	// ReplyTo contains the addresses to use when replying.
	ReplyTo []Participant `json:"reply_to,omitempty"`
	// Date is the Unix timestamp of the draft.
	Date int64 `json:"date,omitempty"`
	// Body is the HTML or plain text content of the draft.
	Body string `json:"body,omitempty"`
	// Snippet is a preview of the draft body.
	Snippet string `json:"snippet,omitempty"`
	// Starred indicates whether the draft is starred.
	Starred bool `json:"starred,omitempty"`
	// Folders contains the folder IDs where this draft is stored.
	Folders []string `json:"folders,omitempty"`
	// Labels contains Gmail label IDs (Gmail only).
	Labels []string `json:"labels,omitempty"`
	// Attachments contains metadata about files attached to this draft.
	Attachments []Attachment `json:"attachments,omitempty"`
	// CreatedAt is the Unix timestamp when the draft was created.
	CreatedAt int64 `json:"created_at,omitempty"`
}

// ListOptions specifies options for listing drafts.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of drafts to return (default 50, max 200).
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
	// Subject filters drafts by subject line (substring match).
	Subject *string `json:"subject,omitempty"`
	// AnyEmail filters drafts where any participant matches these emails.
	AnyEmail []string `json:"any_email,omitempty"`
	// To filters drafts by recipient email address.
	To *string `json:"to,omitempty"`
	// CC filters drafts by CC recipient email address.
	CC *string `json:"cc,omitempty"`
	// BCC filters drafts by BCC recipient email address.
	BCC *string `json:"bcc,omitempty"`
	// Unread filters by unread status.
	Unread *bool `json:"unread,omitempty"`
	// Starred filters by starred status.
	Starred *bool `json:"starred,omitempty"`
	// ThreadID filters drafts by thread ID.
	ThreadID *string `json:"thread_id,omitempty"`
	// HasAttachment filters drafts with or without attachments.
	HasAttachment *bool `json:"has_attachment,omitempty"`
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
	// Subject is the subject line of the draft.
	Subject string `json:"subject,omitempty"`
	// Body is the HTML or plain text content.
	Body string `json:"body,omitempty"`
	// From overrides the default sender address.
	From []Participant `json:"from,omitempty"`
	// To is the list of primary recipients.
	To []Participant `json:"to,omitempty"`
	// CC is the list of carbon copy recipients.
	CC []Participant `json:"cc,omitempty"`
	// BCC is the list of blind carbon copy recipients.
	BCC []Participant `json:"bcc,omitempty"`
	// ReplyTo sets custom reply-to addresses.
	ReplyTo []Participant `json:"reply_to,omitempty"`
	// ReplyToMessageID is the ID of the message being replied to.
	ReplyToMessageID string `json:"reply_to_message_id,omitempty"`
	// TrackingOptions enables open/link tracking when sent.
	TrackingOptions *TrackingOptions `json:"tracking_options,omitempty"`
	// Attachments is the list of files to attach.
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
}

// UpdateRequest represents a request to update a draft.
type UpdateRequest struct {
	// Subject is the subject line of the draft.
	Subject string `json:"subject,omitempty"`
	// Body is the HTML or plain text content.
	Body string `json:"body,omitempty"`
	// From overrides the default sender address.
	From []Participant `json:"from,omitempty"`
	// To is the list of primary recipients.
	To []Participant `json:"to,omitempty"`
	// CC is the list of carbon copy recipients.
	CC []Participant `json:"cc,omitempty"`
	// BCC is the list of blind carbon copy recipients.
	BCC []Participant `json:"bcc,omitempty"`
	// ReplyTo sets custom reply-to addresses.
	ReplyTo []Participant `json:"reply_to,omitempty"`
	// ReplyToMessageID is the ID of the message being replied to.
	ReplyToMessageID string `json:"reply_to_message_id,omitempty"`
	// Starred sets the starred status.
	Starred *bool `json:"starred,omitempty"`
	// Attachments is the list of files to attach.
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
}

// DateTime returns the draft date as time.Time.
func (d *Draft) DateTime() time.Time {
	return time.Unix(d.Date, 0)
}

// CreatedDateTime returns the created date as time.Time.
func (d *Draft) CreatedDateTime() time.Time {
	return time.Unix(d.CreatedAt, 0)
}
