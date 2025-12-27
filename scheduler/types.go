package scheduler

import "time"

// Configuration represents a Scheduler configuration.
type Configuration struct {
	ID                  string              `json:"id,omitempty"`
	Participants        []Participant       `json:"participants,omitempty"`
	Availability        *Availability       `json:"availability,omitempty"`
	EventBooking        *EventBooking       `json:"event_booking,omitempty"`
	Scheduler           *SchedulerSettings  `json:"scheduler,omitempty"`
	AppearanceSettings  *AppearanceSettings `json:"appearance,omitempty"`
	RequiresSessionAuth bool                `json:"requires_session_auth,omitempty"`
}

// Participant represents a scheduling participant.
type Participant struct {
	Name         string              `json:"name,omitempty"`
	Email        string              `json:"email,omitempty"`
	IsOrganizer  bool                `json:"is_organizer,omitempty"`
	Availability *AvailabilityRules  `json:"availability,omitempty"`
	Booking      *ParticipantBooking `json:"booking,omitempty"`
	Timezone     string              `json:"timezone,omitempty"`
}

// AvailabilityRules defines availability for a participant.
type AvailabilityRules struct {
	CalendarIDs   []string       `json:"calendar_ids,omitempty"`
	OpenHours     []OpenHours    `json:"open_hours,omitempty"`
	BufferBefore  int            `json:"buffer_before,omitempty"`
	BufferAfter   int            `json:"buffer_after,omitempty"`
	RoundTo       int            `json:"round_to,omitempty"`
	AvailableDays []AvailableDay `json:"available_days,omitempty"`
}

// OpenHours defines open hours for scheduling.
type OpenHours struct {
	Days       []int  `json:"days,omitempty"`
	Timezone   string `json:"timezone,omitempty"`
	StartTime  string `json:"start,omitempty"`
	EndTime    string `json:"end,omitempty"`
	ObjectType string `json:"object_type,omitempty"`
}

// AvailableDay defines availability for a specific day.
type AvailableDay struct {
	Day   int        `json:"day,omitempty"`
	Hours []TimeSlot `json:"hours,omitempty"`
}

// TimeSlot represents a time slot.
type TimeSlot struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

// ParticipantBooking defines booking settings for a participant.
type ParticipantBooking struct {
	CalendarID string `json:"calendar_id,omitempty"`
}

// Availability defines overall availability settings.
type Availability struct {
	DurationMinutes   int                `json:"duration_minutes,omitempty"`
	IntervalMinutes   int                `json:"interval_minutes,omitempty"`
	AvailabilityRules *AvailabilityRules `json:"availability_rules,omitempty"`
}

// EventBooking defines event booking settings.
type EventBooking struct {
	Title                 string            `json:"title,omitempty"`
	Description           string            `json:"description,omitempty"`
	Location              string            `json:"location,omitempty"`
	Timezone              string            `json:"timezone,omitempty"`
	BookingType           string            `json:"booking_type,omitempty"`
	ConferencingProvider  string            `json:"conferencing,omitempty"`
	DisableEmails         bool              `json:"disable_emails,omitempty"`
	Reminders             []Reminder        `json:"reminders,omitempty"`
	AdditionalFields      []AdditionalField `json:"additional_fields,omitempty"`
	HideParticipants      bool              `json:"hide_participants,omitempty"`
	MinBookingNotice      int               `json:"min_booking_notice,omitempty"`
	MinCancellationNotice int               `json:"min_cancellation_notice,omitempty"`
}

// Reminder defines a booking reminder.
type Reminder struct {
	Type          string `json:"type,omitempty"`
	MinutesBefore int    `json:"minutes_before_event,omitempty"`
	WebhookURL    string `json:"webhook_url,omitempty"`
	Recipient     string `json:"recipient,omitempty"`
	EmailSubject  string `json:"email_subject,omitempty"`
}

// AdditionalField defines a custom field for bookings.
type AdditionalField struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Label    string `json:"label,omitempty"`
	Required bool   `json:"required,omitempty"`
}

// SchedulerSettings defines scheduler UI settings.
type SchedulerSettings struct {
	AvailableDaysInFuture    int    `json:"available_days_in_future,omitempty"`
	MinCancellationNotice    int    `json:"min_cancellation_notice,omitempty"`
	ReschedulingURL          string `json:"rescheduling_url,omitempty"`
	CancellationURL          string `json:"cancellation_url,omitempty"`
	OrganizerConfirmationURL string `json:"organizer_confirmation_url,omitempty"`
	ConfirmationRedirectURL  string `json:"confirmation_redirect_url,omitempty"`
	HideCancellationOptions  bool   `json:"hide_cancellation_options,omitempty"`
	HideReschedulingOptions  bool   `json:"hide_rescheduling_options,omitempty"`
	HideAdditionalGuests     bool   `json:"hide_additional_guests,omitempty"`
	CancellationPolicy       string `json:"cancellation_policy,omitempty"`
}

