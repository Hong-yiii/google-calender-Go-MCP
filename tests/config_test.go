package tests

import (
	"os"
	"path/filepath"
	"testing"

	"google_cal_mcp_golang/calendar"

	"github.com/joho/godotenv"
)

func TestConfigValidation(t *testing.T) {
	// Test invalid configuration (missing credentials)
	os.Clearenv()

	_, err := calendar.LoadConfig()
	if err == nil {
		t.Error("Expected error for missing credentials, got nil")
	}

	// Test valid minimal configuration
	os.Setenv("GOOGLE_CALENDAR_CREDENTIALS_JSON", `{"type": "service_account", "project_id": "test"}`)

	config, err := calendar.LoadConfig()
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}

	if config == nil {
		t.Error("Expected config to be non-nil")
	}

	// Verify defaults
	if config.CalendarID != "primary" {
		t.Errorf("Expected default calendar ID 'primary', got: %s", config.CalendarID)
	}

	if config.TimeZone != "UTC" {
		t.Errorf("Expected default timezone 'UTC', got: %s", config.TimeZone)
	}
}

func TestEnvFileLoading(t *testing.T) {
	// Create a temporary .env file
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	envContent := `
GOOGLE_CALENDAR_CREDENTIALS_JSON={"type": "service_account", "project_id": "test-env"}
GOOGLE_CALENDAR_ID=env-test@example.com
GOOGLE_CALENDAR_TIMEZONE=Europe/London
ENVIRONMENT=test
DEBUG=true
`
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}

	// Clear existing environment variables
	os.Clearenv()

	// Load the .env file
	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	// Test configuration loading
	config, err := calendar.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config from .env: %v", err)
	}

	// Verify values from .env file
	if config.CalendarID != "env-test@example.com" {
		t.Errorf("Expected calendar ID 'env-test@example.com', got: %s", config.CalendarID)
	}

	if config.TimeZone != "Europe/London" {
		t.Errorf("Expected timezone 'Europe/London', got: %s", config.TimeZone)
	}

	if !config.Debug {
		t.Error("Expected debug mode to be true")
	}
}

func TestConfigString(t *testing.T) {
	os.Setenv("GOOGLE_CALENDAR_CREDENTIALS_JSON", `{"type": "service_account"}`)
	os.Setenv("GOOGLE_CALENDAR_ID", "test@example.com")
	os.Setenv("GOOGLE_CALENDAR_TIMEZONE", "America/New_York")
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DEBUG", "true")

	config, err := calendar.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	configStr := config.String()
	if configStr == "" {
		t.Error("Config string should not be empty")
	}

	// Should not contain sensitive information
	if contains(configStr, "service_account") {
		t.Error("Config string should not contain sensitive credential information")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
