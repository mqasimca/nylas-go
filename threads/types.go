package threads

import (
	"time"

	"github.com/mqasimca/nylas-go/common"
)

// Participant is an alias for common.Participant.
type Participant = common.Participant

// Thread represents an email thread (conversation) containing one or more messages.
type Thread struct {
	// ID is the unique identifier for this thread.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this thread belongs to.
	GrantID string `json:"grant_id"`
	// Object is the object type, always "thread".
	Object string `json:"object,omitempty"`
	// LatestDraftOrMessage contains a reference to the most recent message or draft.
	LatestDraftOrMessage *MessageRef `json:"latest_draft_or_message,omitempty"`
	// HasAttachments indicates whether any message in the thread has attachments.
	HasAttachments bool `json:"has_attachments,omitempty"`
	// HasDrafts indicates whether the thread contains any drafts.
	HasDrafts bool `json:"has_drafts,omitempty"`
	// Starred indicates whether the thread is starred/flagged.
	Starred bool `json:"starred,omitempty"`
	// Unread indicates whether any message in the thread is unread.
	Unread bool `json:"unread,omitempty"`
	// EarliestMessageDate is the Unix timestamp of the oldest message.
	EarliestMessageDate int64 `json:"earliest_message_date,omitempty"`
	// LatestMessageDate is the Unix timestamp of the newest message.
	LatestMessageDate int64 `json:"latest_message_date,omitempty"`
	// MessageIDs contains the IDs of all messages in this thread.
	MessageIDs []string `json:"message_ids,omitempty"`
	// DraftIDs contains the IDs of all drafts in this thread.
	DraftIDs []string `json:"draft_ids,omitempty"`
	// Participants contains all email addresses involved in this thread.
	Participants []Participant `json:"participants,omitempty"`
	// Snippet is a preview of the latest message body.
	Snippet string `json:"snippet,omitempty"`
	// Subject is the subject line of the thread.
	Subject string `json:"subject,omitempty"`
	// Folders contains the folder IDs where this thread appears.
	Folders []string `json:"folders,omitempty"`
	// Labels contains Gmail label IDs (Gmail only).
	Labels []string `json:"labels,omitempty"`
}

// MessageRef is a reference to the latest message or draft in a thread.
type MessageRef struct {
	// ID is the message or draft ID.
	ID string `json:"id"`
	// Object is "message" or "draft".
	Object string `json:"object,omitempty"`
	// Subject is the subject line.
	Subject string `json:"subject,omitempty"`
	// From contains the sender(s).
	From []Participant `json:"from,omitempty"`
	// To contains the recipient(s).
	To []Participant `json:"to,omitempty"`
	// Date is the Unix timestamp when sent.
	Date int64 `json:"date,omitempty"`
	// Snippet is a preview of the message body.
	Snippet string `json:"snippet,omitempty"`
}

// ListOptions specifies options for listing threads.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of threads to return (default 50, max 200).
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
	// Subject filters threads by subject line (substring match).
	Subject *string `json:"subject,omitempty"`
	// AnyEmail filters threads where any participant matches these emails.
	AnyEmail []string `json:"any_email,omitempty"`
	// From filters threads by sender email address.
	From *string `json:"from,omitempty"`
	// To filters threads by recipient email address.
	To *string `json:"to,omitempty"`
	// CC filters threads by CC recipient email address.
	CC *string `json:"cc,omitempty"`
	// BCC filters threads by BCC recipient email address.
	BCC *string `json:"bcc,omitempty"`
	// In filters threads by folder ID.
	In *string `json:"in,omitempty"`
	// Unread filters by unread status.
	Unread *bool `json:"unread,omitempty"`
	// Starred filters by starred status.
	Starred *bool `json:"starred,omitempty"`
	// LatestMessageAfter filters threads with latest message after this Unix timestamp.
	LatestMessageAfter *int64 `json:"latest_message_after,omitempty"`
	// LatestMessageBefore filters threads with latest message before this Unix timestamp.
	LatestMessageBefore *int64 `json:"latest_message_before,omitempty"`
	// HasAttachment filters threads with or without attachments.
	HasAttachment *bool `json:"has_attachment,omitempty"`
	// SearchQueryNative is a provider-specific search query.
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
	if o.LatestMessageAfter != nil {
		v["latest_message_after"] = *o.LatestMessageAfter
	}
	if o.LatestMessageBefore != nil {
		v["latest_message_before"] = *o.LatestMessageBefore
	}
	if o.HasAttachment != nil {
		v["has_attachment"] = *o.HasAttachment
	}
	if o.SearchQueryNative != nil {
		v["search_query_native"] = *o.SearchQueryNative
	}
	return v
}

// UpdateRequest represents a request to update a thread's metadata.
type UpdateRequest struct {
	// Unread sets the read/unread status for all messages in the thread.
	Unread *bool `json:"unread,omitempty"`
	// Starred sets the starred/flagged status.
	Starred *bool `json:"starred,omitempty"`
	// Folders moves all messages in the thread to these folder IDs.
	Folders []string `json:"folders,omitempty"`
}

// EarliestMessageDateTime returns the earliest message date as time.Time.
func (t *Thread) EarliestMessageDateTime() time.Time {
	return time.Unix(t.EarliestMessageDate, 0)
}

// LatestMessageDateTime returns the latest message date as time.Time.
func (t *Thread) LatestMessageDateTime() time.Time {
	return time.Unix(t.LatestMessageDate, 0)
}

// MessageCount returns the number of messages in the thread.
func (t *Thread) MessageCount() int {
	return len(t.MessageIDs)
}

// DraftCount returns the number of drafts in the thread.
func (t *Thread) DraftCount() int {
	return len(t.DraftIDs)
}
