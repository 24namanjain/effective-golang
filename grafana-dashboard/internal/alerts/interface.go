package alerts

import (
	"context"
	"fmt"
	"sync"
	"system-monitor/internal/config"
	"system-monitor/internal/datasource"
	"time"

	"github.com/sirupsen/logrus"
)

// AlertState tracks the current state of alerts
type AlertState struct {
	CPUWarning      bool
	CPUCritical     bool
	MemoryWarning   bool
	MemoryCritical  bool
	LatencyWarning  bool
	LatencyCritical bool
}

// AlertBackend defines the interface for different alert backends
type AlertBackend interface {
	// SendAlert sends an alert
	SendAlert(ctx context.Context, alert *Alert) error

	// HealthCheck checks if the alert backend is healthy
	HealthCheck(ctx context.Context) error

	// Close closes the alert backend connection
	Close() error
}

// Alert represents an alert message
type Alert struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Severity  string                 `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AlertManager manages alert processing
type AlertManager struct {
	config    *config.Config
	backend   AlertBackend
	state     AlertState
	lastAlert map[string]time.Time
	mu        sync.RWMutex
	stopChan  chan struct{}
}

// NewAlertManager creates a new alert manager
func NewAlertManager(config *config.Config, backend AlertBackend) *AlertManager {
	return &AlertManager{
		config:    config,
		backend:   backend,
		lastAlert: make(map[string]time.Time),
		stopChan:  make(chan struct{}),
	}
}

// Start begins the alert manager
func (am *AlertManager) Start() error {
	// Send startup alert
	am.sendStartupAlert()
	return nil
}

// Stop stops the alert manager
func (am *AlertManager) Stop() error {
	am.sendShutdownAlert()
	return nil
}

// ProcessMetrics processes metrics and sends alerts if thresholds are exceeded
func (am *AlertManager) ProcessMetrics(metrics *datasource.Metrics) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Check CPU alerts
	am.checkCPUAlerts(metrics)

	// Check memory alerts
	am.checkMemoryAlerts(metrics)

	// Check latency alerts
	am.checkLatencyAlerts(metrics)
}

// GetAlertState returns the current alert state
func (am *AlertManager) GetAlertState() AlertState {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.state
}

// checkCPUAlerts checks CPU usage and sends alerts if needed
func (am *AlertManager) checkCPUAlerts(metrics *datasource.Metrics) {
	if am.config == nil {
		return
	}

	cpuThreshold := am.config.CPUThreshold
	cpuUsage := metrics.CPU

	// Check if we can send an alert (cooldown period)
	alertKey := "cpu_warning"
	if !am.canSendAlert(alertKey) {
		return
	}

	// Check CPU threshold
	if cpuUsage > cpuThreshold {
		// Determine severity
		severity := "warning"
		if cpuUsage > cpuThreshold*1.5 {
			severity = "critical"
		}

		// Create alert
		alert := &Alert{
			ID:        generateAlertID(),
			Type:      "cpu_high_usage",
			Title:     "High CPU Usage Alert",
			Message:   fmt.Sprintf("CPU usage is %.1f%% (threshold: %.1f%%)", cpuUsage, cpuThreshold),
			Severity:  severity,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"cpu_usage": cpuUsage,
				"threshold": cpuThreshold,
				"host":      "localhost",
			},
		}

		// Send alert
		ctx := context.Background()
		if err := am.backend.SendAlert(ctx, alert); err != nil {
			logrus.Errorf("Failed to send CPU alert: %v", err)
			return
		}

		// Update state and mark alert as sent
		am.lastAlert[alertKey] = time.Now()
		if severity == "critical" {
			am.state.CPUCritical = true
		} else {
			am.state.CPUWarning = true
		}

		logrus.Infof("🚨 CPU Alert sent: %.1f%% usage (threshold: %.1f%%)", cpuUsage, cpuThreshold)
	} else {
		// Reset state if CPU usage is back to normal
		am.state.CPUWarning = false
		am.state.CPUCritical = false
	}
}

// checkMemoryAlerts checks memory usage and sends alerts if needed
func (am *AlertManager) checkMemoryAlerts(metrics *datasource.Metrics) {
	if am.config == nil {
		return
	}

	memoryThreshold := am.config.MemoryThreshold
	memoryUsage := metrics.Memory.Percent

	// Check if we can send an alert (cooldown period)
	alertKey := "memory_warning"
	if !am.canSendAlert(alertKey) {
		return
	}

	// Check memory threshold
	if memoryUsage > memoryThreshold {
		// Determine severity
		severity := "warning"
		if memoryUsage > memoryThreshold*1.2 {
			severity = "critical"
		}

		// Create alert
		alert := &Alert{
			ID:        generateAlertID(),
			Type:      "memory_high_usage",
			Title:     "High Memory Usage Alert",
			Message:   fmt.Sprintf("Memory usage is %.1f%% (threshold: %.1f%%)", memoryUsage, memoryThreshold),
			Severity:  severity,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"memory_usage": memoryUsage,
				"threshold":    memoryThreshold,
				"host":         "localhost",
			},
		}

		// Send alert
		ctx := context.Background()
		if err := am.backend.SendAlert(ctx, alert); err != nil {
			logrus.Errorf("Failed to send memory alert: %v", err)
			return
		}

		// Update state and mark alert as sent
		am.lastAlert[alertKey] = time.Now()
		if severity == "critical" {
			am.state.MemoryCritical = true
		} else {
			am.state.MemoryWarning = true
		}

		logrus.Infof("🚨 Memory Alert sent: %.1f%% usage (threshold: %.1f%%)", memoryUsage, memoryThreshold)
	} else {
		// Reset state if memory usage is back to normal
		am.state.MemoryWarning = false
		am.state.MemoryCritical = false
	}
}

// checkLatencyAlerts checks latency and sends alerts if needed
func (am *AlertManager) checkLatencyAlerts(metrics *datasource.Metrics) {
	if am.config == nil {
		return
	}

	latencyThreshold := am.config.LatencyThreshold
	latency := metrics.Latency.HTTPLatency

	// Check if we can send an alert (cooldown period)
	alertKey := "latency_warning"
	if !am.canSendAlert(alertKey) {
		return
	}

	// Check latency threshold
	if latency > latencyThreshold {
		// Determine severity
		severity := "warning"
		if latency > latencyThreshold*2 {
			severity = "critical"
		}

		// Create alert
		alert := &Alert{
			ID:        generateAlertID(),
			Type:      "latency_high",
			Title:     "High Latency Alert",
			Message:   fmt.Sprintf("HTTP latency is %dms (threshold: %dms)", latency, latencyThreshold),
			Severity:  severity,
			Timestamp: time.Now(),
			Metadata: map[string]interface{}{
				"latency":   latency,
				"threshold": latencyThreshold,
				"host":      "localhost",
			},
		}

		// Send alert
		ctx := context.Background()
		if err := am.backend.SendAlert(ctx, alert); err != nil {
			logrus.Errorf("Failed to send latency alert: %v", err)
			return
		}

		// Update state and mark alert as sent
		am.lastAlert[alertKey] = time.Now()
		if severity == "critical" {
			am.state.LatencyCritical = true
		} else {
			am.state.LatencyWarning = true
		}

		logrus.Infof("🚨 Latency Alert sent: %dms (threshold: %dms)", latency, latencyThreshold)
	} else {
		// Reset state if latency is back to normal
		am.state.LatencyWarning = false
		am.state.LatencyCritical = false
	}
}

// sendStartupAlert sends a startup notification
func (am *AlertManager) sendStartupAlert() {
	if am.config == nil {
		return
	}

	alert := &Alert{
		ID:        generateAlertID(),
		Type:      "system_startup",
		Title:     "System Monitor Started",
		Message:   "System monitoring service has started successfully",
		Severity:  "info",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"data_source":   am.config.DataSourceType,
			"alert_backend": am.config.AlertBackendType,
			"host":          "localhost",
		},
	}

	ctx := context.Background()
	if err := am.backend.SendAlert(ctx, alert); err != nil {
		logrus.Errorf("Failed to send startup alert: %v", err)
	}
}

// sendShutdownAlert sends a shutdown notification
func (am *AlertManager) sendShutdownAlert() {
	if am.config == nil {
		return
	}

	alert := &Alert{
		ID:        generateAlertID(),
		Type:      "system_shutdown",
		Title:     "System Monitor Stopping",
		Message:   "System monitoring service is shutting down",
		Severity:  "info",
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"host": "localhost",
		},
	}

	ctx := context.Background()
	if err := am.backend.SendAlert(ctx, alert); err != nil {
		logrus.Errorf("Failed to send shutdown alert: %v", err)
	}
}

// canSendAlert checks if enough time has passed since the last alert
func (am *AlertManager) canSendAlert(alertKey string) bool {
	if am.config == nil {
		return false
	}

	lastAlertTime, exists := am.lastAlert[alertKey]
	if !exists {
		return true
	}

	cooldown := am.config.AlertCooldown
	return time.Since(lastAlertTime) > cooldown
}

// generateAlertID generates a unique alert ID
func generateAlertID() string {
	return fmt.Sprintf("alert-%d", time.Now().UnixNano())
}
