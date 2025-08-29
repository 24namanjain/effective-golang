package alerts

import (
	"context"

	"github.com/sirupsen/logrus"
)

// NoOpAlertBackend implements AlertBackend for testing (does nothing)
type NoOpAlertBackend struct{}

// NewNoOpAlertBackend creates a new no-op alert backend
func NewNoOpAlertBackend() *NoOpAlertBackend {
	return &NoOpAlertBackend{}
}

// SendAlert logs the alert but doesn't send it anywhere
func (n *NoOpAlertBackend) SendAlert(ctx context.Context, alert *Alert) error {
	logrus.Infof("🚨 [NO-OP] Alert would be sent: %s - %s (Severity: %s)",
		alert.Type, alert.Title, alert.Severity)

	if len(alert.Metadata) > 0 {
		logrus.Infof("📋 Alert metadata: %+v", alert.Metadata)
	}

	return nil
}

// HealthCheck always returns healthy
func (n *NoOpAlertBackend) HealthCheck(ctx context.Context) error {
	return nil
}

// Close does nothing
func (n *NoOpAlertBackend) Close() error {
	return nil
}
