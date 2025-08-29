package alerts

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// SlackAlertBackend implements AlertBackend for Slack
type SlackAlertBackend struct {
	client  *slack.Client
	channel string
}

// NewSlackAlertBackend creates a new Slack alert backend
func NewSlackAlertBackend(token, channel string) (*SlackAlertBackend, error) {
	client := slack.New(token)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.AuthTestContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Slack: %w", err)
	}

	return &SlackAlertBackend{
		client:  client,
		channel: channel,
	}, nil
}

// SendAlert sends an alert to Slack
func (sab *SlackAlertBackend) SendAlert(ctx context.Context, alert *Alert) error {
	// Build the message
	text := fmt.Sprintf("%s *%s*\n%s", getEmoji(alert.Severity), alert.Title, alert.Message)

	// Add metadata
	if len(alert.Metadata) > 0 {
		text += "\n\n*Details:*\n"
		for key, value := range alert.Metadata {
			text += fmt.Sprintf("• %s: %v\n", key, value)
		}
	}

	text += fmt.Sprintf("\n*Timestamp:* %s", alert.Timestamp.Format(time.RFC3339))

	// Send to Slack
	_, _, err := sab.client.PostMessageContext(ctx, sab.channel,
		slack.MsgOptionText(text, false),
		slack.MsgOptionUsername("System Monitor"),
		slack.MsgOptionIconEmoji(":computer:"),
	)

	if err != nil {
		return fmt.Errorf("failed to send Slack message: %w", err)
	}

	logrus.Infof("Sent Slack alert: %s - %s", alert.Severity, alert.Title)
	return nil
}

// HealthCheck checks if the Slack backend is healthy
func (sab *SlackAlertBackend) HealthCheck(ctx context.Context) error {
	_, err := sab.client.AuthTestContext(ctx)
	return err
}

// Close closes the Slack backend connection
func (sab *SlackAlertBackend) Close() error {
	// No specific cleanup needed for Slack client
	return nil
}

// getEmoji returns the appropriate emoji for the alert severity
func getEmoji(severity string) string {
	switch severity {
	case "critical":
		return "🚨"
	case "warning":
		return "⚠️"
	case "info":
		return "ℹ️"
	case "success":
		return "✅"
	default:
		return "ℹ️"
	}
}
