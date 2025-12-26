package threads

import (
	"time"

	"github.com/mqasimca/nylas-go/common"
)

// Participant is an alias for common.Participant.
type Participant = common.Participant

// Thread represents an email thread (conversation).
type Thread struct {
	ID                   string        `json:"id"`
	GrantID              string        `json:"grant_id"`
	Object               string        `json:"object,omitempty"`
	LatestDraftOrMessage *MessageRef   `json:"latest_draft_or_message,omitempty"`
	HasAttachments       bool          `json:"has_attachments,omitempty"`
	HasDrafts            bool          `json:"has_drafts,omitempty"`
	Starred              bool          `json:"starred,omitempty"`
	Unread               bool          `json:"unread,omitempty"`
	EarliestMessageDate  int64         `json:"earliest_message_date,omitempty"`
	LatestMessageDate    int64         `json:"latest_message_date,omitempty"`
	MessageIDs           []string      `json:"message_ids,omitempty"`
	DraftIDs             []string      `json:"draft_ids,omitempty"`
	Participants         []Participant `json:"participants,omitempty"`
	Snippet              string        `json:"snippet,omitempty"`
	Subject              string        `json:"subject,omitempty"`
	Folders              []string      `json:"folders,omitempty"`
	Labels               []string      `json:"labels,omitempty"`
}

// MessageRef is a reference to the latest message or draft in a thread.
type MessageRef struct {
	ID      string        `json:"id"`
	Object  string        `json:"object,omitempty"`
	Subject string        `json:"subject,omitempty"`
	From    []Participant `json:"from,omitempty"`
	To      []Participant `json:"to,omitempty"`
	Date    int64         `json:"date,omitempty"`
	Snippet string        `json:"snippet,omitempty"`
}

// ListOptions specifies options for listing threads.
type ListOptions struct {
	Limit               *int     `json:"limit,omitempty"`
	PageToken           string   `json:"page_token,omitempty"`
	Subject             *string  `json:"subject,omitempty"`
	AnyEmail            []string `json:"any_email,omitempty"`
	From                *string  `json:"from,omitempty"`
	To                  *string  `json:"to,omitempty"`
	CC                  *string  `json:"cc,omitempty"`
	BCC                 *string  `json:"bcc,omitempty"`
	In                  *string  `json:"in,omitempty"`
	Unread              *bool    `json:"unread,omitempty"`
	Starred             *bool    `json:"starred,omitempty"`
	LatestMessageAfter  *int64   `json:"latest_message_after,omitempty"`
	LatestMessageBefore *int64   `json:"latest_message_before,omitempty"`
	HasAttachment       *bool    `json:"has_attachment,omitempty"`
	SearchQueryNative   *string  `json:"search_query_native,omitempty"`
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

// UpdateRequest represents a request to update a thread.
type UpdateRequest struct {
	Unread  *bool    `json:"unread,omitempty"`
	Starred *bool    `json:"starred,omitempty"`
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
