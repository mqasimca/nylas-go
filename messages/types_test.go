package messages

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
			name: "with thread_id",
			opts: &ListOptions{ThreadID: ptrStr("thread-123")},
			want: 1,
		},
		{
			name: "with received_after",
			opts: &ListOptions{ReceivedAfter: ptrInt64(1704067200)},
			want: 1,
		},
		{
			name: "with received_before",
			opts: &ListOptions{ReceivedBefore: ptrInt64(1704153600)},
			want: 1,
		},
		{
			name: "with has_attachment",
			opts: &ListOptions{HasAttachment: ptrBool(true)},
			want: 1,
		},
		{
			name: "with fields",
			opts: &ListOptions{Fields: ptrStr("include_headers")},
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
				Limit:             ptr(100),
				PageToken:         "token123",
				Subject:           ptrStr("Hello"),
				AnyEmail:          []string{"a@test.com"},
				From:              ptrStr("from@test.com"),
				To:                ptrStr("to@test.com"),
				CC:                ptrStr("cc@test.com"),
				BCC:               ptrStr("bcc@test.com"),
				In:                ptrStr("INBOX"),
				Unread:            ptrBool(false),
				Starred:           ptrBool(true),
				ThreadID:          ptrStr("thread-1"),
				ReceivedAfter:     ptrInt64(1000),
				ReceivedBefore:    ptrInt64(2000),
				HasAttachment:     ptrBool(true),
				Fields:            ptrStr("headers"),
				SearchQueryNative: ptrStr("search query"),
			},
			want: 17,
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

func TestMessage_DateTime(t *testing.T) {
	msg := &Message{
		Date: 1704067200, // 2024-01-01 00:00:00 UTC
	}

	dt := msg.DateTime()
	expected := time.Unix(1704067200, 0)

	if !dt.Equal(expected) {
		t.Errorf("DateTime() = %v, want %v", dt, expected)
	}
}

func TestMessage_CreatedDateTime(t *testing.T) {
	msg := &Message{
		CreatedAt: 1704067200,
	}

	dt := msg.CreatedDateTime()
	expected := time.Unix(1704067200, 0)

	if !dt.Equal(expected) {
		t.Errorf("CreatedDateTime() = %v, want %v", dt, expected)
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

func TestAttachment(t *testing.T) {
	a := Attachment{
		ID:          "att-1",
		Filename:    "test.pdf",
		ContentType: "application/pdf",
		Size:        1024,
	}

	if a.ID != "att-1" {
		t.Errorf("ID = %s, want att-1", a.ID)
	}
	if a.Size != 1024 {
		t.Errorf("Size = %d, want 1024", a.Size)
	}
}

// helpers
func ptr(v int) *int          { return &v }
func ptrBool(v bool) *bool    { return &v }
func ptrStr(v string) *string { return &v }
func ptrInt64(v int64) *int64 { return &v }
