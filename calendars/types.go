package calendars

// Calendar represents a calendar in the Nylas API.
type Calendar struct {
	// ID is the unique identifier for this calendar.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this calendar belongs to.
	GrantID string `json:"grant_id"`
	// Name is the display name of the calendar.
	Name string `json:"name"`
	// Description is the calendar description.
	Description string `json:"description,omitempty"`
	// Location is the geographic location of the calendar.
	Location string `json:"location,omitempty"`
	// Timezone is the IANA timezone (e.g., "America/New_York").
	Timezone string `json:"timezone,omitempty"`
	// IsPrimary indicates whether this is the user's primary calendar.
	IsPrimary bool `json:"is_primary,omitempty"`
	// ReadOnly indicates whether the calendar is read-only.
	ReadOnly bool `json:"read_only,omitempty"`
	// IsOwnedByUser indicates whether the user owns this calendar.
	IsOwnedByUser bool `json:"is_owned_by_user,omitempty"`
	// Object is the object type, always "calendar".
	Object string `json:"object,omitempty"`
	// HexColor is the background color in hex format (e.g., "#0099EE").
	HexColor string `json:"hex_color,omitempty"`
	// HexForegroundColor is the text color in hex format.
	HexForegroundColor string `json:"hex_foreground_color,omitempty"`
	// Metadata contains custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListOptions specifies options for listing calendars.
type ListOptions struct {
	// Limit is the maximum number of calendars to return.
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
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
	return v
}

// CreateRequest represents a request to create a calendar.
type CreateRequest struct {
	// Name is the display name of the calendar (required).
	Name string `json:"name"`
	// Description is the calendar description.
	Description string `json:"description,omitempty"`
	// Location is the geographic location of the calendar.
	Location string `json:"location,omitempty"`
	// Timezone is the IANA timezone (e.g., "America/New_York").
	Timezone string `json:"timezone,omitempty"`
	// Metadata contains custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateRequest represents a request to update a calendar.
type UpdateRequest struct {
	// Name is the display name of the calendar.
	Name *string `json:"name,omitempty"`
	// Description is the calendar description.
	Description *string `json:"description,omitempty"`
	// Location is the geographic location of the calendar.
	Location *string `json:"location,omitempty"`
	// Timezone is the IANA timezone.
	Timezone *string `json:"timezone,omitempty"`
	// HexColor is the background color in hex format.
	HexColor *string `json:"hex_color,omitempty"`
	// Metadata contains custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// AvailabilityRequest represents a request to check availability across participants.
type AvailabilityRequest struct {
	// StartTime is the start of the availability window (Unix timestamp).
	StartTime int64 `json:"start_time"`
	// EndTime is the end of the availability window (Unix timestamp).
	EndTime int64 `json:"end_time"`
	// DurationMinutes is the required meeting duration.
	DurationMinutes int `json:"duration_minutes"`
	// IntervalMinutes is the interval between available slots (default: DurationMinutes).
	IntervalMinutes *int `json:"interval_minutes,omitempty"`
	// RoundTo30Minutes rounds time slots to 30-minute boundaries.
	RoundTo30Minutes bool `json:"round_to_30_minutes,omitempty"`
	// Participants is the list of participants to check availability for.
	Participants []AvailabilityParticipant `json:"participants"`
	// AvailabilityRules specifies additional rules for availability calculation.
	AvailabilityRules *AvailabilityRules `json:"availability_rules,omitempty"`
}

// AvailabilityParticipant represents a participant in an availability check.
type AvailabilityParticipant struct {
	// Email is the participant's email address.
	Email string `json:"email"`
	// CalendarIDs limits availability check to specific calendars.
	CalendarIDs []string `json:"calendar_ids,omitempty"`
	// OpenHours defines when this participant is available.
	OpenHours []OpenHours `json:"open_hours,omitempty"`
}

// OpenHours represents recurring availability windows.
type OpenHours struct {
	// Days is the days of the week (0=Sunday, 6=Saturday).
	Days []int `json:"days"`
	// Timezone is the IANA timezone for these hours.
	Timezone string `json:"timezone"`
	// Start is the start time in HH:MM format (24-hour).
	Start string `json:"start"`
	// End is the end time in HH:MM format (24-hour).
	End string `json:"end"`
	// Exdates is a list of dates to exclude (YYYY-MM-DD format).
	Exdates []string `json:"exdates,omitempty"`
}

// AvailabilityRules specifies rules for calculating availability.
type AvailabilityRules struct {
	// AvailabilityMethod is "collective", "max-fairness", or "max-availability".
	AvailabilityMethod string `json:"availability_method,omitempty"`
	// Buffer specifies buffer time before and after meetings.
	Buffer *Buffer `json:"buffer,omitempty"`
	// DefaultOpenHours applies to all participants without specific open hours.
	DefaultOpenHours []OpenHours `json:"default_open_hours,omitempty"`
	// RoundRobinGroupID groups participants for round-robin scheduling.
	RoundRobinGroupID string `json:"round_robin_group_id,omitempty"`
}

// Buffer represents buffer time before and after meetings.
type Buffer struct {
	// Before is the buffer time in minutes before meetings.
	Before int `json:"before"`
	// After is the buffer time in minutes after meetings.
	After int `json:"after"`
}

// AvailabilityResponse contains available time slots from an availability check.
type AvailabilityResponse struct {
	// TimeSlots is the list of available time slots.
	TimeSlots []TimeSlot `json:"time_slots"`
}

// TimeSlot represents an available time slot.
type TimeSlot struct {
	// StartTime is the start of the slot (Unix timestamp).
	StartTime int64 `json:"start_time"`
	// EndTime is the end of the slot (Unix timestamp).
	EndTime int64 `json:"end_time"`
	// Emails lists the participants available during this slot.
	Emails []string `json:"emails,omitempty"`
}

// FreeBusyRequest represents a request to get free/busy information.
type FreeBusyRequest struct {
	// StartTime is the start of the query window (Unix timestamp).
	StartTime int64 `json:"start_time"`
	// EndTime is the end of the query window (Unix timestamp).
	EndTime int64 `json:"end_time"`
	// Emails is the list of email addresses to check.
	Emails []string `json:"emails"`
}

// FreeBusyResponse represents free/busy data for a single email address.
type FreeBusyResponse struct {
	// Email is the email address this data is for.
	Email string `json:"email"`
	// TimeSlots lists the busy periods.
	TimeSlots []BusySlot `json:"time_slots"`
}

// BusySlot represents a busy time period.
type BusySlot struct {
	// StartTime is the start of the busy period (Unix timestamp).
	StartTime int64 `json:"start_time"`
	// EndTime is the end of the busy period (Unix timestamp).
	EndTime int64 `json:"end_time"`
	// Status is the availability status (e.g., "busy", "tentative").
	Status string `json:"status,omitempty"`
}
