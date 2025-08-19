package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	SlackClientID     string
	SlackAppID        string
	SlackClientSecret string
	SlackSigningSecret string
	SlackAppLevelToken string
	SlackBotToken     string
	SlackChannel      string
	Environment       string
	APIAddress        string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		SlackClientID:      getEnv("SLACK_CLIENT_ID", "9371894822341.9375537877046"),
		SlackAppID:         getEnv("SLACK_APP_ID", "A09B1FTRT1C"),
		SlackClientSecret:  getEnv("SLACK_CLIENT_SECRET", "1b5dc7da4540da29a5e02ec9bbaf69e5"),
		SlackSigningSecret: getEnv("SLACK_SIGNING_SECRET", ""),
		SlackAppLevelToken: getEnv("SLACK_APP_LEVEL_TOKEN", ""),
		SlackBotToken:      getEnv("SLACK_BOT_TOKEN", ""),
		SlackChannel:       getEnv("SLACK_CHANNEL", "#general"),
		Environment:        getEnv("ENVIRONMENT", "development"),
		APIAddress:         getEnv("API_ADDR", ":8081"),
	}

	// Validate required fields and common misconfigurations
	if config.SlackBotToken == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN is required (must be a Bot User OAuth Token starting with 'xoxb-')")
	}
	if strings.HasPrefix(config.SlackBotToken, "xapp-") {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN is an app-level token (xapp-). Use a Bot User OAuth Token (xoxb-) from OAuth & Permissions after installing the app")
	}
	if !(strings.HasPrefix(config.SlackBotToken, "xoxb-") || strings.HasPrefix(config.SlackBotToken, "xoxp-")) {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN must start with 'xoxb-' (bot) or 'xoxp-' (user). Recommended: bot token 'xoxb-'")
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
