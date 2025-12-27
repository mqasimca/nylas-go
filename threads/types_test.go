package threads

import (
	"testing"
	"time"
)

func TestListOptions_Values(t *testing.T) {
	tests := []struct {
		name string
		opts *ListOptions
		want int // number of expected params
	}{
		{
			name: "nil options",
			opts: nil,
			want: 0,
		},
		{
			name: "empty options",
			opts: &ListOptions{},
			want: 0,
		},
		{
			name: "with limit",
			opts: &ListOptions{Limit: ptr(10)},
			want: 1,
		},
		{
			name: "with multiple options",
			opts: &ListOptions{
				Limit:     ptr(50),
				Unread:    ptrBool(true),
				PageToken: "abc",
			},
			want: 3,
		},
		{
			name: "with any_email",
			opts: &ListOptions{
				AnyEmail: []string{"a@test.com", "b@test.com"},
			},
			want: 1,
		},
		{
			name: "with latest_message_after",
			opts: &ListOptions{
				LatestMessageAfter: ptrInt64(1704067200),
			},
			want: 1,
		},
		{
			name: "with subject",
			opts: &ListOptions{Subject: ptrStr("test subject")},
			want: 1,
		},
		{
			name: "with from",
			opts: &ListOptions{From: ptrStr("sender@test.com")},
			want: 1,
		},
		{
			name: "with to",
			opts: &ListOptions{To: ptrStr("recipient@test.com")},
			want: 1,
		},
		{
			name: "with cc",
			opts: &ListOptions{CC: ptrStr("cc@test.com")},
			want: 1,
		},
		{
			name: "with bcc",
			opts: &ListOptions{BCC: ptrStr("bcc@test.com")},
			want: 1,
		},
		{
			name: "with in",
			opts: &ListOptions{In: ptrStr("INBOX")},
			want: 1,
		},
		{
			name: "with starred",
			opts: &ListOptions{Starred: ptrBool(true)},
			want: 1,
		},
		{
			name: "with latest_message_before",
			opts: &ListOptions{LatestMessageBefore: ptrInt64(1704153600)},
			want: 1,
		},
		{
			name: "with has_attachment",
			opts: &ListOptions{HasAttachment: ptrBool(true)},
			want: 1,
		},
		{
			name: "with search_query_native",
			opts: &ListOptions{SearchQueryNative: ptrStr("from:test@example.com")},
			want: 1,
		},
		{
			name: "with all options",
			opts: &ListOptions{
				Limit:               ptr(100),
				PageToken:           "token123",
				Subject:             ptrStr("Hello"),
				AnyEmail:            []string{"a@test.com"},
				From:                ptrStr("from@test.com"),
				To:                  ptrStr("to@test.com"),
				CC:                  ptrStr("cc@test.com"),
				BCC:                 ptrStr("bcc@test.com"),
				In:                  ptrStr("INBOX"),
				Unread:              ptrBool(false),
				Starred:             ptrBool(true),
				LatestMessageAfter:  ptrInt64(1000),
				LatestMessageBefore: ptrInt64(2000),
				HasAttachment:       ptrBool(true),
				SearchQueryNative:   ptrStr("search query"),
			},
			want: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.opts.Values()
			if v == nil && tt.want > 0 {
				t.Error("Values() returned nil, expected map")
				return
			}
			if v != nil && len(v) != tt.want {
				t.Errorf("Values() len = %d, want %d", len(v), tt.want)
			}
		})
	}
}

func TestThread_EarliestMessageDateTime(t *testing.T) {
	thread := &Thread{
		EarliestMessageDate: 1704067200, // 2024-01-01 00:00:00 UTC
	}

	dt := thread.EarliestMessageDateTime()
	expected := time.Unix(1704067200, 0)

	if !dt.Equal(expected) {
		t.Errorf("EarliestMessageDateTime() = %v, want %v", dt, expected)
	}
}

func TestThread_LatestMessageDateTime(t *testing.T) {
	thread := &Thread{
		LatestMessageDate: 1704153600, // 2024-01-02 00:00:00 UTC
	}

	dt := thread.LatestMessageDateTime()
	expected := time.Unix(1704153600, 0)

	if !dt.Equal(expected) {
		t.Errorf("LatestMessageDateTime() = %v, want %v", dt, expected)
	}
}

func TestThread_MessageCount(t *testing.T) {
	thread := &Thread{
		MessageIDs: []string{"msg-1", "msg-2", "msg-3"},
	}

	if count := thread.MessageCount(); count != 3 {
		t.Errorf("MessageCount() = %d, want 3", count)
	}
}

func TestThread_DraftCount(t *testing.T) {
	thread := &Thread{
		DraftIDs: []string{"draft-1", "draft-2"},
	}

	if count := thread.DraftCount(); count != 2 {
		t.Errorf("DraftCount() = %d, want 2", count)
	}
}

func TestParticipant(t *testing.T) {
	p := Participant{
		Name:  "Test User",
		Email: "test@example.com",
	}

	if p.Name != "Test User" {
		t.Errorf("Name = %s, want Test User", p.Name)
	}
	if p.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", p.Email)
	}
}

func TestMessageRef(t *testing.T) {
	ref := &MessageRef{
		ID:      "msg-123",
		Object:  "message",
		Subject: "Test Subject",
		Date:    1704067200,
	}

	if ref.ID != "msg-123" {
		t.Errorf("ID = %s, want msg-123", ref.ID)
	}
	if ref.Subject != "Test Subject" {
		t.Errorf("Subject = %s, want Test Subject", ref.Subject)
	}
}

// helpers
func ptr(v int) *int          { return &v }
func ptrBool(v bool) *bool    { return &v }
func ptrInt64(v int64) *int64 { return &v }
func ptrStr(v string) *string { return &v }
