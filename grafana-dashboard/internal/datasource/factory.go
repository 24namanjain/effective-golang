package datasource

import (
	"fmt"
	"system-monitor/internal/config"
)

// Factory creates data source instances
type Factory struct{}

// NewFactory creates a new data source factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateDataSource creates a data source based on configuration
func (f *Factory) CreateDataSource(cfg *config.Config) (DataSource, error) {
	switch cfg.DataSourceType {
	case config.DataSourceLocal:
		return f.createLocalDataSource(cfg)
	case config.DataSourceGrafana:
		return f.createGrafanaDataSource(cfg)
	case config.DataSourcePrometheus:
		return f.createPrometheusDataSource(cfg)
	default:
		return nil, fmt.Errorf("unsupported data source type: %s", cfg.DataSourceType)
	}
}

// createLocalDataSource creates a local data source
func (f *Factory) createLocalDataSource(cfg *config.Config) (DataSource, error) {
	return NewLocalDataSource(1000), nil // Keep last 1000 metrics
}

// createGrafanaDataSource creates a Grafana data source
func (f *Factory) createGrafanaDataSource(cfg *config.Config) (DataSource, error) {
	dsConfig := NewDataSourceConfig(DataSourceGrafana, cfg.GrafanaURL)

	if cfg.GrafanaAPIKey != "" {
		dsConfig.WithAPIKey(cfg.GrafanaAPIKey)
	} else if cfg.GrafanaUsername != "" && cfg.GrafanaPassword != "" {
		dsConfig.WithCredentials(cfg.GrafanaUsername, cfg.GrafanaPassword)
	}

	return NewGrafanaDataSource(dsConfig), nil
}

// createPrometheusDataSource creates a Prometheus data source
func (f *Factory) createPrometheusDataSource(cfg *config.Config) (DataSource, error) {
	// For now, return a simplified implementation
	// In a real implementation, you'd create a PrometheusDataSource
	dsConfig := NewDataSourceConfig(DataSourcePrometheus, cfg.PrometheusURL)
	return NewGrafanaDataSource(dsConfig), nil // Reuse Grafana implementation for now
}
