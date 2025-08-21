/*
# Monitoring a Go Project with Prometheus, Grafana, and Loki

This documentation provides a step-by-step guide to set up monitoring and observability for your Go project using Prometheus (metrics), Grafana (visualization), and Loki (logs).

---

## 1. Instrumenting Your Go Application

### a. Metrics with Prometheus

1. **Add Prometheus Client to Your Project**

   ```bash
   go get github.com/prometheus/client_golang/prometheus
   go get github.com/prometheus/client_golang/prometheus/promhttp
   ```

2. **Expose Metrics Endpoint**

   In your `main.go`:

   ```go
   import (
       "net/http"
       "github.com/prometheus/client_golang/prometheus/promhttp"
   )

   func main() {
       // ... your app setup

       // Expose metrics at /metrics
       http.Handle("/metrics", promhttp.Handler())
       http.ListenAndServe(":8080", nil)
   }
   ```

3. **Add Custom Metrics (Optional)**

   ```go
   var (
       opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
           Name: "myapp_processed_ops_total",
           Help: "The total number of processed events",
       })
   )

   func init() {
       prometheus.MustRegister(opsProcessed)
   }
   ```

---

### b. Logging with Loki

1. **Structured Logging**

   Use a structured logger (e.g., [logrus](https://github.com/sirupsen/logrus) or [zap](https://github.com/uber-go/zap)):

   ```bash
   go get github.com/sirupsen/logrus
   ```

   ```go
   import log "github.com/sirupsen/logrus"

   func main() {
       log.SetFormatter(&log.JSONFormatter{})
       log.Info("Application started")
   }
   ```

2. **Shipping Logs to Loki**

   - Use [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/), a log collector that ships logs to Loki.
   - Configure Promtail to watch your application's log files or stdout.

---

## 2. Deploying Prometheus, Grafana, and Loki

- **Docker Compose Example:**

  ```yaml
  version: '3'

  services:
    prometheus:
      image: prom/prometheus
      volumes:
        - ./prometheus.yml:/etc/prometheus/prometheus.yml
      ports:
        - "9090:9090"

    grafana:
      image: grafana/grafana
      ports:
        - "3000:3000"
      environment:
        - GF_SECURITY_ADMIN_PASSWORD=admin

    loki:
      image: grafana/loki
      ports:
        - "3100:3100"
      command: -config.file=/etc/loki/local-config.yaml

    promtail:
      image: grafana/promtail
      volumes:
        - ./promtail-config.yaml:/etc/promtail/config.yaml
        - /var/log:/var/log
      command: -config.file=/etc/promtail/config.yaml
  ```

- **Configure Prometheus to scrape your Go app:**

  ```yaml
  # prometheus.yml
  scrape_configs:
    - job_name: 'go-app'
      static_configs:
        - targets: ['host.docker.internal:8080']
  ```

---

## 3. Visualizing in Grafana

- **Add Prometheus as a Data Source** in Grafana (`http://localhost:9090`).
- **Add Loki as a Data Source** (`http://localhost:3100`).
- **Import Dashboards** for Go metrics and logs, or create your own.

---

## 4. References

- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Grafana Loki](https://grafana.com/oss/loki/)
- [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)

*/
