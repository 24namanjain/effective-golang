package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"system-monitor/internal/alerts"
	"system-monitor/internal/config"
	"system-monitor/internal/dashboard"
	"system-monitor/internal/datasource"

	"github.com/sirupsen/logrus"
)

func main() {
	// Set up logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("🚀 Starting System Monitor with Pluggable Architecture")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	logrus.Infof("Configuration loaded - Data Source: %s, Alert Backend: %s",
		cfg.DataSourceType, cfg.AlertBackendType)

	// Create data source factory and data source
	dataSourceFactory := datasource.NewFactory()
	dataSource, err := dataSourceFactory.CreateDataSource(cfg)
	if err != nil {
		logrus.Fatalf("Failed to create data source: %v", err)
	}

	// Create alert backend factory and alert backend
	alertBackendFactory := alerts.NewAlertBackendFactory()
	alertBackend, err := alertBackendFactory.CreateAlertBackend(cfg)
	if err != nil {
		logrus.Fatalf("Failed to create alert backend: %v", err)
	}

	// Create alert manager
	alertManager := alerts.NewAlertManager(cfg, alertBackend)

	// Create dashboard server
	dashboardServer := dashboard.NewServer(cfg, dataSource, alertManager)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start alert manager
	if err := alertManager.Start(); err != nil {
		logrus.Fatalf("Failed to start alert manager: %v", err)
	}

	// Start data source collection if it's a local data source
	if localDS, ok := dataSource.(*datasource.LocalDataSource); ok {
		go func() {
			localDS.Start(ctx, cfg.MetricsInterval)
		}()
	}

	// Start alert processing in a goroutine
	go func() {
		ticker := time.NewTicker(cfg.MetricsInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				latestMetrics, err := dataSource.GetLatestMetrics(ctx)
				if err != nil {
					logrus.Errorf("Failed to get latest metrics: %v", err)
					continue
				}
				if latestMetrics != nil {
					alertManager.ProcessMetrics(latestMetrics)
				}
			}
		}
	}()

	// Start dashboard server in a goroutine
	go func() {
		if err := dashboardServer.Start(); err != nil {
			logrus.Errorf("Dashboard server error: %v", err)
		}
	}()

	logrus.Infof("✅ System Monitor is running!")
	logrus.Infof("📊 Dashboard available at: %s", cfg.GetDashboardURL())
	logrus.Infof("🔔 Alert backend: %s", cfg.AlertBackendType)
	logrus.Infof("📈 Data source: %s", cfg.DataSourceType)
	logrus.Infof("⏰ Metrics collection interval: %v", cfg.MetricsInterval)
	logrus.Infof("⏰ Alert cooldown period: %v", cfg.AlertCooldown)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logrus.Info("🛑 Shutting down System Monitor...")

	// Graceful shutdown
	cancel()

	// Stop alert manager
	if err := alertManager.Stop(); err != nil {
		logrus.Errorf("Error stopping alert manager: %v", err)
	}

	// Close data source
	if err := dataSource.Close(); err != nil {
		logrus.Errorf("Error closing data source: %v", err)
	}

	// Close alert backend
	if err := alertBackend.Close(); err != nil {
		logrus.Errorf("Error closing alert backend: %v", err)
	}

	logrus.Info("✅ System Monitor shutdown complete")
}
