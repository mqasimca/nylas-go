package nylas

import (
	"net/url"
	"testing"
)

func TestSetQueryParams(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		want   map[string][]string
	}{
		{
			name:   "empty params",
			params: map[string]any{},
			want:   map[string][]string{},
		},
		{
			name:   "string param",
			params: map[string]any{"subject": "test"},
			want:   map[string][]string{"subject": {"test"}},
		},
		{
			name:   "int param",
			params: map[string]any{"limit": 10},
			want:   map[string][]string{"limit": {"10"}},
		},
		{
			name:   "int64 param",
			params: map[string]any{"received_after": int64(1704067200)},
			want:   map[string][]string{"received_after": {"1704067200"}},
		},
		{
			name:   "bool param",
			params: map[string]any{"unread": true},
			want:   map[string][]string{"unread": {"true"}},
		},
		{
			name:   "string slice param",
			params: map[string]any{"any_email": []string{"a@test.com", "b@test.com"}},
			want:   map[string][]string{"any_email": {"a@test.com", "b@test.com"}},
		},
		{
			name: "mixed params",
			params: map[string]any{
				"limit":   50,
				"unread":  true,
				"subject": "hello",
			},
			want: map[string][]string{
				"limit":   {"50"},
				"unread":  {"true"},
				"subject": {"hello"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			setQueryParams(q, tt.params)

			for k, expected := range tt.want {
				got := q[k]
				if len(got) != len(expected) {
					t.Errorf("param %s: got %v, want %v", k, got, expected)
					continue
				}
				for i, v := range expected {
					if got[i] != v {
						t.Errorf("param %s[%d]: got %s, want %s", k, i, got[i], v)
					}
				}
			}
		})
	}
}
