package datasource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// GrafanaDataSource implements DataSource for Grafana
type GrafanaDataSource struct {
	config     *DataSourceConfig
	httpClient *http.Client
}

// NewGrafanaDataSource creates a new Grafana data source
func NewGrafanaDataSource(config *DataSourceConfig) *GrafanaDataSource {
	return &GrafanaDataSource{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// GetLatestMetrics returns the most recent metrics from Grafana
func (ds *GrafanaDataSource) GetLatestMetrics(ctx context.Context) (*Metrics, error) {
	// Query Grafana API for latest metrics
	// This is a simplified implementation - in practice, you'd query specific panels/dashboards

	// For now, we'll simulate getting metrics from Grafana
	// In a real implementation, you'd make API calls to Grafana's query endpoints

	cpuQuery := fmt.Sprintf("%s/api/datasources/proxy/1/api/v1/query?query=100%%20-%%20(avg%%20by%%20(instance)%%20(irate(node_cpu_seconds_total{mode=\"idle\"}[5m]))%%20*%%20100)", ds.config.URL)
	memoryQuery := fmt.Sprintf("%s/api/datasources/proxy/1/api/v1/query?query=(node_memory_MemTotal_bytes%%20-%%20node_memory_MemAvailable_bytes)%%20/%%20node_memory_MemTotal_bytes%%20*%%20100", ds.config.URL)

	// Simulate API calls (in real implementation, you'd make actual HTTP requests)
	cpuUsage, err := ds.queryGrafanaMetric(ctx, cpuQuery)
	if err != nil {
		logrus.Warnf("Failed to query CPU metric: %v", err)
		cpuUsage = 0
	}

	memoryUsage, err := ds.queryGrafanaMetric(ctx, memoryQuery)
	if err != nil {
		logrus.Warnf("Failed to query memory metric: %v", err)
		memoryUsage = 0
	}

	// Simulate latency metric
	latencyUsage := int64(50 + (time.Now().UnixNano() % 200))

	return &Metrics{
		Timestamp: time.Now(),
		CPU:       cpuUsage,
		Memory: MemoryInfo{
			Percent: memoryUsage,
		},
		Latency: LatencyInfo{
			HTTPLatency: latencyUsage,
		},
	}, nil
}

// GetMetricsHistory returns historical metrics from Grafana
func (ds *GrafanaDataSource) GetMetricsHistory(ctx context.Context, start, end time.Time) ([]*Metrics, error) {
	// Query Grafana API for historical metrics
	// This would involve making range queries to Grafana's API

	// For now, return a simplified implementation
	var metrics []*Metrics
	duration := end.Sub(start)
	points := int(duration.Minutes()) / 5 // 5-minute intervals

	for i := 0; i < points; i++ {
		timestamp := start.Add(time.Duration(i*5) * time.Minute)
		metrics = append(metrics, &Metrics{
			Timestamp: timestamp,
			CPU:       float64(30 + (i % 70)), // Simulated CPU usage
			Memory: MemoryInfo{
				Percent: float64(40 + (i % 60)), // Simulated memory usage
			},
			Latency: LatencyInfo{
				HTTPLatency: int64(50 + (i % 200)), // Simulated latency
			},
		})
	}

	return metrics, nil
}

// GetCPUHistory returns CPU usage history from Grafana
func (ds *GrafanaDataSource) GetCPUHistory(ctx context.Context, duration time.Duration) ([]float64, error) {
	end := time.Now()
	start := end.Add(-duration)

	metrics, err := ds.GetMetricsHistory(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var result []float64
	for _, metric := range metrics {
		result = append(result, metric.CPU)
	}

	return result, nil
}

// GetMemoryHistory returns memory usage history from Grafana
func (ds *GrafanaDataSource) GetMemoryHistory(ctx context.Context, duration time.Duration) ([]float64, error) {
	end := time.Now()
	start := end.Add(-duration)

	metrics, err := ds.GetMetricsHistory(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var result []float64
	for _, metric := range metrics {
		result = append(result, metric.Memory.Percent)
	}

	return result, nil
}

// GetLatencyHistory returns latency history from Grafana
func (ds *GrafanaDataSource) GetLatencyHistory(ctx context.Context, duration time.Duration) ([]int64, error) {
	end := time.Now()
	start := end.Add(-duration)

	metrics, err := ds.GetMetricsHistory(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, metric := range metrics {
		result = append(result, metric.Latency.HTTPLatency)
	}

	return result, nil
}

// GetTimestamps returns timestamps for metrics from Grafana
func (ds *GrafanaDataSource) GetTimestamps(ctx context.Context, duration time.Duration) ([]time.Time, error) {
	end := time.Now()
	start := end.Add(-duration)

	metrics, err := ds.GetMetricsHistory(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var result []time.Time
	for _, metric := range metrics {
		result = append(result, metric.Timestamp)
	}

	return result, nil
}

// HealthCheck checks if the Grafana data source is healthy
func (ds *GrafanaDataSource) HealthCheck(ctx context.Context) error {
	healthURL := fmt.Sprintf("%s/api/health", ds.config.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Add authentication if configured
	if ds.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+ds.config.APIKey)
	} else if ds.config.Username != "" && ds.config.Password != "" {
		req.SetBasicAuth(ds.config.Username, ds.config.Password)
	}

	resp, err := ds.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check Grafana health: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Grafana health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Close closes the data source connection
func (ds *GrafanaDataSource) Close() error {
	// No specific cleanup needed for HTTP client
	return nil
}

// queryGrafanaMetric queries a specific metric from Grafana
func (ds *GrafanaDataSource) queryGrafanaMetric(ctx context.Context, queryURL string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", queryURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured
	if ds.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+ds.config.APIKey)
	} else if ds.config.Username != "" && ds.config.Password != "" {
		req.SetBasicAuth(ds.config.Username, ds.config.Password)
	}

	resp, err := ds.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to query Grafana: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Grafana query failed with status: %d", resp.StatusCode)
	}

	// Parse the response (simplified - in practice, you'd parse the actual Prometheus response format)
	var response struct {
		Data struct {
			Result []struct {
				Value []interface{} `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Data.Result) == 0 || len(response.Data.Result[0].Value) < 2 {
		return 0, fmt.Errorf("no data in response")
	}

	// Extract the metric value
	valueStr, ok := response.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("invalid value format in response")
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse metric value: %w", err)
	}

	return value, nil
}
