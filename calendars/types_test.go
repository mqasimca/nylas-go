package calendars

import "testing"

func TestListOptions_Values(t *testing.T) {
	tests := []struct {
		name string
		opts *ListOptions
		want map[string]any
	}{
		{
			name: "nil options",
			opts: nil,
			want: nil,
		},
		{
			name: "empty options",
			opts: &ListOptions{},
			want: map[string]any{},
		},
		{
			name: "with limit",
			opts: &ListOptions{Limit: ptr(10)},
			want: map[string]any{"limit": 10},
		},
		{
			name: "with page token",
			opts: &ListOptions{PageToken: "token123"},
			want: map[string]any{"page_token": "token123"},
		},
		{
			name: "all options",
			opts: &ListOptions{Limit: ptr(25), PageToken: "abc"},
			want: map[string]any{"limit": 25, "page_token": "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Values()
			if tt.want == nil {
				if got != nil {
					t.Errorf("Values() = %v, want nil", got)
				}
				return
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("Values()[%s] = %v, want %v", k, got[k], v)
				}
			}
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
