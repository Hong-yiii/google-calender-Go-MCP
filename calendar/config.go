package calendar

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadConfig loads configuration from environment variables
func LoadConfig() (*CalendarConfig, error) {
	config := &CalendarConfig{
		CredentialsJSON: getEnvWithDefault("GOOGLE_CALENDAR_CREDENTIALS_JSON", ""),
		CalendarID:      getEnvWithDefault("GOOGLE_CALENDAR_ID", DefaultCalendarID),
		TimeZone:        getEnvWithDefault("GOOGLE_CALENDAR_TIMEZONE", DefaultTimeZone),
		ServerName:      getEnvWithDefault("MCP_SERVER_NAME", "Google Calendar MCP Server"),
		ServerVersion:   getEnvWithDefault("MCP_SERVER_VERSION", "1.0.0"),
		LogLevel:        getEnvWithDefault("LOG_LEVEL", "info"),
		Environment:     getEnvWithDefault("ENVIRONMENT", "development"),
		Debug:           getEnvBool("DEBUG", false),
	}

	if err := validateConfig(config); err != nil {
		return nil, NewConfigurationError(ErrCodeConfigurationError, "Invalid configuration", err)
	}

	return config, nil
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets a boolean environment variable with a default value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// validateConfig validates the configuration
func validateConfig(config *CalendarConfig) error {
	var errors []string

	// Validate required fields
	if config.CredentialsJSON == "" {
		errors = append(errors, "GOOGLE_CALENDAR_CREDENTIALS_JSON is required")
	}

	// Validate credentials file exists if it's a file path
	if config.CredentialsJSON != "" && !strings.HasPrefix(config.CredentialsJSON, "{") {
		if _, err := os.Stat(config.CredentialsJSON); os.IsNotExist(err) {
			errors = append(errors, fmt.Sprintf("Credentials file not found: %s", config.CredentialsJSON))
		}
	}

	// Validate calendar ID format
	if config.CalendarID == "" {
		errors = append(errors, "Calendar ID cannot be empty")
	}

	// Validate timezone
	if config.TimeZone == "" {
		errors = append(errors, "Timezone cannot be empty")
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[strings.ToLower(config.LogLevel)] {
		errors = append(errors, fmt.Sprintf("Invalid log level: %s. Valid levels are: debug, info, warn, error, fatal", config.LogLevel))
	}

	// Validate environment
	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	if !validEnvironments[strings.ToLower(config.Environment)] {
		errors = append(errors, fmt.Sprintf("Invalid environment: %s. Valid environments are: development, staging, production, test", config.Environment))
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(code, message string, cause error) CalendarError {
	return &calendarError{
		code:    code,
		errType: ErrorTypeInvalidInput,
		message: message,
		cause:   cause,
	}
}

// IsProduction returns true if the environment is production
func (c *CalendarConfig) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

// IsDevelopment returns true if the environment is development
func (c *CalendarConfig) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

// GetLogLevel returns the log level in lowercase
func (c *CalendarConfig) GetLogLevel() string {
	return strings.ToLower(c.LogLevel)
}

// String returns a string representation of the config (without sensitive data)
func (c *CalendarConfig) String() string {
	return fmt.Sprintf("CalendarConfig{CalendarID: %s, TimeZone: %s, Environment: %s, Debug: %t}",
		c.CalendarID, c.TimeZone, c.Environment, c.Debug)
}
