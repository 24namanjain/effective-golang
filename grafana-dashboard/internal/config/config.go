package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// DataSourceType represents the type of data source
type DataSourceType string

const (
	DataSourceLocal      DataSourceType = "local"
	DataSourcePrometheus DataSourceType = "prometheus"
	DataSourceGrafana    DataSourceType = "grafana"
)

// AlertBackendType represents the type of alerting backend
type AlertBackendType string

const (
	AlertBackendSlack   AlertBackendType = "slack"
	AlertBackendEmail   AlertBackendType = "email"
	AlertBackendWebhook AlertBackendType = "webhook"
)

// Config holds all configuration for the monitoring system
type Config struct {
	// Data Source Configuration
	DataSourceType DataSourceType
	DataSourceURL  string

	// Grafana Configuration
	GrafanaURL      string
	GrafanaUsername string
	GrafanaPassword string
	GrafanaAPIKey   string

	// Prometheus Configuration
	PrometheusURL string

	// Alert Configuration
	AlertBackendType AlertBackendType
	SlackBotToken    string
	SlackChannel     string
	WebhookURL       string

	// Threshold Configuration
	CPUThreshold     float64
	MemoryThreshold  float64
	LatencyThreshold int64

	// Alert Settings
	AlertCooldown time.Duration

	// Dashboard Settings
	DashboardPort string

	// Metrics Collection
	MetricsInterval time.Duration

	// Environment
	Environment string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		DataSourceType:   getDataSourceType("DATA_SOURCE_TYPE", DataSourceLocal),
		DataSourceURL:    getEnv("DATA_SOURCE_URL", "http://localhost:9090"),
		GrafanaURL:       getEnv("GRAFANA_URL", "http://localhost:3000"),
		GrafanaUsername:  getEnv("GRAFANA_USERNAME", "admin"),
		GrafanaPassword:  getEnv("GRAFANA_PASSWORD", "admin123"),
		GrafanaAPIKey:    getEnv("GRAFANA_API_KEY", ""),
		PrometheusURL:    getEnv("PROMETHEUS_URL", "http://localhost:9090"),
		AlertBackendType: getAlertBackendType("ALERT_BACKEND_TYPE", AlertBackendSlack),
		SlackBotToken:    getEnv("SLACK_BOT_TOKEN", ""),
		SlackChannel:     getEnv("SLACK_CHANNEL", "#alerts"),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		CPUThreshold:     getEnvAsFloat("CPU_THRESHOLD", 80.0),
		MemoryThreshold:  getEnvAsFloat("MEMORY_THRESHOLD", 85.0),
		LatencyThreshold: getEnvAsInt64("LATENCY_THRESHOLD", 500),
		AlertCooldown:    getEnvAsDuration("ALERT_COOLDOWN", 5*time.Minute),
		DashboardPort:    getEnv("DASHBOARD_PORT", "8080"),
		MetricsInterval:  getEnvAsDuration("METRICS_INTERVAL", 5*time.Second),
		Environment:      getEnv("ENVIRONMENT", "development"),
	}

	// Validate configuration based on data source type
	if err := config.validateDataSourceConfig(); err != nil {
		return nil, err
	}

	// Validate alert backend configuration
	if err := config.validateAlertBackendConfig(); err != nil {
		return nil, err
	}

	// Validate thresholds
	if err := config.validateThresholds(); err != nil {
		return nil, err
	}

	return config, nil
}

// validateDataSourceConfig validates data source specific configuration
func (c *Config) validateDataSourceConfig() error {
	switch c.DataSourceType {
	case DataSourceGrafana:
		if c.GrafanaURL == "" {
			return fmt.Errorf("GRAFANA_URL is required when using grafana data source")
		}
		if c.GrafanaAPIKey == "" && (c.GrafanaUsername == "" || c.GrafanaPassword == "") {
			return fmt.Errorf("either GRAFANA_API_KEY or GRAFANA_USERNAME/GRAFANA_PASSWORD is required")
		}
	case DataSourcePrometheus:
		if c.PrometheusURL == "" {
			return fmt.Errorf("PROMETHEUS_URL is required when using prometheus data source")
		}
	}
	return nil
}

// validateAlertBackendConfig validates alert backend specific configuration
func (c *Config) validateAlertBackendConfig() error {
	switch c.AlertBackendType {
	case AlertBackendSlack:
		if c.SlackBotToken == "" {
			return fmt.Errorf("SLACK_BOT_TOKEN is required when using slack alert backend")
		}
		if !strings.HasPrefix(c.SlackBotToken, "xoxb-") {
			return fmt.Errorf("SLACK_BOT_TOKEN must start with 'xoxb-'")
		}
	case AlertBackendWebhook:
		if c.WebhookURL == "" {
			return fmt.Errorf("WEBHOOK_URL is required when using webhook alert backend")
		}
	}
	return nil
}

// validateThresholds validates threshold values
func (c *Config) validateThresholds() error {
	if c.CPUThreshold < 0 || c.CPUThreshold > 100 {
		return fmt.Errorf("CPU_THRESHOLD must be between 0 and 100")
	}
	if c.MemoryThreshold < 0 || c.MemoryThreshold > 100 {
		return fmt.Errorf("MEMORY_THRESHOLD must be between 0 and 100")
	}
	if c.LatencyThreshold < 0 {
		return fmt.Errorf("LATENCY_THRESHOLD must be positive")
	}
	return nil
}

// getDataSourceType gets data source type from environment
func getDataSourceType(key string, fallback DataSourceType) DataSourceType {
	if value := os.Getenv(key); value != "" {
		switch DataSourceType(value) {
		case DataSourceLocal, DataSourcePrometheus, DataSourceGrafana:
			return DataSourceType(value)
		}
	}
	return fallback
}

// getAlertBackendType gets alert backend type from environment
func getAlertBackendType(key string, fallback AlertBackendType) AlertBackendType {
	if value := os.Getenv(key); value != "" {
		switch AlertBackendType(value) {
		case AlertBackendSlack, AlertBackendEmail, AlertBackendWebhook, "noop":
			return AlertBackendType(value)
		}
	}
	return fallback
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsFloat gets an environment variable as float64 with a fallback default value
func getEnvAsFloat(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

// getEnvAsInt64 gets an environment variable as int64 with a fallback default value
func getEnvAsInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as time.Duration with a fallback default value
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
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

// GetDashboardURL returns the full dashboard URL
func (c *Config) GetDashboardURL() string {
	return fmt.Sprintf("http://localhost:%s", c.DashboardPort)
}

// IsGrafanaEnabled returns true if Grafana integration is enabled
func (c *Config) IsGrafanaEnabled() bool {
	return c.DataSourceType == DataSourceGrafana
}

// IsPrometheusEnabled returns true if Prometheus integration is enabled
func (c *Config) IsPrometheusEnabled() bool {
	return c.DataSourceType == DataSourcePrometheus
}

// IsLocalDataSource returns true if using local data source
func (c *Config) IsLocalDataSource() bool {
	return c.DataSourceType == DataSourceLocal
}
