package events

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
			name: "with calendar_id",
			opts: &ListOptions{CalendarID: "cal-123"},
			want: map[string]any{"calendar_id": "cal-123"},
		},
		{
			name: "with time range",
			opts: &ListOptions{
				Start: ptr(int64(1700000000)),
				End:   ptr(int64(1700100000)),
			},
			want: map[string]any{"start": int64(1700000000), "end": int64(1700100000)},
		},
		{
			name: "with expand_recurring",
			opts: &ListOptions{ExpandRecurring: ptr(true)},
			want: map[string]any{"expand_recurring": true},
		},
		{
			name: "with show_cancelled",
			opts: &ListOptions{ShowCancelled: ptr(true)},
			want: map[string]any{"show_cancelled": true},
		},
		{
			name: "all common options",
			opts: &ListOptions{
				Limit:           ptr(50),
				PageToken:       "token123",
				CalendarID:      "cal-456",
				Start:           ptr(int64(1700000000)),
				End:             ptr(int64(1700100000)),
				ExpandRecurring: ptr(true),
				ShowCancelled:   ptr(false),
				Title:           ptr("Meeting"),
			},
			want: map[string]any{
				"limit":            50,
				"page_token":       "token123",
				"calendar_id":      "cal-456",
				"start":            int64(1700000000),
				"end":              int64(1700100000),
				"expand_recurring": true,
				"show_cancelled":   false,
				"title":            "Meeting",
			},
		},
		{
			name: "with description",
			opts: &ListOptions{Description: ptr("team sync")},
			want: map[string]any{"description": "team sync"},
		},
		{
			name: "with location",
			opts: &ListOptions{Location: ptr("Conference Room A")},
			want: map[string]any{"location": "Conference Room A"},
		},
		{
			name: "with attendees",
			opts: &ListOptions{Attendees: ptr("user@example.com")},
			want: map[string]any{"attendees": "user@example.com"},
		},
		{
			name: "with busy filter",
			opts: &ListOptions{Busy: ptr(true)},
			want: map[string]any{"busy": true},
		},
		{
			name: "with master_event_id",
			opts: &ListOptions{MasterEventID: ptr("master-123")},
			want: map[string]any{"master_event_id": "master-123"},
		},
		{
			name: "with ical_uid",
			opts: &ListOptions{ICalUID: ptr("uid@example.com")},
			want: map[string]any{"ical_uid": "uid@example.com"},
		},
		{
			name: "with updated_after",
			opts: &ListOptions{UpdatedAfter: ptr(int64(1699900000))},
			want: map[string]any{"updated_after": int64(1699900000)},
		},
		{
			name: "with updated_before",
			opts: &ListOptions{UpdatedBefore: ptr(int64(1700200000))},
			want: map[string]any{"updated_before": int64(1700200000)},
		},
		{
			name: "with metadata_pair",
			opts: &ListOptions{MetadataPair: ptr("key:value")},
			want: map[string]any{"metadata_pair": "key:value"},
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

func TestEvent_IsAllDay(t *testing.T) {
	tests := []struct {
		name  string
		event *Event
		want  bool
	}{
		{
			name:  "timespan event",
			event: &Event{When: When{Object: "timespan"}},
			want:  false,
		},
		{
			name:  "date event",
			event: &Event{When: When{Object: "date"}},
			want:  true,
		},
		{
			name:  "datespan event",
			event: &Event{When: When{Object: "datespan"}},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsAllDay(); got != tt.want {
				t.Errorf("IsAllDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_IsRecurring(t *testing.T) {
	tests := []struct {
		name  string
		event *Event
		want  bool
	}{
		{
			name:  "non-recurring",
			event: &Event{},
			want:  false,
		},
		{
			name:  "nil recurrence",
			event: &Event{Recurrence: nil},
			want:  false,
		},
		{
			name:  "empty rrule",
			event: &Event{Recurrence: &Recurrence{}},
			want:  false,
		},
		{
			name:  "with rrule",
			event: &Event{Recurrence: &Recurrence{RRule: "RRULE:FREQ=WEEKLY"}},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.event.IsRecurring(); got != tt.want {
				t.Errorf("IsRecurring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvent_StartDateTime(t *testing.T) {
	startTime := int64(1700000000)
	event := &Event{
		When: When{
			StartTime: &startTime,
		},
	}

	if got := event.StartDateTime().Unix(); got != startTime {
		t.Errorf("StartDateTime() = %d, want %d", got, startTime)
	}

	// Test with time field
	timeVal := int64(1700000000)
	eventWithTime := &Event{
		When: When{
			Time: &timeVal,
		},
	}
	if got := eventWithTime.StartDateTime().Unix(); got != timeVal {
		t.Errorf("StartDateTime() with time = %d, want %d", got, timeVal)
	}

	// Test with no time
	eventNoTime := &Event{When: When{}}
	if !eventNoTime.StartDateTime().IsZero() {
		t.Error("StartDateTime() with no time should be zero")
	}
}

func TestEvent_EndDateTime(t *testing.T) {
	endTime := int64(1700003600)
	event := &Event{
		When: When{
			EndTime: &endTime,
		},
	}

	if got := event.EndDateTime().Unix(); got != endTime {
		t.Errorf("EndDateTime() = %d, want %d", got, endTime)
	}

	// Test with no end time
	eventNoEnd := &Event{When: When{}}
	if !eventNoEnd.EndDateTime().IsZero() {
		t.Error("EndDateTime() with no end time should be zero")
	}
}

func TestEvent_CreatedDateTime(t *testing.T) {
	createdAt := int64(1699999999)
	event := &Event{
		CreatedAt: createdAt,
	}

	if got := event.CreatedDateTime().Unix(); got != createdAt {
		t.Errorf("CreatedDateTime() = %d, want %d", got, createdAt)
	}
}

func TestEvent_UpdatedDateTime(t *testing.T) {
	updatedAt := int64(1700000001)
	event := &Event{
		UpdatedAt: updatedAt,
	}

	if got := event.UpdatedDateTime().Unix(); got != updatedAt {
		t.Errorf("UpdatedDateTime() = %d, want %d", got, updatedAt)
	}
}

func ptr[T any](v T) *T {
	return &v
}
