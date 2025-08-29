package dashboard

import (
	"encoding/json"
	"net/http"
	"time"

	"system-monitor/internal/alerts"
	"system-monitor/internal/config"
	"system-monitor/internal/datasource"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server handles the web dashboard
type Server struct {
	config     *config.Config
	dataSource datasource.DataSource
	alerts     *alerts.AlertManager
	router     *mux.Router
}

// NewServer creates a new dashboard server
func NewServer(cfg *config.Config, dataSource datasource.DataSource, alertManager *alerts.AlertManager) *Server {
	s := &Server{
		config:     cfg,
		dataSource: dataSource,
		alerts:     alertManager,
		router:     mux.NewRouter(),
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all the routes for the dashboard
func (s *Server) setupRoutes() {
	// Static files
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// API routes
	s.router.HandleFunc("/api/metrics", s.handleGetMetrics).Methods("GET")
	s.router.HandleFunc("/api/metrics/latest", s.handleGetLatestMetrics).Methods("GET")
	s.router.HandleFunc("/api/metrics/history", s.handleGetMetricsHistory).Methods("GET")
	s.router.HandleFunc("/api/alerts/state", s.handleGetAlertState).Methods("GET")
	s.router.HandleFunc("/api/charts/cpu", s.handleGetCPUChart).Methods("GET")
	s.router.HandleFunc("/api/charts/memory", s.handleGetMemoryChart).Methods("GET")
	s.router.HandleFunc("/api/charts/latency", s.handleGetLatencyChart).Methods("GET")
	s.router.HandleFunc("/api/config", s.handleGetConfig).Methods("GET")
	s.router.HandleFunc("/api/config", s.handleUpdateConfig).Methods("PUT")
	s.router.HandleFunc("/api/health", s.handleHealthCheck).Methods("GET")

	// Dashboard page
	s.router.HandleFunc("/", s.handleDashboard).Methods("GET")
}

// Start starts the web server
func (s *Server) Start() error {
	addr := ":" + s.config.DashboardPort
	logrus.Infof("Starting dashboard server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// handleDashboard serves the main dashboard page
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/dashboard.html")
}

// handleGetMetrics returns all metrics
func (s *Server) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metrics, err := s.dataSource.GetMetricsHistory(ctx, time.Now().Add(-1*time.Hour), time.Now())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, metrics)
}

// handleGetLatestMetrics returns the most recent metrics
func (s *Server) handleGetLatestMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	metrics, err := s.dataSource.GetLatestMetrics(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, metrics)
}

// handleGetMetricsHistory returns metrics history with optional time range
func (s *Server) handleGetMetricsHistory(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "Invalid start time format", http.StatusBadRequest)
			return
		}
	} else {
		start = time.Now().Add(-1 * time.Hour) // Default to last hour
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "Invalid end time format", http.StatusBadRequest)
			return
		}
	} else {
		end = time.Now()
	}

	ctx := r.Context()
	metrics, err := s.dataSource.GetMetricsHistory(ctx, start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSON(w, metrics)
}

// handleGetAlertState returns the current alert state
func (s *Server) handleGetAlertState(w http.ResponseWriter, r *http.Request) {
	state := s.alerts.GetAlertState()
	sendJSON(w, state)
}

// handleGetCPUChart returns CPU usage chart data
func (s *Server) handleGetCPUChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cpuData, err := s.dataSource.GetCPUHistory(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamps, err := s.dataSource.GetTimestamps(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chart := generateCPUChart(timestamps, cpuData)
	sendJSON(w, chart)
}

// handleGetMemoryChart returns memory usage chart data
func (s *Server) handleGetMemoryChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	memoryData, err := s.dataSource.GetMemoryHistory(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamps, err := s.dataSource.GetTimestamps(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chart := generateMemoryChart(timestamps, memoryData)
	sendJSON(w, chart)
}

// handleGetLatencyChart returns latency chart data
func (s *Server) handleGetLatencyChart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	latencyData, err := s.dataSource.GetLatencyHistory(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timestamps, err := s.dataSource.GetTimestamps(ctx, 1*time.Hour)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chart := generateLatencyChart(timestamps, latencyData)
	sendJSON(w, chart)
}

// handleGetConfig returns the current configuration
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config := map[string]interface{}{
		"data_source_type":   s.config.DataSourceType,
		"data_source_url":    s.config.DataSourceURL,
		"alert_backend_type": s.config.AlertBackendType,
		"cpu_threshold":      s.config.CPUThreshold,
		"memory_threshold":   s.config.MemoryThreshold,
		"latency_threshold":  s.config.LatencyThreshold,
		"alert_cooldown":     s.config.AlertCooldown.Seconds(),
		"metrics_interval":   s.config.MetricsInterval.Seconds(),
		"dashboard_port":     s.config.DashboardPort,
		"slack_channel":      s.config.SlackChannel,
	}
	sendJSON(w, config)
}

// handleUpdateConfig updates configuration values
func (s *Server) handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update thresholds if provided
	if cpuThreshold, ok := config["cpu_threshold"].(float64); ok {
		s.config.CPUThreshold = cpuThreshold
	}
	if memoryThreshold, ok := config["memory_threshold"].(float64); ok {
		s.config.MemoryThreshold = memoryThreshold
	}
	if latencyThreshold, ok := config["latency_threshold"].(float64); ok {
		s.config.LatencyThreshold = int64(latencyThreshold)
	}

	sendJSON(w, map[string]string{"status": "updated"})
}

