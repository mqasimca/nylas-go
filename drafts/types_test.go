package drafts

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
				Starred:   ptrBool(true),
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
			name: "with thread_id",
			opts: &ListOptions{
				ThreadID: ptrString("thread-123"),
			},
			want: 1,
		},
		{
			name: "with subject",
			opts: &ListOptions{Subject: ptrString("test subject")},
			want: 1,
		},
		{
			name: "with to",
			opts: &ListOptions{To: ptrString("recipient@test.com")},
			want: 1,
		},
		{
			name: "with cc",
			opts: &ListOptions{CC: ptrString("cc@test.com")},
			want: 1,
		},
		{
			name: "with bcc",
			opts: &ListOptions{BCC: ptrString("bcc@test.com")},
			want: 1,
		},
		{
			name: "with unread",
			opts: &ListOptions{Unread: ptrBool(true)},
			want: 1,
		},
		{
			name: "with has_attachment",
			opts: &ListOptions{HasAttachment: ptrBool(true)},
			want: 1,
		},
		{
			name: "with all options",
			opts: &ListOptions{
				Limit:         ptr(100),
				PageToken:     "token123",
				Subject:       ptrString("Hello"),
				AnyEmail:      []string{"a@test.com"},
				To:            ptrString("to@test.com"),
				CC:            ptrString("cc@test.com"),
				BCC:           ptrString("bcc@test.com"),
				Unread:        ptrBool(false),
				Starred:       ptrBool(true),
				ThreadID:      ptrString("thread-1"),
				HasAttachment: ptrBool(true),
			},
			want: 11,
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

func TestDraft_DateTime(t *testing.T) {
	draft := &Draft{
		Date: 1704067200, // 2024-01-01 00:00:00 UTC
	}

	dt := draft.DateTime()
	expected := time.Unix(1704067200, 0)

	if !dt.Equal(expected) {
		t.Errorf("DateTime() = %v, want %v", dt, expected)
	}
}

func TestDraft_CreatedDateTime(t *testing.T) {
	draft := &Draft{
		CreatedAt: 1704067200,
	}

	dt := draft.CreatedDateTime()
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

func TestCreateRequest(t *testing.T) {
	req := &CreateRequest{
		Subject: "Test Subject",
		Body:    "Test Body",
		To:      []Participant{{Email: "to@example.com"}},
	}

	if req.Subject != "Test Subject" {
		t.Errorf("Subject = %s, want Test Subject", req.Subject)
	}
	if len(req.To) != 1 {
		t.Errorf("To count = %d, want 1", len(req.To))
	}
}

func TestTrackingOptions(t *testing.T) {
	opts := &TrackingOptions{
		Opens:         true,
		Links:         true,
		ThreadReplies: false,
		Label:         "campaign-1",
	}

	if !opts.Opens {
		t.Error("Opens = false, want true")
	}
	if opts.Label != "campaign-1" {
		t.Errorf("Label = %s, want campaign-1", opts.Label)
	}
}

// helpers
func ptr(v int) *int             { return &v }
func ptrBool(v bool) *bool       { return &v }
func ptrString(v string) *string { return &v }
