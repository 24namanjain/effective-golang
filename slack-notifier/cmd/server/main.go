package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"slack-notifier/internal/config"
	"slack-notifier/internal/events"
	"slack-notifier/internal/notifier"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	svc, err := notifier.NewNotifierService(cfg, 2)
	if err != nil {
		log.Fatalf("notifier: %v", err)
	}
	if err := svc.Start(); err != nil {
		log.Fatalf("start: %v", err)
	}

	// simple smoke send
	if err := svc.SendMessage("🚀 Slack Notifier is running"); err != nil {
		log.Printf("send message: %v", err)
	}

	// demo events once on boot
	demoEvents(svc)

	// HTTP API for testing
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/send-message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body.Message == "" {
			body.Message = "Hello from API"
		}

		if err := svc.SendMessage(body.Message); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sent"))
	})

	mux.HandleFunc("/send-event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Type     string                 `json:"type"`
			Title    string                 `json:"title"`
			Message  string                 `json:"message"`
			Severity string                 `json:"severity"`
			Channel  string                 `json:"channel"`
			Metadata map[string]interface{} `json:"metadata"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		sev := events.SeverityInfo
		switch body.Severity {
		case "warning":
			sev = events.SeverityWarning
		case "error":
			sev = events.SeverityError
		case "critical":
			sev = events.SeverityCritical
		}
		evt := events.NewEvent(events.EventType(body.Type)).
			WithTitle(body.Title).
			WithMessage(body.Message).
			WithSeverity(sev).
			WithChannel(body.Channel).
			Build()
		for k, v := range body.Metadata {
			evt.Metadata[k] = v
		}
		svc.SendEvent(evt)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("queued"))
	})

	server := &http.Server{Addr: cfg.APIAddress, Handler: mux}
	go func() {
		log.Printf("api listening on %s", cfg.APIAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("api: %v", err)
		}
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
	if err := svc.Stop(); err != nil {
		log.Printf("stop: %v", err)
	}
}

func demoEvents(svc *notifier.Service) {
	// space the demo a bit
	time.Sleep(500 * time.Millisecond)

	login := events.NewEvent(events.EventTypeUserLogin).
		WithTitle("User Login").
		WithMessage("User john@example.com logged in").
		WithSeverity(events.SeverityInfo).
		WithUserID("john@example.com").
		WithMetadata("ip", "203.0.113.10").
		Build()
	svc.SendEvent(login)

	order := events.NewEvent(events.EventTypeOrderCreated).
		WithTitle("Order Created").
		WithMessage("Order ORD-1001 created").
		WithSeverity(events.SeverityInfo).
		WithMetadata("order_id", "ORD-1001").
		WithMetadata("amount", 299.99).
		Build()
	svc.SendEvent(order)

	warn := events.NewEvent(events.EventTypeHighCPUUsage).
		WithTitle("High CPU Usage").
		WithMessage("CPU > 85% for 5m").
		WithSeverity(events.SeverityWarning).
		WithMetadata("cpu", "89%").
		Build()
	svc.SendEvent(warn)
}
