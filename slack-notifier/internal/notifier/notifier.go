package notifier

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"slack-notifier/internal/config"
	"slack-notifier/internal/events"
	slackpkg "slack-notifier/pkg/slack"
)

// Service coordinates event processing and Slack posting.
type Service struct {
	cfg         *config.Config
	slack       *slackpkg.Client
	events      chan *events.Event
	workers     int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewNotifierService(cfg *config.Config, workers int) (*Service, error) {
	s := &Service{
		cfg:     cfg,
		events:  make(chan *events.Event, 128),
		workers: workers,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.slack = slackpkg.NewClient(cfg.SlackBotToken, cfg.SlackChannel)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.slack.TestConnection(ctx); err != nil {
		return nil, fmt.Errorf("slack auth failed: %w", err)
	}
	return s, nil
}

func (s *Service) Start() error {
	log.Printf("notifier: starting with %d workers", s.workers)
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
	// send startup event
	s.SendEvent(events.NewEvent(events.EventTypeSystemStartup).
		WithTitle("Notifier Started").
		WithMessage("Slack notifier is up").
		WithSeverity(events.SeverityInfo).
		Build())
	return nil
}

func (s *Service) Stop() error {
	log.Println("notifier: stopping")
	// send shutdown event (best-effort)
	s.SendEvent(events.NewEvent(events.EventTypeSystemShutdown).
		WithTitle("Notifier Stopping").
		WithMessage("Slack notifier is shutting down").
		WithSeverity(events.SeverityInfo).
		Build())

	s.cancel()
	close(s.events)
	s.wg.Wait()
	return nil
}

func (s *Service) SendEvent(e *events.Event) {
	if e.Channel == "" {
		e.Channel = s.cfg.SlackChannel
	}
	select {
	case s.events <- e:
		// queued
	case <-s.ctx.Done():
		log.Printf("notifier: drop event %s (%s): stopping", e.ID, e.Type)
	}
}

func (s *Service) SendMessage(msg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return s.slack.SendMessage(ctx, msg)
}

func (s *Service) worker(id int) {
	defer s.wg.Done()
	for {
		select {
		case e, ok := <-s.events:
			if !ok {
				return
			}
			s.process(e)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Service) process(e *events.Event) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.slack.SendEvent(ctx, e); err != nil {
		log.Printf("notifier: send failed: %v", err)
		// emit error event without requeue to avoid loops
		errEvent := events.NewEvent(events.EventTypeSystemError).
			WithTitle("Notification Failed").
			WithMessage(fmt.Sprintf("event %s failed: %v", e.ID, err)).
			WithSeverity(events.SeverityError).
			WithMetadata("original_event_id", e.ID).
			WithMetadata("original_event_type", e.Type).
			Build()
		_ = s.slack.SendEvent(ctx, errEvent)
	}
}

func (s *Service) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"queue_size": len(s.events),
		"workers":   s.workers,
	}
}
