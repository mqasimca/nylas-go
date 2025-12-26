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
func ptr(v int) *int       { return &v }
func ptrBool(v bool) *bool { return &v }
