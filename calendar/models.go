package calendar

import (
	"time"
)

// Event represents a calendar event
type Event struct {
	ID          string    `json:"id"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location,omitempty"`
	Attendees   []string  `json:"attendees,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TimeSlot represents a time slot with availability information
type TimeSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Free  bool      `json:"free"`
}

// CalendarConfig holds configuration for calendar operations
type CalendarConfig struct {
	CredentialsJSON string `json:"credentials_json"`
	CalendarID      string `json:"calendar_id"`
	TimeZone        string `json:"timezone"`
	ServerName      string `json:"server_name"`
	ServerVersion   string `json:"server_version"`
	LogLevel        string `json:"log_level"`
	Environment     string `json:"environment"`
	Debug           bool   `json:"debug"`
}

// CalendarInfo represents basic calendar information
type CalendarInfo struct {
	ID          string `json:"id"`
	Summary     string `json:"summary"`
	Description string `json:"description,omitempty"`
	TimeZone    string `json:"timezone"`
	Location    string `json:"location,omitempty"`
}

// EventCreateRequest represents a request to create an event
type EventCreateRequest struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Location    string    `json:"location,omitempty"`
	Attendees   []string  `json:"attendees,omitempty"`
}

// EventUpdateRequest represents a request to update an event
type EventUpdateRequest struct {
	Summary     *string    `json:"summary,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Attendees   []string   `json:"attendees,omitempty"`
}

// AvailabilityRequest represents a request to check availability
type AvailabilityRequest struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	CalendarID string    `json:"calendar_id,omitempty"`
}

// SearchRequest represents a request to search events
type SearchRequest struct {
	Query      string     `json:"query"`
	StartTime  *time.Time `json:"start_time,omitempty"`
	EndTime    *time.Time `json:"end_time,omitempty"`
	MaxResults int        `json:"max_results,omitempty"`
}

// ListEventsRequest represents a request to list events
type ListEventsRequest struct {
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	MaxResults int       `json:"max_results,omitempty"`
}

// EventStatus constants
const (
	EventStatusConfirmed = "confirmed"
	EventStatusTentative = "tentative"
	EventStatusCancelled = "cancelled"
)

// Default values
const (
	DefaultMaxResults = 50
	DefaultTimeZone   = "UTC"
	DefaultCalendarID = "primary"
)
