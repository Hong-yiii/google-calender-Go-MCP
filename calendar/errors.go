package calendar

import (
	"fmt"
)

// ErrorType represents different types of calendar errors
type ErrorType string

const (
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypePermission     ErrorType = "PERMISSION_DENIED"
	ErrorTypeInvalidInput   ErrorType = "INVALID_INPUT"
	ErrorTypeQuotaExceeded  ErrorType = "QUOTA_EXCEEDED"
	ErrorTypeInternal       ErrorType = "INTERNAL_ERROR"
	ErrorTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
	ErrorTypeNetwork        ErrorType = "NETWORK_ERROR"
	ErrorTypeTimeout        ErrorType = "TIMEOUT_ERROR"
	ErrorTypeConflict       ErrorType = "CONFLICT_ERROR"
)

// CalendarError represents a calendar-specific error
type CalendarError interface {
	error
	Code() string
	Type() ErrorType
	Details() string
}

// calendarError implements CalendarError interface
type calendarError struct {
	code    string
	errType ErrorType
	message string
	details string
	cause   error
}

func (e *calendarError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *calendarError) Code() string {
	return e.code
}

func (e *calendarError) Type() ErrorType {
	return e.errType
}

func (e *calendarError) Details() string {
	return e.details
}

func (e *calendarError) Unwrap() error {
	return e.cause
}

// Error constructors

// NewNotFoundError creates a new not found error
func NewNotFoundError(code, message string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeNotFound,
		message: message,
	}
}

// NewPermissionError creates a new permission denied error
func NewPermissionError(code, message string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypePermission,
		message: message,
	}
}

// NewInvalidInputError creates a new invalid input error
func NewInvalidInputError(code, message, details string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeInvalidInput,
		message: message,
		details: details,
	}
}

// NewQuotaExceededError creates a new quota exceeded error
func NewQuotaExceededError(code, message string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeQuotaExceeded,
		message: message,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(code, message string, cause error) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeInternal,
		message: message,
		cause:   cause,
	}
}

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(code, message string, cause error) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeAuthentication,
		message: message,
		cause:   cause,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(code, message string, cause error) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeNetwork,
		message: message,
		cause:   cause,
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(code, message string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeTimeout,
		message: message,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(code, message string) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeConflict,
		message: message,
	}
}

// Common error codes
const (
	ErrCodeEventNotFound      = "EVENT_NOT_FOUND"
	ErrCodeCalendarNotFound   = "CALENDAR_NOT_FOUND"
	ErrCodeInvalidTimeFormat  = "INVALID_TIME_FORMAT"
	ErrCodeInvalidTimeRange   = "INVALID_TIME_RANGE"
	ErrCodeMissingCredentials = "MISSING_CREDENTIALS"
	ErrCodeInvalidCredentials = "INVALID_CREDENTIALS"
	ErrCodeAPIQuotaExceeded   = "API_QUOTA_EXCEEDED"
	ErrCodePermissionDenied   = "PERMISSION_DENIED"
	ErrCodeNetworkTimeout     = "NETWORK_TIMEOUT"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrCodeInvalidEventData   = "INVALID_EVENT_DATA"
	ErrCodeEventConflict      = "EVENT_CONFLICT"
	ErrCodeConfigurationError = "CONFIGURATION_ERROR"
)

// ErrorResponse represents an error response for MCP tools
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Details string `json:"details,omitempty"`
}

// ToErrorResponse converts a CalendarError to an ErrorResponse
func ToErrorResponse(err CalendarError) *ErrorResponse {
	return &ErrorResponse{
		Code:    err.Code(),
		Message: err.Error(),
		Type:    string(err.Type()),
		Details: err.Details(),
	}
}

// IsCalendarError checks if an error is a CalendarError
func IsCalendarError(err error) bool {
	_, ok := err.(CalendarError)
	return ok
}

// GetErrorType returns the error type if it's a CalendarError, otherwise returns ErrorTypeInternal
func GetErrorType(err error) ErrorType {
	if calErr, ok := err.(CalendarError); ok {
		return calErr.Type()
	}
	return ErrorTypeInternal
}
