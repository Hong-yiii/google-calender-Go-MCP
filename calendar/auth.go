package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// AuthManager handles authentication for Google Calendar API
type AuthManager struct {
	config *CalendarConfig
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(config *CalendarConfig) *AuthManager {
	return &AuthManager{
		config: config,
	}
}

// GetCalendarService creates and returns an authenticated Google Calendar service
func (a *AuthManager) GetCalendarService(ctx context.Context) (*calendar.Service, error) {
	client, err := a.getAuthenticatedClient(ctx)
	if err != nil {
		return nil, NewAuthenticationError(ErrCodeInvalidCredentials, "Failed to create authenticated client", err)
	}

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to create calendar service", err)
	}

	return service, nil
}

// getAuthenticatedClient returns an authenticated HTTP client
func (a *AuthManager) getAuthenticatedClient(ctx context.Context) (*http.Client, error) {
	// Try to load credentials
	creds, err := a.loadCredentials()
	if err != nil {
		return nil, err
	}

	// Check if it's a service account or OAuth2 credentials
	if a.isServiceAccount(creds) {
		return a.getServiceAccountClient(ctx, creds)
	}

	return a.getOAuth2Client(ctx, creds)
}

// loadCredentials loads credentials from file or environment variable
func (a *AuthManager) loadCredentials() ([]byte, error) {
	credentialsSource := a.config.CredentialsJSON

	if credentialsSource == "" {
		return nil, NewAuthenticationError(ErrCodeMissingCredentials, "No credentials provided", nil)
	}

	// Check if it's a JSON string or file path
	if strings.HasPrefix(credentialsSource, "{") {
		// It's a JSON string
		return []byte(credentialsSource), nil
	}

	// It's a file path
	if _, err := os.Stat(credentialsSource); os.IsNotExist(err) {
		return nil, NewAuthenticationError(ErrCodeMissingCredentials, fmt.Sprintf("Credentials file not found: %s", credentialsSource), err)
	}

	creds, err := ioutil.ReadFile(credentialsSource)
	if err != nil {
		return nil, NewAuthenticationError(ErrCodeInvalidCredentials, "Failed to read credentials file", err)
	}

	return creds, nil
}

// isServiceAccount checks if the credentials are for a service account
func (a *AuthManager) isServiceAccount(creds []byte) bool {
	var credData map[string]interface{}
	if err := json.Unmarshal(creds, &credData); err != nil {
		return false
	}

	credType, exists := credData["type"]
	return exists && credType == "service_account"
}

// getServiceAccountClient creates an authenticated client using service account credentials
func (a *AuthManager) getServiceAccountClient(ctx context.Context, creds []byte) (*http.Client, error) {
	config, err := google.JWTConfigFromJSON(creds, calendar.CalendarScope)
	if err != nil {
		return nil, NewAuthenticationError(ErrCodeInvalidCredentials, "Failed to parse service account credentials", err)
	}

	client := config.Client(ctx)
	return client, nil
}

// getOAuth2Client creates an authenticated client using OAuth2 credentials
func (a *AuthManager) getOAuth2Client(ctx context.Context, creds []byte) (*http.Client, error) {
	_, err := google.ConfigFromJSON(creds, calendar.CalendarScope)
	if err != nil {
		return nil, NewAuthenticationError(ErrCodeInvalidCredentials, "Failed to parse OAuth2 credentials", err)
	}

	// For OAuth2, we need to handle the token exchange
	// In a real application, you would implement the full OAuth2 flow
	// For now, we'll return an error indicating that OAuth2 flow needs to be implemented
	return nil, NewAuthenticationError(ErrCodeInvalidCredentials, "OAuth2 flow not implemented. Please use service account credentials", nil)
}

// ValidateCredentials validates that the credentials are valid and can access the calendar
func (a *AuthManager) ValidateCredentials(ctx context.Context) error {
	service, err := a.GetCalendarService(ctx)
	if err != nil {
		return err
	}

	// Try to get calendar info to validate access
	_, err = service.Calendars.Get(a.config.CalendarID).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "notFound") {
			return NewNotFoundError(ErrCodeCalendarNotFound, fmt.Sprintf("Calendar not found: %s", a.config.CalendarID))
		}
		if strings.Contains(err.Error(), "forbidden") {
			return NewPermissionError(ErrCodePermissionDenied, "Access denied to calendar")
		}
		return NewInternalError(ErrCodeServiceUnavailable, "Failed to validate calendar access", err)
	}

	log.Printf("Successfully validated credentials for calendar: %s", a.config.CalendarID)
	return nil
}

// GetCalendarInfo retrieves basic information about the configured calendar
func (a *AuthManager) GetCalendarInfo(ctx context.Context) (*CalendarInfo, error) {
	service, err := a.GetCalendarService(ctx)
	if err != nil {
		return nil, err
	}

	calendarInfo, err := service.Calendars.Get(a.config.CalendarID).Context(ctx).Do()
	if err != nil {
		if strings.Contains(err.Error(), "notFound") {
			return nil, NewNotFoundError(ErrCodeCalendarNotFound, fmt.Sprintf("Calendar not found: %s", a.config.CalendarID))
		}
		return nil, NewInternalError(ErrCodeServiceUnavailable, "Failed to get calendar info", err)
	}

	return &CalendarInfo{
		ID:          calendarInfo.Id,
		Summary:     calendarInfo.Summary,
		Description: calendarInfo.Description,
		TimeZone:    calendarInfo.TimeZone,
		Location:    calendarInfo.Location,
	}, nil
}

// RefreshCredentials refreshes the authentication credentials if needed
func (a *AuthManager) RefreshCredentials(ctx context.Context) error {
	// For service accounts, credentials don't need refreshing
	// For OAuth2, this would handle token refresh
	return a.ValidateCredentials(ctx)
}
