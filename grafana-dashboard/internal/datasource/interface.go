package datasource

import (
	"context"
	"time"
)

// Metrics represents system metrics at a point in time
type Metrics struct {
	Timestamp time.Time   `json:"timestamp"`
	CPU       float64     `json:"cpu"`
	Memory    MemoryInfo  `json:"memory"`
	Latency   LatencyInfo `json:"latency"`
}

// MemoryInfo contains memory-related metrics
type MemoryInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Percent   float64 `json:"percent"`
}

// LatencyInfo contains latency-related metrics
type LatencyInfo struct {
	HTTPLatency int64 `json:"http_latency"`
	DBLatency   int64 `json:"db_latency"`
	APILatency  int64 `json:"api_latency"`
}

// DataSource defines the interface for different data sources
type DataSource interface {
	// GetLatestMetrics returns the most recent metrics
	GetLatestMetrics(ctx context.Context) (*Metrics, error)

	// GetMetricsHistory returns historical metrics within a time range
	GetMetricsHistory(ctx context.Context, start, end time.Time) ([]*Metrics, error)

	// GetCPUHistory returns CPU usage history
	GetCPUHistory(ctx context.Context, duration time.Duration) ([]float64, error)

	// GetMemoryHistory returns memory usage history
	GetMemoryHistory(ctx context.Context, duration time.Duration) ([]float64, error)

	// GetLatencyHistory returns latency history
	GetLatencyHistory(ctx context.Context, duration time.Duration) ([]int64, error)

	// GetTimestamps returns timestamps for metrics
	GetTimestamps(ctx context.Context, duration time.Duration) ([]time.Time, error)

	// HealthCheck checks if the data source is healthy
	HealthCheck(ctx context.Context) error

	// Close closes the data source connection
	Close() error
}

// DataSourceType represents the type of data source
type DataSourceType string

const (
	DataSourceLocal      DataSourceType = "local"
	DataSourcePrometheus DataSourceType = "prometheus"
	DataSourceGrafana    DataSourceType = "grafana"
)

// DataSourceFactory creates data source instances
type DataSourceFactory interface {
	CreateDataSource(dataSourceType DataSourceType, config map[string]interface{}) (DataSource, error)
}

// DataSourceConfig holds configuration for data sources
type DataSourceConfig struct {
	Type     DataSourceType
	URL      string
	Username string
	Password string
	APIKey   string
	Timeout  time.Duration
}

// NewDataSourceConfig creates a new data source configuration
func NewDataSourceConfig(dataSourceType DataSourceType, url string) *DataSourceConfig {
	return &DataSourceConfig{
		Type:    dataSourceType,
		URL:     url,
		Timeout: 30 * time.Second,
	}
}

// WithCredentials adds authentication credentials
func (c *DataSourceConfig) WithCredentials(username, password string) *DataSourceConfig {
	c.Username = username
	c.Password = password
	return c
}

// WithAPIKey adds API key authentication
func (c *DataSourceConfig) WithAPIKey(apiKey string) *DataSourceConfig {
	c.APIKey = apiKey
	return c
}

// WithTimeout sets the timeout for data source operations
func (c *DataSourceConfig) WithTimeout(timeout time.Duration) *DataSourceConfig {
	c.Timeout = timeout
	return c
}
