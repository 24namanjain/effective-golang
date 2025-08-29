# System Monitor - How It Works

A simple Go application that monitors your system (CPU, Memory, Latency) and sends Slack alerts when thresholds are exceeded.

## 🎯 What This App Does

1. **Collects Metrics**: Every 5 seconds, it checks your system's CPU usage, memory usage, and HTTP latency
2. **Checks Thresholds**: Compares the metrics against your configured thresholds (e.g., CPU > 10%)
3. **Sends Alerts**: If thresholds are exceeded, it sends a message to your Slack channel
4. **Shows Dashboard**: Provides a web interface at `http://localhost:8080` to view current metrics

## 🔄 How It Works (Step by Step)

### 1. Application Startup
```
┌─────────────────────────────────────────────────────────────┐
│                    Application Starts                       │
├─────────────────────────────────────────────────────────────┤
│ 1. Loads configuration from environment variables          │
│ 2. Creates data source (collects system metrics)           │
│ 3. Creates alert backend (sends Slack messages)            │
│ 4. Starts web dashboard server                             │
│ 5. Sends startup notification to Slack                     │
└─────────────────────────────────────────────────────────────┘
```

### 2. Continuous Monitoring Loop
```
┌─────────────────────────────────────────────────────────────┐
│                    Every 5 Seconds                         │
├─────────────────────────────────────────────────────────────┤
│ 1. Collect current metrics:                                │
│    • CPU usage (e.g., 28.5%)                              │
│    • Memory usage (e.g., 65.2%)                           │
│    • HTTP latency (e.g., 45ms)                            │
│                                                           │
│ 2. Check each metric against thresholds:                   │
│    • CPU > 10%? → Send Slack alert                        │
│    • Memory > 80%? → Send Slack alert                     │
│    • Latency > 100ms? → Send Slack alert                  │
│                                                           │
│ 3. Update web dashboard with new data                     │
└─────────────────────────────────────────────────────────────┘
```

### 3. Alert Logic
```
┌─────────────────────────────────────────────────────────────┐
│                    Alert Processing                        │
├─────────────────────────────────────────────────────────────┤
│ IF metric > threshold:                                     │
│   • Create alert message                                   │
│   • Determine severity (warning/critical)                  │
│   • Send to Slack channel                                  │
│   • Log alert sent                                         │
│   • Wait 5 minutes before sending another alert           │
│                                                           │
│ ELSE:                                                      │
│   • Reset alert state (back to normal)                    │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### 1. Set Up Environment Variables
```bash
# Required for Slack alerts
export SLACK_BOT_TOKEN="xoxb-your-bot-token"
export SLACK_CHANNEL="#alerts"

# Optional: Customize thresholds (defaults shown)
export CPU_THRESHOLD="10.0"        # Alert if CPU > 10%
export MEMORY_THRESHOLD="80.0"     # Alert if Memory > 80%
export LATENCY_THRESHOLD="100"     # Alert if Latency > 100ms
```

### 2. Build and Run
```bash
# Build the application
go build -o bin/system-monitor cmd/monitor/main.go

# Run with Slack alerts
SLACK_BOT_TOKEN="xoxb-your-token" SLACK_CHANNEL="#alerts" ./bin/system-monitor
```

### 3. View Dashboard
Open your browser to: `http://localhost:8080`

## 📊 What You'll See

### In Slack Channel
```
🚨 High CPU Usage Alert
CPU usage is 28.5% (threshold: 10.0%)
```

### In Web Dashboard
- Real-time CPU, Memory, and Latency graphs
- Current alert status (warning/critical)
- System health indicators
- Configuration settings

### In Application Logs
```
INFO[2025-08-29T15:37:48+05:30] ✅ System Monitor is running!
INFO[2025-08-29T15:37:48+05:30] 📊 Dashboard available at: http://localhost:8080
INFO[2025-08-29T15:37:58+05:30] 🚨 CPU Alert sent: 28.5% usage (threshold: 10.0%)
```

## ⚙️ Configuration Options

### Thresholds
- `CPU_THRESHOLD`: CPU usage percentage (default: 80%)
- `MEMORY_THRESHOLD`: Memory usage percentage (default: 80%)
- `LATENCY_THRESHOLD`: HTTP latency in milliseconds (default: 100ms)

### Timing
- `METRICS_INTERVAL`: How often to check metrics (default: 5s)
- `ALERT_COOLDOWN`: Wait time between alerts (default: 5m)

### Application
- `DASHBOARD_PORT`: Web dashboard port (default: 8080)
- `ENVIRONMENT`: Environment name (default: development)

## 🏗️ Code Structure

```
grafana-dashboard/
├── cmd/monitor/main.go          # Main application entry point
├── internal/
│   ├── config/config.go         # Configuration loading
│   ├── datasource/
│   │   ├── interface.go         # Data source interface
│   │   ├── local.go            # Collects local system metrics
│   │   └── factory.go          # Creates data sources
│   ├── alerts/
│   │   ├── interface.go        # Alert backend interface
│   │   ├── slack.go           # Sends Slack messages
│   │   └── factory.go         # Creates alert backends
│   └── dashboard/server.go     # Web dashboard server
└── web/                        # Dashboard HTML/CSS/JS files
```

## 🔧 Key Components Explained

### 1. Data Source (`internal/datasource/local.go`)
- Uses `gopsutil` library to collect system metrics
- Runs every 5 seconds to get current CPU, memory, and latency
- Simulates HTTP latency by measuring localhost response time

### 2. Alert Manager (`internal/alerts/interface.go`)
- Compares metrics against configured thresholds
- Implements cooldown logic (prevents spam)
- Determines alert severity (warning vs critical)
- Manages alert state (active/inactive)

### 3. Slack Backend (`internal/alerts/slack.go`)
- Formats alert messages with emojis and details
- Sends messages to configured Slack channel
- Handles Slack API authentication and errors

### 4. Dashboard Server (`internal/dashboard/server.go`)
- Serves web interface at `http://localhost:8080`
- Provides API endpoints for metrics data
- Updates charts in real-time via JavaScript

## 🐛 Troubleshooting

### No Slack Alerts?
1. Check your `SLACK_BOT_TOKEN` is correct
2. Ensure the bot is added to your `SLACK_CHANNEL`
3. Verify the channel name starts with `#`

### Dashboard Not Loading?
1. Check if port 8080 is available
2. Look for error messages in the application logs
3. Try accessing `http://localhost:8080/api/health`

### No Metrics Data?
1. Check if the application has permission to read system metrics
2. Look for errors in the application logs
3. Verify the data source is working

## 🎯 Example Usage

```bash
# Run with custom CPU threshold (10%)
SLACK_BOT_TOKEN="xoxb-your-token" \
SLACK_CHANNEL="#alerts" \
CPU_THRESHOLD=10 \
./bin/system-monitor

# Run with shorter intervals for testing
SLACK_BOT_TOKEN="xoxb-your-token" \
SLACK_CHANNEL="#alerts" \
METRICS_INTERVAL=2s \
ALERT_COOLDOWN=1m \
./bin/system-monitor
```

That's it! The application is designed to be simple and focused on one thing: monitoring your system and alerting you when something goes wrong.
