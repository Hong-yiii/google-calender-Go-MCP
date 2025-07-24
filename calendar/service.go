package calendar

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

// CalendarService defines the interface for calendar operations
type CalendarService interface {
	// Core operations
	CheckAvailability(ctx context.Context, startTime, endTime time.Time) ([]TimeSlot, error)
	CreateEvent(ctx context.Context, event *EventCreateRequest) (*Event, error)
	ListEvents(ctx context.Context, timeMin, timeMax time.Time) ([]*Event, error)
	UpdateEvent(ctx context.Context, eventID string, update *EventUpdateRequest) (*Event, error)
	DeleteEvent(ctx context.Context, eventID string) error

	// Utility operations
	GetCalendarInfo(ctx context.Context) (*CalendarInfo, error)
	SearchEvents(ctx context.Context, query string, timeMin, timeMax time.Time) ([]*Event, error)
}

// googleCalendarService implements CalendarService using Google Calendar API
type googleCalendarService struct {
	authManager *AuthManager
	config      *CalendarConfig
}

// NewCalendarService creates a new calendar service
func NewCalendarService(config *CalendarConfig) (CalendarService, error) {
	authManager := NewAuthManager(config)

	// Validate credentials during initialization
	ctx := context.Background()
	if err := authManager.ValidateCredentials(ctx); err != nil {
		return nil, fmt.Errorf("failed to validate credentials: %w", err)
	}

	return &googleCalendarService{
		authManager: authManager,
		config:      config,
	}, nil
}

// CheckAvailability checks for available time slots in the given time range
func (s *googleCalendarService) CheckAvailability(ctx context.Context, startTime, endTime time.Time) ([]TimeSlot, error) {
	if startTime.After(endTime) {
		return nil, NewInvalidInputError(ErrCodeInvalidTimeRange, "Start time must be before end time", "")
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	// Get events in the time range
	events, err := service.Events.List(s.config.CalendarID).
		TimeMin(startTime.Format(time.RFC3339)).
		TimeMax(endTime.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Context(ctx).
		Do()

	if err != nil {
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to retrieve events", err)
	}

	// Convert events to time slots and find free slots
	return s.calculateFreeTimeSlots(startTime, endTime, events.Items), nil
}

// CreateEvent creates a new calendar event
func (s *googleCalendarService) CreateEvent(ctx context.Context, eventReq *EventCreateRequest) (*Event, error) {
	if err := s.validateEventCreateRequest(eventReq); err != nil {
		return nil, err
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	// Convert our event to Google Calendar event
	googleEvent := &calendar.Event{
		Summary:     eventReq.Summary,
		Description: eventReq.Description,
		Location:    eventReq.Location,
		Start: &calendar.EventDateTime{
			DateTime: eventReq.StartTime.Format(time.RFC3339),
			TimeZone: s.config.TimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: eventReq.EndTime.Format(time.RFC3339),
			TimeZone: s.config.TimeZone,
		},
	}

	// Add attendees if provided
	if len(eventReq.Attendees) > 0 {
		googleEvent.Attendees = make([]*calendar.EventAttendee, len(eventReq.Attendees))
		for i, email := range eventReq.Attendees {
			googleEvent.Attendees[i] = &calendar.EventAttendee{
				Email: email,
			}
		}
	}

	createdEvent, err := service.Events.Insert(s.config.CalendarID, googleEvent).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return nil, NewPermissionError(ErrCodePermissionDenied, "Permission denied to create event")
		}
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to create event", err)
	}

	return s.convertGoogleEventToEvent(createdEvent), nil
}

// ListEvents retrieves events in the specified time range
func (s *googleCalendarService) ListEvents(ctx context.Context, timeMin, timeMax time.Time) ([]*Event, error) {
	if timeMin.After(timeMax) {
		return nil, NewInvalidInputError(ErrCodeInvalidTimeRange, "Start time must be before end time", "")
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	googleEvents, err := service.Events.List(s.config.CalendarID).
		TimeMin(timeMin.Format(time.RFC3339)).
		TimeMax(timeMax.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		MaxResults(int64(DefaultMaxResults)).
		Context(ctx).
		Do()

	if err != nil {
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to retrieve events", err)
	}

	events := make([]*Event, len(googleEvents.Items))
	for i, googleEvent := range googleEvents.Items {
		events[i] = s.convertGoogleEventToEvent(googleEvent)
	}

	return events, nil
}

// UpdateEvent updates an existing calendar event
func (s *googleCalendarService) UpdateEvent(ctx context.Context, eventID string, update *EventUpdateRequest) (*Event, error) {
	if eventID == "" {
		return nil, NewInvalidInputError(ErrCodeInvalidEventData, "Event ID is required", "")
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	// Get the existing event
	existingEvent, err := service.Events.Get(s.config.CalendarID, eventID).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "notFound") {
			return nil, NewNotFoundError(ErrCodeEventNotFound, fmt.Sprintf("Event not found: %s", eventID))
		}
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to retrieve event", err)
	}

	// Apply updates
	s.applyEventUpdates(existingEvent, update)

	// Update the event
	updatedEvent, err := service.Events.Update(s.config.CalendarID, eventID, existingEvent).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return nil, NewPermissionError(ErrCodePermissionDenied, "Permission denied to update event")
		}
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to update event", err)
	}

	return s.convertGoogleEventToEvent(updatedEvent), nil
}

