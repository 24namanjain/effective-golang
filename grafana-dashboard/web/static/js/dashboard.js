class SystemMonitorDashboard {
    constructor() {
        this.charts = {};
        this.updateInterval = 5000; // 5 seconds
        this.config = {};
        this.alertState = {};
        this.init();
    }

    async init() {
        console.log('🚀 Initializing System Monitor Dashboard');
        
        // Initialize charts
        this.initCharts();
        
        // Load initial configuration
        await this.loadConfig();
        
        // Load initial data
        await this.loadLatestMetrics();
        await this.loadAlertState();
        await this.loadSystemInfo();
        
        // Set up event listeners
        this.setupEventListeners();
        
        // Start real-time updates
        this.startRealTimeUpdates();
        
        console.log('✅ Dashboard initialized successfully');
    }

    initCharts() {
        // Initialize CPU chart
        this.charts.cpu = echarts.init(document.getElementById('cpu-chart'));
        this.charts.memory = echarts.init(document.getElementById('memory-chart'));
        this.charts.latency = echarts.init(document.getElementById('latency-chart'));

        // Set up basic chart options
        const basicOption = {
            grid: {
                left: '10%',
                right: '10%',
                top: '10%',
                bottom: '15%'
            },
            tooltip: {
                trigger: 'axis',
                formatter: function(params) {
                    const data = params[0];
                    return `${data.name}<br/>${data.seriesName}: ${data.value}`;
                }
            },
            xAxis: {
                type: 'category',
                data: [],
                axisLabel: {
                    rotate: 45,
                    fontSize: 10
                }
            },
            yAxis: {
                type: 'value',
                axisLabel: {
                    fontSize: 10
                }
            },
            series: [{
                data: [],
                type: 'line',
                smooth: true,
                lineStyle: {
                    width: 2
                },
                areaStyle: {
                    opacity: 0.1
                }
            }]
        };

        // Apply basic options to all charts
        this.charts.cpu.setOption(basicOption);
        this.charts.memory.setOption(basicOption);
        this.charts.latency.setOption(basicOption);

        // Handle window resize
        window.addEventListener('resize', () => {
            Object.values(this.charts).forEach(chart => chart.resize());
        });
    }

    async loadConfig() {
        try {
            const response = await fetch('/api/config');
            if (!response.ok) throw new Error('Failed to load config');
            
            this.config = await response.json();
            this.updateConfigUI();
        } catch (error) {
            console.error('Error loading config:', error);
        }
    }

    updateConfigUI() {
        document.getElementById('cpu-threshold').value = this.config.cpu_threshold || 80;
        document.getElementById('memory-threshold').value = this.config.memory_threshold || 85;
        document.getElementById('latency-threshold').value = this.config.latency_threshold || 500;
        
        document.getElementById('data-source-type').textContent = this.config.data_source_type || '--';
        document.getElementById('alert-backend-type').textContent = this.config.alert_backend_type || '--';
    }

    async loadLatestMetrics() {
        try {
            const response = await fetch('/api/metrics/latest');
            if (!response.ok) throw new Error('Failed to load metrics');
            
            const metrics = await response.json();
            if (metrics) {
                this.updateMetricsDisplay(metrics);
            }
        } catch (error) {
            console.error('Error loading metrics:', error);
        }
    }

    async loadChartData() {
        try {
            // Load CPU chart data
            const cpuResponse = await fetch('/api/charts/cpu');
            if (cpuResponse.ok) {
                const cpuData = await cpuResponse.json();
                this.updateChart('cpu', cpuData);
            }

            // Load Memory chart data
            const memoryResponse = await fetch('/api/charts/memory');
            if (memoryResponse.ok) {
                const memoryData = await memoryResponse.json();
                this.updateChart('memory', memoryData);
            }

            // Load Latency chart data
            const latencyResponse = await fetch('/api/charts/latency');
            if (latencyResponse.ok) {
                const latencyData = await latencyResponse.json();
                this.updateChart('latency', latencyData);
            }
        } catch (error) {
            console.error('Error loading chart data:', error);
        }
    }

    updateChart(chartType, data) {
        const chart = this.charts[chartType];
        if (!chart || !data) return;

        const option = {
            xAxis: {
                data: data.x_axis || []
            },
            series: [{
                data: data.series?.[0]?.data || [],
                name: data.series?.[0]?.name || chartType.toUpperCase(),
                color: data.series?.[0]?.color || '#667eea'
            }]
        };

        chart.setOption(option);
    }

    updateMetricsDisplay(metrics) {
        // Update current values
        document.getElementById('cpu-value').textContent = `${metrics.cpu?.toFixed(1)}%`;
        document.getElementById('memory-value').textContent = `${metrics.memory?.percent?.toFixed(1)}%`;
        document.getElementById('latency-value').textContent = `${metrics.latency?.http_latency}ms`;

        // Update last update time
        document.getElementById('last-update').textContent = new Date().toLocaleTimeString();

        // Update status indicators based on thresholds
        this.updateStatusIndicators(metrics);
    }

    updateStatusIndicators(metrics) {
        const cpuThreshold = parseFloat(document.getElementById('cpu-threshold').value);
        const memoryThreshold = parseFloat(document.getElementById('memory-threshold').value);
        const latencyThreshold = parseInt(document.getElementById('latency-threshold').value);

        // Update data source status
        const dataSourceStatus = document.querySelector('#data-source-status .status-dot');
        dataSourceStatus.className = 'status-dot healthy';

        // Update alert status based on metrics
        const alertStatus = document.querySelector('#alert-status .status-dot');
        let hasAlerts = false;

        if (metrics.cpu > cpuThreshold) hasAlerts = true;
        if (metrics.memory?.percent > memoryThreshold) hasAlerts = true;
        if (metrics.latency?.http_latency > latencyThreshold) hasAlerts = true;

        alertStatus.className = `status-dot ${hasAlerts ? 'warning' : 'healthy'}`;
    }

    async loadAlertState() {
        try {
            const response = await fetch('/api/alerts/state');
            if (!response.ok) throw new Error('Failed to load alert state');
            
            this.alertState = await response.json();
            this.updateAlertDisplay();
        } catch (error) {
            console.error('Error loading alert state:', error);
        }
    }

    updateAlertDisplay() {
        const alertsGrid = document.getElementById('alerts-grid');
        alertsGrid.innerHTML = '';

        const alerts = [
            { name: 'CPU', warning: this.alertState.cpu_warning, critical: this.alertState.cpu_critical },
            { name: 'Memory', warning: this.alertState.memory_warning, critical: this.alertState.memory_critical },
            { name: 'Latency', warning: this.alertState.latency_warning, critical: this.alertState.latency_critical }
        ];

        alerts.forEach(alert => {
            const alertItem = document.createElement('div');
            alertItem.className = 'alert-item';
            
            if (alert.critical) {
                alertItem.classList.add('critical');
                alertItem.innerHTML = `🚨 ${alert.name} Critical`;
            } else if (alert.warning) {
                alertItem.classList.add('warning');
                alertItem.innerHTML = `⚠️ ${alert.name} Warning`;
            } else {
                alertItem.classList.add('success');
                alertItem.innerHTML = `✅ ${alert.name} Normal`;
            }
            
            alertsGrid.appendChild(alertItem);
        });
    }

    async loadSystemInfo() {
        try {
            const response = await fetch('/api/health');
            if (!response.ok) throw new Error('Failed to load system info');
            
            const health = await response.json();
            document.getElementById('uptime').textContent = health.uptime || '--';
        } catch (error) {
            console.error('Error loading system info:', error);
        }
    }

    setupEventListeners() {
        // Update configuration button
        document.getElementById('update-config').addEventListener('click', async () => {
            await this.updateConfiguration();
        });

        // Threshold input changes
        ['cpu-threshold', 'memory-threshold', 'latency-threshold'].forEach(id => {
            document.getElementById(id).addEventListener('change', () => {
                this.updateStatusIndicators(this.lastMetrics);
            });
        });
    }

    async updateConfiguration() {
        const config = {
            cpu_threshold: parseFloat(document.getElementById('cpu-threshold').value),
            memory_threshold: parseFloat(document.getElementById('memory-threshold').value),
            latency_threshold: parseInt(document.getElementById('latency-threshold').value)
        };

        try {
            const response = await fetch('/api/config', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(config)
            });

            if (response.ok) {
                console.log('Configuration updated successfully');
                this.showNotification('Configuration updated successfully', 'success');
            } else {
                throw new Error('Failed to update configuration');
            }
        } catch (error) {
            console.error('Error updating configuration:', error);
            this.showNotification('Failed to update configuration', 'error');
        }
    }

    startRealTimeUpdates() {
        setInterval(async () => {
            await this.loadLatestMetrics();
            await this.loadAlertState();
            await this.loadSystemInfo();
        }, this.updateInterval);

        // Update charts less frequently
        setInterval(async () => {
            await this.loadChartData();
        }, this.updateInterval * 2);
    }

    showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 1000;
            animation: slideIn 0.3s ease-out;
            background: ${type === 'success' ? '#48bb78' : type === 'error' ? '#f56565' : '#3182ce'};
        `;

        document.body.appendChild(notification);

        // Remove notification after 3 seconds
        setTimeout(() => {
            notification.style.animation = 'slideOut 0.3s ease-out';
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
            }, 300);
        }, 3000);
    }
}

// Add CSS animations for notifications
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new SystemMonitorDashboard();
});
