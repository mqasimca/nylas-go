package notetakers

import "time"

// Notetaker represents a Nylas Notetaker bot.
type Notetaker struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	JoinTime        int64            `json:"join_time,omitempty"`
	MeetingLink     string           `json:"meeting_link,omitempty"`
	MeetingProvider string           `json:"meeting_provider,omitempty"`
	State           string           `json:"state,omitempty"`
	MeetingSettings *MeetingSettings `json:"meeting_settings,omitempty"`
	CreatedAt       int64            `json:"created_at,omitempty"`
	UpdatedAt       int64            `json:"updated_at,omitempty"`
}

// MeetingSettings defines settings for a notetaker session.
type MeetingSettings struct {
	VideoRecording    bool `json:"video_recording,omitempty"`
	AudioRecording    bool `json:"audio_recording,omitempty"`
	Transcription     bool `json:"transcription,omitempty"`
	Summary           bool `json:"summary,omitempty"`
	ActionItems       bool `json:"action_items,omitempty"`
	LeaveAfterSilence int  `json:"leave_after_silence_seconds,omitempty"`
}

// CreateRequest represents a request to create/invite a notetaker.
type CreateRequest struct {
	MeetingLink     string           `json:"meeting_link"`
	JoinTime        *int64           `json:"join_time,omitempty"`
	Name            string           `json:"name,omitempty"`
	MeetingSettings *MeetingSettings `json:"meeting_settings,omitempty"`
}

// HistoryEvent represents an event in a notetaker's history.
type HistoryEvent struct {
	CreatedAt int64          `json:"created_at,omitempty"`
	EventType string         `json:"event_type,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
}

// History represents the history of a notetaker.
type History struct {
	Events []HistoryEvent `json:"events,omitempty"`
}

// Media represents media files from a notetaker session.
type Media struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	URL         string `json:"url,omitempty"`
	Status      string `json:"status,omitempty"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

// ListOptions specifies options for listing notetakers.
type ListOptions struct {
	Limit     *int   `json:"limit,omitempty"`
	PageToken string `json:"page_token,omitempty"`
	State     string `json:"state,omitempty"`
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
	if o.State != "" {
		v["state"] = o.State
	}
	return v
}

// JoinDateTime returns the join time as time.Time.
func (n *Notetaker) JoinDateTime() time.Time {
	return time.Unix(n.JoinTime, 0)
}

// CreatedDateTime returns the created time as time.Time.
func (n *Notetaker) CreatedDateTime() time.Time {
	return time.Unix(n.CreatedAt, 0)
}

// Notetaker states.
const (
	StateScheduled  = "scheduled"
	StateConnecting = "connecting"
	StateWaiting    = "waiting_for_entry"
	StateJoined     = "joined"
	StateRecording  = "recording"
	StateCompleted  = "completed"
	StateCancelled  = "cancelled"
	StateFailed     = "failed"
)

// Media types.
const (
	MediaTypeVideo       = "video"
	MediaTypeAudio       = "audio"
	MediaTypeTranscript  = "transcript"
	MediaTypeSummary     = "summary"
	MediaTypeActionItems = "action_items"
)
