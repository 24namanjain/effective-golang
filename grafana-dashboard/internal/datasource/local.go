package datasource

import (
	"context"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
)

// LocalDataSource implements DataSource for local system metrics
type LocalDataSource struct {
	metrics         []*Metrics
	maxHistory      int
	mu              sync.RWMutex
	stopChan        chan struct{}
	latencyMeasurer *LatencyMeasurer
}

// NewLocalDataSource creates a new local data source
func NewLocalDataSource(maxHistory int) *LocalDataSource {
	return &LocalDataSource{
		metrics:         make([]*Metrics, 0, maxHistory),
		maxHistory:      maxHistory,
		stopChan:        make(chan struct{}),
		latencyMeasurer: NewLatencyMeasurer(),
	}
}

// Start begins collecting metrics at the specified interval
func (ds *LocalDataSource) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logrus.Info("Starting local metrics collection")

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Stopping local metrics collection")
			return
		case <-ds.stopChan:
			logrus.Info("Stopping local metrics collection")
			return
		case <-ticker.C:
			metrics := ds.collectMetrics()
			ds.addMetrics(metrics)
		}
	}
}

// Stop stops the metrics collection
func (ds *LocalDataSource) Stop() {
	close(ds.stopChan)
}

// GetLatestMetrics returns the most recent metrics
func (ds *LocalDataSource) GetLatestMetrics(ctx context.Context) (*Metrics, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if len(ds.metrics) == 0 {
		return nil, nil
	}

	return ds.metrics[len(ds.metrics)-1], nil
}

// GetMetricsHistory returns historical metrics within a time range
func (ds *LocalDataSource) GetMetricsHistory(ctx context.Context, start, end time.Time) ([]*Metrics, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	var result []*Metrics
	for _, metric := range ds.metrics {
		if metric.Timestamp.After(start) && metric.Timestamp.Before(end) {
			result = append(result, metric)
		}
	}
	return result, nil
}

// GetCPUHistory returns CPU usage history
func (ds *LocalDataSource) GetCPUHistory(ctx context.Context, duration time.Duration) ([]float64, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var result []float64

	for _, metric := range ds.metrics {
		if metric.Timestamp.After(cutoff) {
			result = append(result, metric.CPU)
		}
	}
	return result, nil
}

// GetMemoryHistory returns memory usage history
func (ds *LocalDataSource) GetMemoryHistory(ctx context.Context, duration time.Duration) ([]float64, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var result []float64

	for _, metric := range ds.metrics {
		if metric.Timestamp.After(cutoff) {
			result = append(result, metric.Memory.Percent)
		}
	}
	return result, nil
}

// GetLatencyHistory returns latency history
func (ds *LocalDataSource) GetLatencyHistory(ctx context.Context, duration time.Duration) ([]int64, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var result []int64

	for _, metric := range ds.metrics {
		if metric.Timestamp.After(cutoff) {
			result = append(result, metric.Latency.HTTPLatency)
		}
	}
	return result, nil
}

// GetTimestamps returns timestamps for metrics
func (ds *LocalDataSource) GetTimestamps(ctx context.Context, duration time.Duration) ([]time.Time, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var result []time.Time

	for _, metric := range ds.metrics {
		if metric.Timestamp.After(cutoff) {
			result = append(result, metric.Timestamp)
		}
	}
	return result, nil
}

// HealthCheck checks if the data source is healthy
func (ds *LocalDataSource) HealthCheck(ctx context.Context) error {
	// Local data source is always healthy if it's running
	return nil
}

// Close closes the data source connection
func (ds *LocalDataSource) Close() error {
	ds.Stop()
	return nil
}

// collectMetrics gathers current system metrics
func (ds *LocalDataSource) collectMetrics() *Metrics {
	now := time.Now()

	// Collect CPU metrics
	cpuPercent, err := cpu.Percent(0, false)
	var cpuUsage float64
	if err != nil {
		logrus.Errorf("Failed to collect CPU metrics: %v", err)
		cpuUsage = 0
	} else if len(cpuPercent) > 0 {
		cpuUsage = cpuPercent[0]
	}

	// Collect memory metrics
	memInfo, err := mem.VirtualMemory()
	var memoryInfo MemoryInfo
	if err != nil {
		logrus.Errorf("Failed to collect memory metrics: %v", err)
		memoryInfo = MemoryInfo{}
	} else {
		memoryInfo = MemoryInfo{
			Total:     memInfo.Total,
			Used:      memInfo.Used,
			Available: memInfo.Available,
			Percent:   memInfo.UsedPercent,
		}
	}

	// Collect latency metrics
	latencyInfo := ds.latencyMeasurer.MeasureLatency()

	return &Metrics{
		Timestamp: now,
		CPU:       cpuUsage,
		Memory:    memoryInfo,
		Latency:   latencyInfo,
	}
}

// addMetrics adds new metrics to the history
func (ds *LocalDataSource) addMetrics(metrics *Metrics) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.metrics = append(ds.metrics, metrics)

	// Keep only the last maxHistory metrics
	if len(ds.metrics) > ds.maxHistory {
		ds.metrics = ds.metrics[len(ds.metrics)-ds.maxHistory:]
	}
}

// LatencyMeasurer handles latency measurements
type LatencyMeasurer struct {
	rand *time.Time
}

// NewLatencyMeasurer creates a new latency measurer
func NewLatencyMeasurer() *LatencyMeasurer {
	return &LatencyMeasurer{}
}

// MeasureLatency measures latency for different services
func (lm *LatencyMeasurer) MeasureLatency() LatencyInfo {
	// Simulate latency measurements
	now := time.Now()
	seed := now.UnixNano()

	// Simple simulation - in real implementation, you'd measure actual latency
	httpLatency := int64(50 + (seed % 200)) // 50-250ms
	dbLatency := int64(10 + (seed % 40))    // 10-50ms
	apiLatency := int64(30 + (seed % 100))  // 30-130ms

	return LatencyInfo{
		HTTPLatency: httpLatency,
		DBLatency:   dbLatency,
		APILatency:  apiLatency,
	}
}