// DeleteEvent deletes a calendar event
func (s *googleCalendarService) DeleteEvent(ctx context.Context, eventID string) error {
	if eventID == "" {
		return NewInvalidInputError(ErrCodeInvalidEventData, "Event ID is required", "")
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return err
	}

	err = service.Events.Delete(s.config.CalendarID, eventID).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "notFound") {
			return NewNotFoundError(ErrCodeEventNotFound, fmt.Sprintf("Event not found: %s", eventID))
		}
		if strings.Contains(err.Error(), "forbidden") {
			return NewPermissionError(ErrCodePermissionDenied, "Permission denied to delete event")
		}
		return NewInternalError(ErrCodeServiceUnavailable, "Failed to delete event", err)
	}

	log.Printf("Successfully deleted event: %s", eventID)
	return nil
}

// GetCalendarInfo retrieves basic calendar information
func (s *googleCalendarService) GetCalendarInfo(ctx context.Context) (*CalendarInfo, error) {
	return s.authManager.GetCalendarInfo(ctx)
}

// SearchEvents searches for events matching the given query
func (s *googleCalendarService) SearchEvents(ctx context.Context, query string, timeMin, timeMax time.Time) ([]*Event, error) {
	if query == "" {
		return nil, NewInvalidInputError(ErrCodeInvalidEventData, "Search query is required", "")
	}

	service, err := s.authManager.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	call := service.Events.List(s.config.CalendarID).
		Q(query).
		SingleEvents(true).
		OrderBy("startTime").
		MaxResults(int64(DefaultMaxResults)).
		Context(ctx)

	if !timeMin.IsZero() {
		call = call.TimeMin(timeMin.Format(time.RFC3339))
	}
	if !timeMax.IsZero() {
		call = call.TimeMax(timeMax.Format(time.RFC3339))
	}

	googleEvents, err := call.Do()
	if err != nil {
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to search events", err)
	}

	events := make([]*Event, len(googleEvents.Items))
	for i, googleEvent := range googleEvents.Items {
		events[i] = s.convertGoogleEventToEvent(googleEvent)
	}

	return events, nil
}

// Helper methods

