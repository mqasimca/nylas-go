package events

import "time"

// Event represents a calendar event in the Nylas API.
type Event struct {
	// ID is the unique identifier for this event.
	ID string `json:"id"`
	// GrantID is the ID of the grant (connected account) this event belongs to.
	GrantID string `json:"grant_id"`
	// CalendarID is the ID of the calendar containing this event.
	CalendarID string `json:"calendar_id"`
	// Title is the event title/summary.
	Title string `json:"title,omitempty"`
	// Description is the event description (may contain HTML).
	Description string `json:"description,omitempty"`
	// Location is the event location (physical or virtual).
	Location string `json:"location,omitempty"`
	// When contains the event timing information.
	When When `json:"when"`
	// Participants is the list of attendees.
	Participants []Participant `json:"participants,omitempty"`
	// Status is the event status ("confirmed", "tentative", "cancelled").
	Status string `json:"status,omitempty"`
	// Busy indicates whether the event blocks time on the calendar.
	Busy bool `json:"busy,omitempty"`
	// Visibility is "public", "private", or "default".
	Visibility string `json:"visibility,omitempty"`
	// Conferencing contains video conferencing details.
	Conferencing *Conferencing `json:"conferencing,omitempty"`
	// Reminders contains notification settings.
	Reminders *Reminders `json:"reminders,omitempty"`
	// Recurrence contains rules for recurring events.
	Recurrence *Recurrence `json:"recurrence,omitempty"`
	// Metadata contains custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Capacity is the maximum number of attendees.
	Capacity *int `json:"capacity,omitempty"`
	// HideParticipants hides the attendee list from other participants.
	HideParticipants bool `json:"hide_participants,omitempty"`
	// Organizer is the event organizer.
	Organizer *Organizer `json:"organizer,omitempty"`
	// Creator is the user who created the event.
	Creator *Organizer `json:"creator,omitempty"`
	// Object is the object type, always "event".
	Object string `json:"object,omitempty"`
	// MasterEventID is the ID of the master event for recurring instances.
	MasterEventID string `json:"master_event_id,omitempty"`
	// OriginalStartTime is the original start time for modified recurring instances.
	OriginalStartTime *int64 `json:"original_start_time,omitempty"`
	// ICalUID is the iCalendar UID.
	ICalUID string `json:"ical_uid,omitempty"`
	// HTMLLink is the URL to view the event in the provider's web interface.
	HTMLLink string `json:"html_link,omitempty"`
	// ReadOnly indicates whether the event can be modified.
	ReadOnly bool `json:"read_only,omitempty"`
	// CreatedAt is the Unix timestamp when the event was created.
	CreatedAt int64 `json:"created_at,omitempty"`
	// UpdatedAt is the Unix timestamp when the event was last modified.
	UpdatedAt int64 `json:"updated_at,omitempty"`
	// Occurrences lists expanded instances of a recurring event.
	Occurrences []string `json:"occurrences,omitempty"`
	// CancelledOccurrences lists cancelled instances of a recurring event.
	CancelledOccurrences []string `json:"cancelled_occurrences,omitempty"`
	// Resources lists rooms and equipment booked for the event.
	Resources []Resource `json:"resources,omitempty"`
	// ColorID is the provider-specific color identifier.
	ColorID string `json:"color_id,omitempty"`
}

// When represents the time information for an event.
// It can be a timespan (start_time/end_time), datespan (start_date/end_date),
// single date, or single time depending on the Object field.
type When struct {
	// Object is "timespan", "datespan", "date", or "time".
	Object string `json:"object,omitempty"`
	// StartTime is the start time as Unix timestamp (for timespan).
	StartTime *int64 `json:"start_time,omitempty"`
	// EndTime is the end time as Unix timestamp (for timespan).
	EndTime *int64 `json:"end_time,omitempty"`
	// StartTimezone is the IANA timezone for StartTime.
	StartTimezone string `json:"start_timezone,omitempty"`
	// EndTimezone is the IANA timezone for EndTime.
	EndTimezone string `json:"end_timezone,omitempty"`
	// Date is a single all-day date in YYYY-MM-DD format.
	Date string `json:"date,omitempty"`
	// StartDate is the start of a multi-day event in YYYY-MM-DD format.
	StartDate string `json:"start_date,omitempty"`
	// EndDate is the end of a multi-day event in YYYY-MM-DD format.
	EndDate string `json:"end_date,omitempty"`
	// Time is a single point in time as Unix timestamp.
	Time *int64 `json:"time,omitempty"`
	// Timezone is the IANA timezone for Date or Time fields.
	Timezone string `json:"timezone,omitempty"`
}

