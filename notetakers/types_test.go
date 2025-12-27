package notetakers

import (
	"testing"
	"time"
)

func TestListOptions_Values(t *testing.T) {
	t.Run("nil options", func(t *testing.T) {
		var o *ListOptions
		if o.Values() != nil {
			t.Error("expected nil for nil options")
		}
	})

	t.Run("empty options", func(t *testing.T) {
		o := &ListOptions{}
		v := o.Values()
		if len(v) != 0 {
			t.Errorf("expected empty map, got %d entries", len(v))
		}
	})

	t.Run("with all options", func(t *testing.T) {
		limit := 25
		o := &ListOptions{
			Limit:     &limit,
			PageToken: "next-page-token",
			State:     StateScheduled,
		}
		v := o.Values()

		if v["limit"] != 25 {
			t.Errorf("expected limit=25, got %v", v["limit"])
		}
		if v["page_token"] != "next-page-token" {
			t.Errorf("expected page_token=next-page-token, got %v", v["page_token"])
		}
		if v["state"] != StateScheduled {
			t.Errorf("expected state=%s, got %v", StateScheduled, v["state"])
		}
	})

	t.Run("partial options - limit only", func(t *testing.T) {
		limit := 50
		o := &ListOptions{
			Limit: &limit,
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["limit"] != 50 {
			t.Errorf("expected limit=50, got %v", v["limit"])
		}
	})

	t.Run("partial options - state only", func(t *testing.T) {
		o := &ListOptions{
			State: StateRecording,
		}
		v := o.Values()

		if len(v) != 1 {
			t.Errorf("expected 1 entry, got %d", len(v))
		}
		if v["state"] != StateRecording {
			t.Errorf("expected state=%s, got %v", StateRecording, v["state"])
		}
	})
}

func TestNotetaker_TimeHelpers(t *testing.T) {
	// Use known Unix timestamps
	joinUnix := int64(1705315800)    // 2024-01-15 10:30:00 UTC
	createdUnix := int64(1705312200) // 2024-01-15 09:30:00 UTC (1 hour before)

	nt := &Notetaker{
		ID:        "notetaker123",
		JoinTime:  joinUnix,
		CreatedAt: createdUnix,
	}

	t.Run("JoinDateTime", func(t *testing.T) {
		join := nt.JoinDateTime()
		expected := time.Unix(joinUnix, 0)

		if !join.Equal(expected) {
			t.Errorf("JoinDateTime = %v, want %v", join, expected)
		}
	})

	t.Run("CreatedDateTime", func(t *testing.T) {
		created := nt.CreatedDateTime()
		expected := time.Unix(createdUnix, 0)

		if !created.Equal(expected) {
			t.Errorf("CreatedDateTime = %v, want %v", created, expected)
		}
	})

	t.Run("created before join", func(t *testing.T) {
		if !nt.CreatedDateTime().Before(nt.JoinDateTime()) {
			t.Error("expected CreatedDateTime to be before JoinDateTime")
		}
	})
}

func TestNotetaker_States(t *testing.T) {
	states := []string{
		StateScheduled,
		StateConnecting,
		StateWaiting,
		StateJoined,
		StateRecording,
		StateCompleted,
		StateCancelled,
		StateFailed,
	}

	expectedStates := map[string]bool{
		"scheduled":         true,
		"connecting":        true,
		"waiting_for_entry": true,
		"joined":            true,
		"recording":         true,
		"completed":         true,
		"cancelled":         true,
		"failed":            true,
	}

	for _, state := range states {
		if !expectedStates[state] {
			t.Errorf("unexpected state value: %s", state)
		}
	}

	if len(states) != len(expectedStates) {
		t.Errorf("expected %d states, got %d", len(expectedStates), len(states))
	}
}

func TestMediaTypes(t *testing.T) {
	mediaTypes := []string{
		MediaTypeVideo,
		MediaTypeAudio,
		MediaTypeTranscript,
		MediaTypeSummary,
		MediaTypeActionItems,
	}

	expectedTypes := map[string]bool{
		"video":        true,
		"audio":        true,
		"transcript":   true,
		"summary":      true,
		"action_items": true,
	}

	for _, mt := range mediaTypes {
		if !expectedTypes[mt] {
			t.Errorf("unexpected media type: %s", mt)
		}
	}

	if len(mediaTypes) != len(expectedTypes) {
		t.Errorf("expected %d media types, got %d", len(expectedTypes), len(mediaTypes))
	}
}

func TestNotetaker_Structure(t *testing.T) {
	nt := Notetaker{
		ID:              "notetaker123",
		Name:            "Meeting Notetaker",
		MeetingLink:     "https://zoom.us/j/123456789",
		MeetingProvider: "zoom",
		State:           StateScheduled,
		MeetingSettings: &MeetingSettings{
			VideoRecording:    true,
			AudioRecording:    true,
			Transcription:     true,
			Summary:           true,
			ActionItems:       true,
			LeaveAfterSilence: 300,
		},
	}

	if nt.ID != "notetaker123" {
		t.Errorf("expected ID=notetaker123, got %s", nt.ID)
	}
	if nt.MeetingProvider != "zoom" {
		t.Errorf("expected MeetingProvider=zoom, got %s", nt.MeetingProvider)
	}
	if nt.State != StateScheduled {
		t.Errorf("expected State=%s, got %s", StateScheduled, nt.State)
	}
	if nt.MeetingSettings == nil {
		t.Error("expected MeetingSettings to be set")
		return
	}
	if !nt.MeetingSettings.VideoRecording {
		t.Error("expected VideoRecording=true")
	}
	if nt.MeetingSettings.LeaveAfterSilence != 300 {
		t.Errorf("expected LeaveAfterSilence=300, got %d", nt.MeetingSettings.LeaveAfterSilence)
	}
}

func TestCreateRequest_Structure(t *testing.T) {
	joinTime := int64(1705315800)
	req := CreateRequest{
		MeetingLink: "https://meet.google.com/abc-defg-hij",
		JoinTime:    &joinTime,
		Name:        "Test Bot",
		MeetingSettings: &MeetingSettings{
			Transcription: true,
			Summary:       true,
		},
	}

	if req.MeetingLink != "https://meet.google.com/abc-defg-hij" {
		t.Errorf("expected MeetingLink, got %s", req.MeetingLink)
	}
	if req.JoinTime == nil || *req.JoinTime != joinTime {
		t.Error("expected JoinTime to be set")
	}
	if req.MeetingSettings == nil || !req.MeetingSettings.Transcription {
		t.Error("expected MeetingSettings.Transcription=true")
	}
}

func TestMedia_Structure(t *testing.T) {
	media := Media{
		ID:          "media123",
		Type:        MediaTypeTranscript,
		URL:         "https://storage.nylas.com/transcript/123",
		Status:      "ready",
		ExpiresAt:   1705401600,
		ContentType: "application/json",
	}

	if media.ID != "media123" {
		t.Errorf("expected ID=media123, got %s", media.ID)
	}
	if media.Type != MediaTypeTranscript {
		t.Errorf("expected Type=%s, got %s", MediaTypeTranscript, media.Type)
	}
	if media.Status != "ready" {
		t.Errorf("expected Status=ready, got %s", media.Status)
	}
}

func TestHistoryEvent_Structure(t *testing.T) {
	event := HistoryEvent{
		CreatedAt: 1705315800,
		EventType: "joined",
		Data: map[string]any{
			"participant_count": 5,
		},
	}

	if event.EventType != "joined" {
		t.Errorf("expected EventType=joined, got %s", event.EventType)
	}
	if event.Data["participant_count"] != 5 {
		t.Errorf("expected participant_count=5, got %v", event.Data["participant_count"])
	}
}