// calculateFreeTimeSlots calculates free time slots between events
func (s *googleCalendarService) calculateFreeTimeSlots(startTime, endTime time.Time, events []*calendar.Event) []TimeSlot {
	var timeSlots []TimeSlot

	// If no events, the entire range is free
	if len(events) == 0 {
		return []TimeSlot{{
			Start: startTime,
			End:   endTime,
			Free:  true,
		}}
	}

	// Sort events by start time
	sort.Slice(events, func(i, j int) bool {
		iStart, _ := time.Parse(time.RFC3339, events[i].Start.DateTime)
		jStart, _ := time.Parse(time.RFC3339, events[j].Start.DateTime)
		return iStart.Before(jStart)
	})

	currentTime := startTime

	for _, event := range events {
		eventStart, _ := time.Parse(time.RFC3339, event.Start.DateTime)
		eventEnd, _ := time.Parse(time.RFC3339, event.End.DateTime)

		// Skip events outside our time range
		if eventEnd.Before(startTime) || eventStart.After(endTime) {
			continue
		}

		// Adjust event times to our range
		if eventStart.Before(startTime) {
			eventStart = startTime
		}
		if eventEnd.After(endTime) {
			eventEnd = endTime
		}

		// Add free slot before this event
		if currentTime.Before(eventStart) {
			timeSlots = append(timeSlots, TimeSlot{
				Start: currentTime,
				End:   eventStart,
				Free:  true,
			})
		}

		// Add busy slot for this event
		timeSlots = append(timeSlots, TimeSlot{
			Start: eventStart,
			End:   eventEnd,
			Free:  false,
		})

		// Update current time
		if eventEnd.After(currentTime) {
			currentTime = eventEnd
		}
	}

	// Add final free slot if needed
	if currentTime.Before(endTime) {
		timeSlots = append(timeSlots, TimeSlot{
			Start: currentTime,
			End:   endTime,
			Free:  true,
		})
	}

	return timeSlots
}

// convertGoogleEventToEvent converts a Google Calendar event to our Event struct
func (s *googleCalendarService) convertGoogleEventToEvent(googleEvent *calendar.Event) *Event {
	startTime, _ := time.Parse(time.RFC3339, googleEvent.Start.DateTime)
	endTime, _ := time.Parse(time.RFC3339, googleEvent.End.DateTime)
	createdTime, _ := time.Parse(time.RFC3339, googleEvent.Created)
	updatedTime, _ := time.Parse(time.RFC3339, googleEvent.Updated)

	var attendees []string
	for _, attendee := range googleEvent.Attendees {
		attendees = append(attendees, attendee.Email)
	}

	return &Event{
		ID:          googleEvent.Id,
		Summary:     googleEvent.Summary,
		Description: googleEvent.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Location:    googleEvent.Location,
		Attendees:   attendees,
		Status:      googleEvent.Status,
		CreatedAt:   createdTime,
		UpdatedAt:   updatedTime,
	}
}

// validateEventCreateRequest validates an event creation request
func (s *googleCalendarService) validateEventCreateRequest(req *EventCreateRequest) error {
	if req.Summary == "" {
		return NewInvalidInputError(ErrCodeInvalidEventData, "Event summary is required", "")
	}

	if req.StartTime.IsZero() {
		return NewInvalidInputError(ErrCodeInvalidTimeFormat, "Start time is required", "")
	}

	if req.EndTime.IsZero() {
		return NewInvalidInputError(ErrCodeInvalidTimeFormat, "End time is required", "")
	}

	if req.StartTime.After(req.EndTime) {
		return NewInvalidInputError(ErrCodeInvalidTimeRange, "Start time must be before end time", "")
	}

	// Validate attendee email formats (basic validation)
	for _, email := range req.Attendees {
		if !strings.Contains(email, "@") {
			return NewInvalidInputError(ErrCodeInvalidEventData, fmt.Sprintf("Invalid email format: %s", email), "")
		}
	}

	return nil
}

// applyEventUpdates applies updates to an existing Google Calendar event
func (s *googleCalendarService) applyEventUpdates(event *calendar.Event, update *EventUpdateRequest) {
	if update.Summary != nil {
		event.Summary = *update.Summary
	}

	if update.Description != nil {
		event.Description = *update.Description
	}

	if update.Location != nil {
		event.Location = *update.Location
	}

	if update.StartTime != nil {
		event.Start = &calendar.EventDateTime{
			DateTime: update.StartTime.Format(time.RFC3339),
			TimeZone: s.config.TimeZone,
		}
	}

	if update.EndTime != nil {
		event.End = &calendar.EventDateTime{
			DateTime: update.EndTime.Format(time.RFC3339),
			TimeZone: s.config.TimeZone,
		}
	}

	if update.Attendees != nil {
		event.Attendees = make([]*calendar.EventAttendee, len(update.Attendees))
		for i, email := range update.Attendees {
			event.Attendees[i] = &calendar.EventAttendee{
				Email: email,
			}
		}
	}
}