// Participant represents an event participant/attendee.
type Participant struct {
	// Name is the participant's display name.
	Name string `json:"name,omitempty"`
	// Email is the participant's email address.
	Email string `json:"email"`
	// Status is the RSVP status ("yes", "no", "maybe", "noreply").
	Status string `json:"status,omitempty"`
	// Comment is an optional comment from the participant.
	Comment string `json:"comment,omitempty"`
	// PhoneNumber is the participant's phone number (Microsoft only).
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Organizer represents the event organizer or creator.
type Organizer struct {
	// Name is the organizer's display name.
	Name string `json:"name,omitempty"`
	// Email is the organizer's email address.
	Email string `json:"email"`
}

// Conferencing represents video conferencing details for an event.
type Conferencing struct {
	// Provider is "Google Meet", "Zoom Meeting", "Microsoft Teams", etc.
	Provider string `json:"provider,omitempty"`
	// Details contains the meeting URL and access codes.
	Details *ConferencingDetails `json:"details,omitempty"`
	// Autocreate configures automatic conference creation.
	Autocreate *AutocreateConfig `json:"autocreate,omitempty"`
}

// ConferencingDetails contains the video meeting access information.
type ConferencingDetails struct {
	// URL is the meeting join URL.
	URL string `json:"url,omitempty"`
	// MeetingCode is the meeting ID/code.
	MeetingCode string `json:"meeting_code,omitempty"`
	// Password is the meeting password if required.
	Password string `json:"password,omitempty"`
	// Phone contains dial-in phone numbers.
	Phone []string `json:"phone,omitempty"`
	// PIN is the dial-in PIN code.
	PIN string `json:"pin,omitempty"`
}

// AutocreateConfig configures automatic conference creation when the event is created.
type AutocreateConfig struct {
	// ConfGrantID is the grant ID for the conferencing provider.
	ConfGrantID string `json:"conf_grant_id"`
	// ConfSettings contains provider-specific settings.
	ConfSettings map[string]any `json:"conf_settings,omitempty"`
}

// Reminders represents reminder/notification settings for an event.
type Reminders struct {
	// UseDefault uses the calendar's default reminder settings.
	UseDefault bool `json:"use_default,omitempty"`
	// Overrides is a list of custom reminder times.
	Overrides []ReminderOverride `json:"overrides,omitempty"`
}

// ReminderOverride represents a custom reminder notification.
type ReminderOverride struct {
	// ReminderMinutes is how many minutes before the event to send the reminder.
	ReminderMinutes int `json:"reminder_minutes"`
	// ReminderMethod is "email", "popup", "display", or "sound".
	ReminderMethod string `json:"reminder_method,omitempty"`
}

// Recurrence represents recurrence rules for a recurring event.
type Recurrence struct {
	// RRule is the iCalendar RRULE string (e.g., "RRULE:FREQ=WEEKLY;BYDAY=MO").
	RRule string `json:"rrule,omitempty"`
	// Exdate is the iCalendar EXDATE string for excluded dates.
	Exdate string `json:"exdate,omitempty"`
}

// Resource represents a bookable resource (room, equipment) for an event.
type Resource struct {
	// Email is the resource's booking email address.
	Email string `json:"email"`
	// Name is the resource's display name.
	Name string `json:"name,omitempty"`
	// Building is the building where the resource is located.
	Building string `json:"building,omitempty"`
	// Capacity is the room capacity.
	Capacity *int `json:"capacity,omitempty"`
	// FloorName is the floor name.
	FloorName string `json:"floor_name,omitempty"`
	// FloorNumber is the floor number.
	FloorNumber *int `json:"floor_number,omitempty"`
	// FloorSection is the section of the floor.
	FloorSection string `json:"floor_section,omitempty"`
	// Object is the object type, always "resource".
	Object string `json:"object,omitempty"`
}

// ListOptions specifies options for listing events.
// All fields are optional; nil values are not included in the request.
type ListOptions struct {
	// Limit is the maximum number of events to return (default 50, max 200).
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
	// CalendarID filters events by calendar (required if not in path).
	CalendarID string `json:"calendar_id,omitempty"`
	// Start filters events starting on or after this Unix timestamp.
	Start *int64 `json:"start,omitempty"`
	// End filters events starting before this Unix timestamp.
	End *int64 `json:"end,omitempty"`
	// ExpandRecurring expands recurring events into individual instances.
	ExpandRecurring *bool `json:"expand_recurring,omitempty"`
	// ShowCancelled includes cancelled events in results.
	ShowCancelled *bool `json:"show_cancelled,omitempty"`
	// Busy filters by busy/free status.
	Busy *bool `json:"busy,omitempty"`
	// Title filters events by title (substring match).
	Title *string `json:"title,omitempty"`
	// Description filters events by description (substring match).
	Description *string `json:"description,omitempty"`
	// Location filters events by location (substring match).
	Location *string `json:"location,omitempty"`
	// Attendees filters events by attendee email (comma-separated).
	Attendees *string `json:"attendees,omitempty"`
	// MasterEventID filters to instances of a specific recurring event.
	MasterEventID *string `json:"master_event_id,omitempty"`
	// ICalUID filters events by iCalendar UID.
	ICalUID *string `json:"ical_uid,omitempty"`
	// UpdatedAfter filters events updated after this Unix timestamp.
	UpdatedAfter *int64 `json:"updated_after,omitempty"`
	// UpdatedBefore filters events updated before this Unix timestamp.
	UpdatedBefore *int64 `json:"updated_before,omitempty"`
	// MetadataPair filters by metadata key-value pair ("key:value" format).
	MetadataPair *string `json:"metadata_pair,omitempty"`
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
	if o.CalendarID != "" {
		v["calendar_id"] = o.CalendarID
	}
	if o.Start != nil {
		v["start"] = *o.Start
	}
	if o.End != nil {
		v["end"] = *o.End
	}
	if o.ExpandRecurring != nil {
		v["expand_recurring"] = *o.ExpandRecurring
	}
	if o.ShowCancelled != nil {
		v["show_cancelled"] = *o.ShowCancelled
	}
	if o.Busy != nil {
		v["busy"] = *o.Busy
	}
	if o.Title != nil {
		v["title"] = *o.Title
	}
	if o.Description != nil {
		v["description"] = *o.Description
	}
	if o.Location != nil {
		v["location"] = *o.Location
	}
	if o.Attendees != nil {
		v["attendees"] = *o.Attendees
	}
	if o.MasterEventID != nil {
		v["master_event_id"] = *o.MasterEventID
	}
	if o.ICalUID != nil {
		v["ical_uid"] = *o.ICalUID
	}
	if o.UpdatedAfter != nil {
		v["updated_after"] = *o.UpdatedAfter
	}
	if o.UpdatedBefore != nil {
		v["updated_before"] = *o.UpdatedBefore
	}
	if o.MetadataPair != nil {
		v["metadata_pair"] = *o.MetadataPair
	}
	return v
}

// CreateRequest represents a request to create an event.
type CreateRequest struct {
	// Title is the event title/summary.
	Title string `json:"title,omitempty"`
	// Description is the event description.
	Description string `json:"description,omitempty"`
	// Location is the event location.
	Location string `json:"location,omitempty"`
	// When contains the event timing (required).
	When When `json:"when"`
	// Participants is the list of attendees to invite.
	Participants []Participant `json:"participants,omitempty"`
	// Busy marks the event as blocking time (default true).
	Busy *bool `json:"busy,omitempty"`
	// Visibility is "public", "private", or "default".
	Visibility string `json:"visibility,omitempty"`
	// Conferencing adds video conferencing to the event.
	Conferencing *Conferencing `json:"conferencing,omitempty"`
	// Reminders sets notification preferences.
	Reminders *Reminders `json:"reminders,omitempty"`
	// Recurrence makes this a recurring event (RRULE strings).
	Recurrence []string `json:"recurrence,omitempty"`
	// Metadata contains custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Capacity sets the maximum number of attendees.
	Capacity *int `json:"capacity,omitempty"`
	// HideParticipants hides the attendee list.
	HideParticipants *bool `json:"hide_participants,omitempty"`
	// Resources books rooms and equipment.
	Resources []Resource `json:"resources,omitempty"`
	// ColorID sets a provider-specific color.
	ColorID string `json:"color_id,omitempty"`
}

// UpdateRequest represents a request to update an event.
type UpdateRequest struct {
	// Title is the event title/summary.
	Title *string `json:"title,omitempty"`
	// Description is the event description.
	Description *string `json:"description,omitempty"`
	// Location is the event location.
	Location *string `json:"location,omitempty"`
	// When updates the event timing.
	When *When `json:"when,omitempty"`
	// Participants replaces the attendee list.
	Participants []Participant `json:"participants,omitempty"`
	// Busy marks the event as blocking time.
	Busy *bool `json:"busy,omitempty"`
	// Visibility is "public", "private", or "default".
	Visibility *string `json:"visibility,omitempty"`
	// Conferencing updates video conferencing settings.
	Conferencing *Conferencing `json:"conferencing,omitempty"`
	// Reminders updates notification preferences.
	Reminders *Reminders `json:"reminders,omitempty"`
	// Recurrence updates recurrence rules.
	Recurrence []string `json:"recurrence,omitempty"`
	// Metadata updates custom key-value pairs.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Capacity sets the maximum number of attendees.
	Capacity *int `json:"capacity,omitempty"`
	// HideParticipants hides the attendee list.
	HideParticipants *bool `json:"hide_participants,omitempty"`
	// Resources updates room and equipment bookings.
	Resources []Resource `json:"resources,omitempty"`
	// ColorID sets a provider-specific color.
	ColorID *string `json:"color_id,omitempty"`
}

// RSVPRequest represents a request to send an RSVP response to an event invitation.
type RSVPRequest struct {
	// Status is the RSVP response ("yes", "no", "maybe").
	Status string `json:"status"`
	// Comment is an optional message to include with the RSVP.
	Comment string `json:"comment,omitempty"`
}

// StartDateTime returns the event start time as time.Time.
func (e *Event) StartDateTime() time.Time {
	if e.When.StartTime != nil {
		return time.Unix(*e.When.StartTime, 0)
	}
	if e.When.Time != nil {
		return time.Unix(*e.When.Time, 0)
	}
	return time.Time{}
}

// EndDateTime returns the event end time as time.Time.
func (e *Event) EndDateTime() time.Time {
	if e.When.EndTime != nil {
		return time.Unix(*e.When.EndTime, 0)
	}
	return time.Time{}
}

// CreatedDateTime returns the created timestamp as time.Time.
func (e *Event) CreatedDateTime() time.Time {
	return time.Unix(e.CreatedAt, 0)
}

// UpdatedDateTime returns the updated timestamp as time.Time.
func (e *Event) UpdatedDateTime() time.Time {
	return time.Unix(e.UpdatedAt, 0)
}

// IsAllDay returns true if the event is an all-day event.
func (e *Event) IsAllDay() bool {
	return e.When.Object == "date" || e.When.Object == "datespan"
}

// IsRecurring returns true if the event is recurring.
func (e *Event) IsRecurring() bool {
	return e.Recurrence != nil && e.Recurrence.RRule != ""
}

// ImportOptions specifies options for importing events.
// This endpoint returns all events from a calendar, including recurring event
// instances with their parent events and any overrides.
type ImportOptions struct {
	// CalendarID is the calendar to import events from (required).
	CalendarID string `json:"calendar_id"`
	// Start filters events starting on or after this Unix timestamp.
	Start *int64 `json:"start,omitempty"`
	// End filters events starting before this Unix timestamp.
	End *int64 `json:"end,omitempty"`
	// Limit is the maximum number of events to return (default 50, max 200).
	Limit *int `json:"limit,omitempty"`
	// PageToken is the cursor for pagination.
	PageToken string `json:"page_token,omitempty"`
}

// Values converts ImportOptions to URL query parameters.
func (o *ImportOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	v := make(map[string]any)
	if o.CalendarID != "" {
		v["calendar_id"] = o.CalendarID
	}
	if o.Start != nil {
		v["start"] = *o.Start
	}
	if o.End != nil {
		v["end"] = *o.End
	}
	if o.Limit != nil {
		v["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		v["page_token"] = o.PageToken
	}
	return v
}
