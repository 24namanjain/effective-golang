package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test with required environment variables
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-test-token")
	os.Setenv("SLACK_CHANNEL", "#test-channel")
	os.Setenv("ENVIRONMENT", "test")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.SlackBotToken != "xoxb-test-token" {
		t.Errorf("Expected SlackBotToken to be 'xoxb-test-token', got %s", config.SlackBotToken)
	}

	if config.SlackChannel != "#test-channel" {
		t.Errorf("Expected SlackChannel to be '#test-channel', got %s", config.SlackChannel)
	}

	if config.Environment != "test" {
		t.Errorf("Expected Environment to be 'test', got %s", config.Environment)
	}

	// Test default values
	if config.SlackClientID != "9371894822341.9375537877046" {
		t.Errorf("Expected default SlackClientID, got %s", config.SlackClientID)
	}

	if config.SlackAppID != "A09B1FTRT1C" {
		t.Errorf("Expected default SlackAppID, got %s", config.SlackAppID)
	}

	if config.APIAddress != ":8081" {
		t.Errorf("Expected default APIAddress to be ':8081', got %s", config.APIAddress)
	}
}

func TestLoadConfigMissingBotToken(t *testing.T) {
	// Clear the bot token
	os.Unsetenv("SLACK_BOT_TOKEN")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for missing SLACK_BOT_TOKEN, got nil")
	}

	expectedError := "SLACK_BOT_TOKEN is required (must be a Bot User OAuth Token starting with 'xoxb-')"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfigAppLevelToken(t *testing.T) {
	// Test that app-level tokens are rejected
	os.Setenv("SLACK_BOT_TOKEN", "xapp-test-token")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for app-level token, got nil")
	}

	expectedError := "SLACK_BOT_TOKEN is an app-level token (xapp-). Use a Bot User OAuth Token (xoxb-) from OAuth & Permissions after installing the app"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestLoadConfigInvalidToken(t *testing.T) {
	// Test that invalid tokens are rejected
	os.Setenv("SLACK_BOT_TOKEN", "invalid-token")

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error for invalid token, got nil")
	}

	expectedError := "SLACK_BOT_TOKEN must start with 'xoxb-' (bot) or 'xoxp-' (user). Recommended: bot token 'xoxb-'"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnvironmentHelpers(t *testing.T) {
	config := &Config{Environment: "development"}

	if !config.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for 'development' environment")
	}

	if config.IsProduction() {
		t.Error("Expected IsProduction() to return false for 'development' environment")
	}

	config.Environment = "production"

	if config.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false for 'production' environment")
	}

	if !config.IsProduction() {
		t.Error("Expected IsProduction() to return true for 'production' environment")
	}
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_VAR", "test-value")
	result := getEnv("TEST_VAR", "default")
	if result != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", result)
	}

	// Test with environment variable not set
	result = getEnv("NONEXISTENT_VAR", "default-value")
	if result != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", result)
	}
}

func TestValidTokenFormats(t *testing.T) {
	testCases := []struct {
		token   string
		isValid bool
	}{
		{"xoxb-valid-token", true},
		{"xoxp-valid-token", true},
		{"xapp-invalid-token", false},
		{"invalid-token", false},
		{"", false},
	}

	for _, tc := range testCases {
		os.Setenv("SLACK_BOT_TOKEN", tc.token)
		_, err := LoadConfig()

		if tc.isValid && err != nil {
			t.Errorf("Token '%s' should be valid but got error: %v", tc.token, err)
		}
		if !tc.isValid && err == nil {
			t.Errorf("Token '%s' should be invalid but got no error", tc.token)
		}
	}
}
