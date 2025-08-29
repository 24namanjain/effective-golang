package alerts

import (
	"fmt"
	"system-monitor/internal/config"
)

// AlertBackendFactory creates alert backend instances
type AlertBackendFactory struct{}

// NewAlertBackendFactory creates a new alert backend factory
func NewAlertBackendFactory() *AlertBackendFactory {
	return &AlertBackendFactory{}
}

// CreateAlertBackend creates an alert backend based on configuration
func (f *AlertBackendFactory) CreateAlertBackend(cfg *config.Config) (AlertBackend, error) {
	switch cfg.AlertBackendType {
	case config.AlertBackendSlack:
		return f.createSlackAlertBackend(cfg)
	case config.AlertBackendWebhook:
		return f.createWebhookAlertBackend(cfg)
	case config.AlertBackendEmail:
		return f.createEmailAlertBackend(cfg)
	case "noop":
		return f.createNoOpAlertBackend(cfg)
	default:
		return nil, fmt.Errorf("unsupported alert backend type: %s", cfg.AlertBackendType)
	}
}

// createSlackAlertBackend creates a Slack alert backend
func (f *AlertBackendFactory) createSlackAlertBackend(cfg *config.Config) (AlertBackend, error) {
	if cfg.SlackBotToken == "" {
		return nil, fmt.Errorf("SLACK_BOT_TOKEN is required for Slack alert backend")
	}

	return NewSlackAlertBackend(cfg.SlackBotToken, cfg.SlackChannel)
}

// createWebhookAlertBackend creates a webhook alert backend
func (f *AlertBackendFactory) createWebhookAlertBackend(cfg *config.Config) (AlertBackend, error) {
	// Implementation for webhook backend
	return nil, fmt.Errorf("webhook alert backend not implemented yet")
}

// createEmailAlertBackend creates an email alert backend
func (f *AlertBackendFactory) createEmailAlertBackend(cfg *config.Config) (AlertBackend, error) {
	// Implementation for email backend
	return nil, fmt.Errorf("email alert backend not implemented yet")
}

// createNoOpAlertBackend creates a no-op alert backend for testing
func (f *AlertBackendFactory) createNoOpAlertBackend(cfg *config.Config) (AlertBackend, error) {
	return NewNoOpAlertBackend(), nil
}
