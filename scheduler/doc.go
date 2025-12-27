// Package scheduler provides types for the Nylas Scheduler API.
//
// The Scheduler API allows you to create scheduling configurations,
// manage availability, and handle bookings for appointments.
//
// Example usage:
//
//	// Create a configuration
//	config, err := client.Scheduler.CreateConfiguration(ctx, grantID, &scheduler.ConfigurationRequest{
//		Participants: []scheduler.Participant{
//			{Email: "organizer@example.com", IsOrganizer: true},
//		},
//		Availability: &scheduler.Availability{
//			DurationMinutes: 30,
//		},
//		EventBooking: &scheduler.EventBooking{
//			Title: "Meeting",
//		},
//	})
//
//	// List configurations
//	configs, err := client.Scheduler.ListConfigurations(ctx, grantID, nil)
//
//	// Create a booking
//	booking, err := client.Scheduler.CreateBooking(ctx, configID, &scheduler.BookingRequest{
//		StartTime: startUnix,
//		EndTime:   endUnix,
//		Guest:     scheduler.BookingParticipant{Email: "guest@example.com"},
//	})
package scheduler