// AppearanceSettings defines scheduler appearance.
type AppearanceSettings struct {
	SubmitButtonLabel string `json:"submit_button_label,omitempty"`
	ThankYouMessage   string `json:"thank_you_message,omitempty"`
	Color             string `json:"color,omitempty"`
	Logo              string `json:"logo,omitempty"`
}

// ConfigurationRequest represents a request to create/update a configuration.
type ConfigurationRequest struct {
	Participants        []Participant       `json:"participants,omitempty"`
	Availability        *Availability       `json:"availability,omitempty"`
	EventBooking        *EventBooking       `json:"event_booking,omitempty"`
	Scheduler           *SchedulerSettings  `json:"scheduler,omitempty"`
	AppearanceSettings  *AppearanceSettings `json:"appearance,omitempty"`
	RequiresSessionAuth bool                `json:"requires_session_auth,omitempty"`
}

// Session represents a scheduler session.
type Session struct {
	SessionID string `json:"session_id,omitempty"`
}

// SessionRequest represents a request to create a session.
type SessionRequest struct {
	ConfigurationID string `json:"configuration_id"`
	TimeToLive      int    `json:"time_to_live,omitempty"`
}

// Booking represents a scheduled booking.
type Booking struct {
	BookingID        string               `json:"booking_id,omitempty"`
	EventID          string               `json:"event_id,omitempty"`
	Title            string               `json:"title,omitempty"`
	Description      string               `json:"description,omitempty"`
	Organizer        BookingParticipant   `json:"organizer,omitempty"`
	Status           string               `json:"status,omitempty"`
	StartTime        int64                `json:"start_time,omitempty"`
	EndTime          int64                `json:"end_time,omitempty"`
	Participants     []BookingParticipant `json:"participants,omitempty"`
	AdditionalGuests []string             `json:"additional_guests,omitempty"`
	CreatedAt        int64                `json:"created_at,omitempty"`
	UpdatedAt        int64                `json:"updated_at,omitempty"`
}

// BookingParticipant represents a booking participant.
type BookingParticipant struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// BookingRequest represents a request to create a booking.
type BookingRequest struct {
	StartTime        int64              `json:"start_time"`
	EndTime          int64              `json:"end_time"`
	Guest            BookingParticipant `json:"guest"`
	AdditionalGuests []string           `json:"additional_guests,omitempty"`
	AdditionalFields map[string]string  `json:"additional_fields,omitempty"`
}

// ConfirmBookingRequest represents a request to confirm a booking.
type ConfirmBookingRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// RescheduleBookingRequest represents a request to reschedule a booking.
type RescheduleBookingRequest struct {
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Reason    string `json:"reason,omitempty"`
}

// ListConfigurationsOptions specifies options for listing configurations.
type ListConfigurationsOptions struct {
	Limit     *int   `json:"limit,omitempty"`
	PageToken string `json:"page_token,omitempty"`
}

// Values converts ListConfigurationsOptions to URL query parameters.
func (o *ListConfigurationsOptions) Values() map[string]any {
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

// ListBookingsOptions specifies options for listing bookings.
type ListBookingsOptions struct {
	ConfigurationID string `json:"configuration_id,omitempty"`
	Limit           *int   `json:"limit,omitempty"`
	PageToken       string `json:"page_token,omitempty"`
}

// Values converts ListBookingsOptions to URL query parameters.
func (o *ListBookingsOptions) Values() map[string]any {
	if o == nil {
		return nil
	}
	v := make(map[string]any)
	if o.ConfigurationID != "" {
		v["configuration_id"] = o.ConfigurationID
	}
	if o.Limit != nil {
		v["limit"] = *o.Limit
	}
	if o.PageToken != "" {
		v["page_token"] = o.PageToken
	}
	return v
}

// StartDateTime returns the booking start time as time.Time.
func (b *Booking) StartDateTime() time.Time {
	return time.Unix(b.StartTime, 0)
}

// EndDateTime returns the booking end time as time.Time.
func (b *Booking) EndDateTime() time.Time {
	return time.Unix(b.EndTime, 0)
}
