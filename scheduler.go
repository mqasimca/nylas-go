package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mqasimca/nylas-go/scheduler"
)

// ListConfigurations returns all scheduler configurations for a grant.
func (s *SchedulerService) ListConfigurations(ctx context.Context, grantID string, opts *scheduler.ListConfigurationsOptions) (*ListResponse[scheduler.Configuration], error) {
	path := fmt.Sprintf("/v3/grants/%s/scheduling/configurations", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ListConfigurations: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []scheduler.Configuration
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ListConfigurations: %w", err)
	}

	return &ListResponse[scheduler.Configuration]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// GetConfiguration retrieves a single scheduler configuration by ID.
func (s *SchedulerService) GetConfiguration(ctx context.Context, grantID, configID string) (*scheduler.Configuration, error) {
	path := fmt.Sprintf("/v3/grants/%s/scheduling/configurations/%s", grantID, configID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("scheduler.GetConfiguration(%s): %w", configID, err)
	}

	var config scheduler.Configuration
	_, err = s.client.Do(req, &config)
	if err != nil {
		return nil, fmt.Errorf("scheduler.GetConfiguration(%s): %w", configID, err)
	}

	return &config, nil
}

// CreateConfiguration creates a new scheduler configuration.
func (s *SchedulerService) CreateConfiguration(ctx context.Context, grantID string, configReq *scheduler.ConfigurationRequest) (*scheduler.Configuration, error) {
	path := fmt.Sprintf("/v3/grants/%s/scheduling/configurations", grantID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, configReq)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateConfiguration: %w", err)
	}

	var config scheduler.Configuration
	_, err = s.client.Do(req, &config)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateConfiguration: %w", err)
	}

	return &config, nil
}

// UpdateConfiguration updates an existing scheduler configuration.
func (s *SchedulerService) UpdateConfiguration(ctx context.Context, grantID, configID string, configReq *scheduler.ConfigurationRequest) (*scheduler.Configuration, error) {
	path := fmt.Sprintf("/v3/grants/%s/scheduling/configurations/%s", grantID, configID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, configReq)
	if err != nil {
		return nil, fmt.Errorf("scheduler.UpdateConfiguration(%s): %w", configID, err)
	}

	var config scheduler.Configuration
	_, err = s.client.Do(req, &config)
	if err != nil {
		return nil, fmt.Errorf("scheduler.UpdateConfiguration(%s): %w", configID, err)
	}

	return &config, nil
}

// DeleteConfiguration deletes a scheduler configuration.
func (s *SchedulerService) DeleteConfiguration(ctx context.Context, grantID, configID string) error {
	path := fmt.Sprintf("/v3/grants/%s/scheduling/configurations/%s", grantID, configID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("scheduler.DeleteConfiguration(%s): %w", configID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("scheduler.DeleteConfiguration(%s): %w", configID, err)
	}

	return nil
}

// CreateSession creates a new scheduler session for a configuration.
func (s *SchedulerService) CreateSession(ctx context.Context, sessionReq *scheduler.SessionRequest) (*scheduler.Session, error) {
	path := "/v3/scheduling/sessions"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, sessionReq)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateSession: %w", err)
	}

	var session scheduler.Session
	_, err = s.client.Do(req, &session)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateSession: %w", err)
	}

	return &session, nil
}

// ListBookings returns all bookings for a configuration.
func (s *SchedulerService) ListBookings(ctx context.Context, configID string, opts *scheduler.ListBookingsOptions) (*ListResponse[scheduler.Booking], error) {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings", configID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ListBookings: %w", err)
	}

	if opts != nil {
		q := req.URL.Query()
		setQueryParams(q, opts.Values())
		req.URL.RawQuery = q.Encode()
	}

	var data []scheduler.Booking
	nextCursor, requestID, err := s.client.DoList(req, &data)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ListBookings: %w", err)
	}

	return &ListResponse[scheduler.Booking]{
		Data:       data,
		NextCursor: nextCursor,
		RequestID:  requestID,
	}, nil
}

// GetBooking retrieves a single booking by ID.
func (s *SchedulerService) GetBooking(ctx context.Context, configID, bookingID string) (*scheduler.Booking, error) {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings/%s", configID, bookingID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("scheduler.GetBooking(%s): %w", bookingID, err)
	}

	var booking scheduler.Booking
	_, err = s.client.Do(req, &booking)
	if err != nil {
		return nil, fmt.Errorf("scheduler.GetBooking(%s): %w", bookingID, err)
	}

	return &booking, nil
}

// CreateBooking creates a new booking for a configuration.
func (s *SchedulerService) CreateBooking(ctx context.Context, configID string, bookingReq *scheduler.BookingRequest) (*scheduler.Booking, error) {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings", configID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, bookingReq)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateBooking: %w", err)
	}

	var booking scheduler.Booking
	_, err = s.client.Do(req, &booking)
	if err != nil {
		return nil, fmt.Errorf("scheduler.CreateBooking: %w", err)
	}

	return &booking, nil
}

// ConfirmBooking confirms or rejects a pending booking.
func (s *SchedulerService) ConfirmBooking(ctx context.Context, configID, bookingID string, confirm *scheduler.ConfirmBookingRequest) (*scheduler.Booking, error) {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings/%s", configID, bookingID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, confirm)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ConfirmBooking(%s): %w", bookingID, err)
	}

	var booking scheduler.Booking
	_, err = s.client.Do(req, &booking)
	if err != nil {
		return nil, fmt.Errorf("scheduler.ConfirmBooking(%s): %w", bookingID, err)
	}

	return &booking, nil
}

// RescheduleBooking reschedules an existing booking.
func (s *SchedulerService) RescheduleBooking(ctx context.Context, configID, bookingID string, reschedule *scheduler.RescheduleBookingRequest) (*scheduler.Booking, error) {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings/%s/reschedule", configID, bookingID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, reschedule)
	if err != nil {
		return nil, fmt.Errorf("scheduler.RescheduleBooking(%s): %w", bookingID, err)
	}

	var booking scheduler.Booking
	_, err = s.client.Do(req, &booking)
	if err != nil {
		return nil, fmt.Errorf("scheduler.RescheduleBooking(%s): %w", bookingID, err)
	}

	return &booking, nil
}

// CancelBooking cancels an existing booking.
func (s *SchedulerService) CancelBooking(ctx context.Context, configID, bookingID, reason string) error {
	path := fmt.Sprintf("/v3/scheduling/configurations/%s/bookings/%s/cancel", configID, bookingID)

	body := map[string]string{}
	if reason != "" {
		body["reason"] = reason
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return fmt.Errorf("scheduler.CancelBooking(%s): %w", bookingID, err)
	}

	_, err = s.client.Do(req, nil)
	if err != nil {
		return fmt.Errorf("scheduler.CancelBooking(%s): %w", bookingID, err)
	}

	return nil
}