// handleHealthCheck returns the health status
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check data source health
	dataSourceHealth := "healthy"
	if err := s.dataSource.HealthCheck(ctx); err != nil {
		dataSourceHealth = "unhealthy"
		logrus.Errorf("Data source health check failed: %v", err)
	}

	health := map[string]interface{}{
		"status":      "healthy",
		"data_source": dataSourceHealth,
		"timestamp":   time.Now().Format(time.RFC3339),
		"uptime":      time.Since(time.Now()).String(), // This would need to be tracked properly
	}

	if dataSourceHealth == "unhealthy" {
		health["status"] = "degraded"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	sendJSON(w, health)
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// ChartData represents chart data structure
type ChartData struct {
	Title   string                 `json:"title"`
	XAxis   []string               `json:"x_axis"`
	Series  []SeriesData           `json:"series"`
	Options map[string]interface{} `json:"options"`
}

// SeriesData represents a data series for charts
type SeriesData struct {
	Name  string        `json:"name"`
	Data  []interface{} `json:"data"`
	Type  string        `json:"type"`
	Color string        `json:"color,omitempty"`
}

// generateCPUChart creates a CPU usage line chart
func generateCPUChart(timestamps []time.Time, cpuData []float64) ChartData {
	// Convert timestamps to readable format
	xAxis := make([]string, len(timestamps))
	for i, t := range timestamps {
		xAxis[i] = t.Format("15:04:05")
	}

	// Convert CPU data to interface slice
	data := make([]interface{}, len(cpuData))
	for i, v := range cpuData {
		data[i] = v
	}

	return ChartData{
		Title: "CPU Usage Over Time",
		XAxis: xAxis,
		Series: []SeriesData{
			{
				Name:  "CPU Usage (%)",
				Data:  data,
				Type:  "line",
				Color: "#5470c6",
			},
		},
		Options: map[string]interface{}{
			"yAxis": map[string]interface{}{
				"min":  0,
				"max":  100,
				"name": "CPU Usage (%)",
			},
		},
	}
}

// generateMemoryChart creates a memory usage area chart
func generateMemoryChart(timestamps []time.Time, memoryData []float64) ChartData {
	// Convert timestamps to readable format
	xAxis := make([]string, len(timestamps))
	for i, t := range timestamps {
		xAxis[i] = t.Format("15:04:05")
	}

	// Convert memory data to interface slice
	data := make([]interface{}, len(memoryData))
	for i, v := range memoryData {
		data[i] = v
	}

	return ChartData{
		Title: "Memory Usage Over Time",
		XAxis: xAxis,
		Series: []SeriesData{
			{
				Name:  "Memory Usage (%)",
				Data:  data,
				Type:  "line",
				Color: "#91cc75",
			},
		},
		Options: map[string]interface{}{
			"yAxis": map[string]interface{}{
				"min":  0,
				"max":  100,
				"name": "Memory Usage (%)",
			},
		},
	}
}

// generateLatencyChart creates a latency line chart
func generateLatencyChart(timestamps []time.Time, latencyData []int64) ChartData {
	// Convert timestamps to readable format
	xAxis := make([]string, len(timestamps))
	for i, t := range timestamps {
		xAxis[i] = t.Format("15:04:05")
	}

	// Convert latency data to interface slice
	data := make([]interface{}, len(latencyData))
	for i, v := range latencyData {
		data[i] = v
	}

	return ChartData{
		Title: "HTTP Latency Over Time",
		XAxis: xAxis,
		Series: []SeriesData{
			{
				Name:  "HTTP Latency (ms)",
				Data:  data,
				Type:  "line",
				Color: "#ee6666",
			},
		},
		Options: map[string]interface{}{
			"yAxis": map[string]interface{}{
				"min":  0,
				"name": "Latency (ms)",
			},
		},
	}
}
