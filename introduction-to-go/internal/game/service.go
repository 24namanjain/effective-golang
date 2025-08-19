package game

import (
	"context"
	"fmt"
	"sync"
	"time"

	"effective-golang/internal/models"
)

// GameService handles game logic and concurrent operations
type GameService struct {
	gameRepo        models.GameRepository
	userRepo        models.UserRepository
	leaderboardRepo models.LeaderboardRepository
	cacheRepo       models.CacheRepository
	
	// Worker pool for processing game events
	eventWorkers    chan struct{}
	eventQueue      chan *GameEvent
	eventProcessor  *EventProcessor
	
	// Game state management
	activeGames     map[string]*models.Game
	gameMutex       sync.RWMutex
	
	// Configuration
	maxWorkers      int
	queueSize       int
}

// GameEvent represents a game event to be processed
type GameEvent struct {
	GameID    string
	PlayerID  string
	EventType string
	Score     int64
	Data      interface{}
	Timestamp time.Time
}

// GameResult represents the result of a completed game
type GameResult struct {
	GameID     string
	WinnerID   string
	LoserID    string
	WinnerScore int64
	LoserScore  int64
	Duration   time.Duration
	IsTie      bool
}

// EventProcessor handles game event processing
type EventProcessor struct {
	workers    chan struct{}
	queue      chan *GameEvent
	gameSvc    *GameService
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

// Custom errors for game operations
var (
	ErrGameNotFound     = fmt.Errorf("game not found")
	ErrGameAlreadyEnded = fmt.Errorf("game already ended")
	ErrInvalidPlayer    = fmt.Errorf("invalid player")
	ErrGameNotStarted   = fmt.Errorf("game not started")
	ErrEventQueueFull   = fmt.Errorf("event queue is full")
)

// NewGameService creates a new game service
func NewGameService(
	gameRepo models.GameRepository,
	userRepo models.UserRepository,
	leaderboardRepo models.LeaderboardRepository,
	cacheRepo models.CacheRepository,
	maxWorkers, queueSize int,
) *GameService {
	ctx, cancel := context.WithCancel(context.Background())
	
	svc := &GameService{
		gameRepo:        gameRepo,
		userRepo:        userRepo,
		leaderboardRepo: leaderboardRepo,
		cacheRepo:       cacheRepo,
		activeGames:     make(map[string]*models.Game),
		maxWorkers:      maxWorkers,
		queueSize:       queueSize,
	}
	
	// Initialize event processor
	svc.eventWorkers = make(chan struct{}, maxWorkers)
	svc.eventQueue = make(chan *GameEvent, queueSize)
	svc.eventProcessor = &EventProcessor{
		workers: svc.eventWorkers,
		queue:   svc.eventQueue,
		gameSvc: svc,
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// Start event processor
	svc.eventProcessor.Start()
	
	return svc
}

// CreateGame creates a new game between two players
func (s *GameService) CreateGame(ctx context.Context, player1ID, player2ID string) (*models.Game, error) {
	// Validate players exist
	_, err := s.userRepo.GetByID(ctx, player1ID)
	if err != nil {
		return nil, fmt.Errorf("player1 not found: %w", err)
	}
	
	_, err = s.userRepo.GetByID(ctx, player2ID)
	if err != nil {
		return nil, fmt.Errorf("player2 not found: %w", err)
	}
	
	// Create game
	game, err := models.NewGame(player1ID, player2ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}
	
	// Save to database
	if err := s.gameRepo.Create(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to save game: %w", err)
	}
	
	// Add to active games
	s.gameMutex.Lock()
	s.activeGames[game.ID] = game
	s.gameMutex.Unlock()
	
	return game, nil
}

// StartGame starts a game
func (s *GameService) StartGame(ctx context.Context, gameID string) error {
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}
	
	if err := game.Start(); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}
	
	// Update in database
	if err := s.gameRepo.Update(ctx, game); err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}
	
	// Process game start event
	s.QueueEvent(&GameEvent{
		GameID:    gameID,
		EventType: "game_started",
		Timestamp: time.Now(),
	})
	
	return nil
}

// UpdateScore updates a player's score in a game
func (s *GameService) UpdateScore(ctx context.Context, gameID, playerID string, score int64) error {
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return fmt.Errorf("failed to get game: %w", err)
	}
	
	if err := game.UpdateScore(playerID, score); err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}
	
	// Update in database
	if err := s.gameRepo.Update(ctx, game); err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}
	
	// Queue score update event
	s.QueueEvent(&GameEvent{
		GameID:    gameID,
		PlayerID:  playerID,
		EventType: "score_updated",
		Score:     score,
		Timestamp: time.Now(),
	})
	
	return nil
}

// EndGame ends a game and processes results
func (s *GameService) EndGame(ctx context.Context, gameID string) (*GameResult, error) {
	game, err := s.getGame(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	
	if err := game.End(); err != nil {
		return nil, fmt.Errorf("failed to end game: %w", err)
	}
	
	// Update in database
	if err := s.gameRepo.Update(ctx, game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}
	
	// Remove from active games
	s.gameMutex.Lock()
	delete(s.activeGames, gameID)
	s.gameMutex.Unlock()
	
	// Create game result
	result := &GameResult{
		GameID:      gameID,
		WinnerScore: game.Score1,
		LoserScore:  game.Score2,
		Duration:    game.GetDuration(),
		IsTie:       game.GetWinner() == "",
	}
	
	if !result.IsTie {
		result.WinnerID = game.GetWinner()
		if result.WinnerID == game.Player1ID {
			result.LoserID = game.Player2ID
		} else {
			result.LoserID = game.Player1ID
		}
	}
	
	// Queue game end event
	s.QueueEvent(&GameEvent{
		GameID:    gameID,
		EventType: "game_ended",
		Data:      result,
		Timestamp: time.Now(),
	})
	
	return result, nil
}

// GetActiveGames returns all active games
func (s *GameService) GetActiveGames(ctx context.Context) ([]*models.Game, error) {
	s.gameMutex.RLock()
	defer s.gameMutex.RUnlock()
	
	games := make([]*models.Game, 0, len(s.activeGames))
	for _, game := range s.activeGames {
		games = append(games, game)
	}
	
	return games, nil
}

// GetGame returns a game by ID
func (s *GameService) GetGame(ctx context.Context, gameID string) (*models.Game, error) {
	return s.getGame(ctx, gameID)
}

// QueueEvent queues a game event for processing
func (s *GameService) QueueEvent(event *GameEvent) error {
	select {
	case s.eventQueue <- event:
		return nil
	default:
		return ErrEventQueueFull
	}
}

// getGame retrieves a game from cache or database
func (s *GameService) getGame(ctx context.Context, gameID string) (*models.Game, error) {
	// Try to get from active games first
	s.gameMutex.RLock()
	if game, exists := s.activeGames[gameID]; exists {
		s.gameMutex.RUnlock()
		return game, nil
	}
	s.gameMutex.RUnlock()
	
	// Get from database
	game, err := s.gameRepo.GetByID(ctx, gameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}
	
	// Add to active games if it's still active
	if game.State == models.GameStatePlaying {
		s.gameMutex.Lock()
		s.activeGames[gameID] = game
		s.gameMutex.Unlock()
	}
	
	return game, nil
}

// Close shuts down the game service
func (s *GameService) Close() error {
	s.eventProcessor.Stop()
	return nil
}

// EventProcessor methods

// Start starts the event processor
func (ep *EventProcessor) Start() {
	for i := 0; i < cap(ep.workers); i++ {
		ep.workers <- struct{}{}
	}
	
	ep.wg.Add(1)
	go ep.processEvents()
}

// Stop stops the event processor
func (ep *EventProcessor) Stop() {
	ep.cancel()
	ep.wg.Wait()
}

// processEvents processes events from the queue
func (ep *EventProcessor) processEvents() {
	defer ep.wg.Done()
	
	for {
		select {
		case event := <-ep.queue:
			ep.wg.Add(1)
			go ep.processEvent(event)
		case <-ep.ctx.Done():
			return
		}
	}
}

// processEvent processes a single event
func (ep *EventProcessor) processEvent(event *GameEvent) {
	defer ep.wg.Done()
	
	// Acquire worker
	<-ep.workers
	defer func() { ep.workers <- struct{}{} }()
	
	ctx := context.Background()
	
	switch event.EventType {
	case "game_started":
		ep.handleGameStarted(ctx, event)
	case "score_updated":
		ep.handleScoreUpdated(ctx, event)
	case "game_ended":
		ep.handleGameEnded(ctx, event)
	default:
		// Log unknown event type
		fmt.Printf("Unknown event type: %s\n", event.EventType)
	}
}

// handleGameStarted handles game start events
func (ep *EventProcessor) handleGameStarted(ctx context.Context, event *GameEvent) {
	// Cache game state
	cacheKey := fmt.Sprintf("game:%s", event.GameID)
	ep.gameSvc.cacheRepo.Set(ctx, cacheKey, event, 3600)
}

// handleScoreUpdated handles score update events
func (ep *EventProcessor) handleScoreUpdated(ctx context.Context, event *GameEvent) {
	// Update cached game state
	cacheKey := fmt.Sprintf("game:%s", event.GameID)
	ep.gameSvc.cacheRepo.Set(ctx, cacheKey, event, 3600)
}

// handleGameEnded handles game end events
func (ep *EventProcessor) handleGameEnded(ctx context.Context, event *GameEvent) {
	result, ok := event.Data.(*GameResult)
	if !ok {
		return
	}
	
	// Update user statistics
	ep.updateUserStats(ctx, result)
	
	// Update leaderboards
	ep.updateLeaderboards(ctx, result)
	
	// Clean up cached game state
	cacheKey := fmt.Sprintf("game:%s", event.GameID)
	ep.gameSvc.cacheRepo.Delete(ctx, cacheKey)
}

// updateUserStats updates user statistics after a game
func (ep *EventProcessor) updateUserStats(ctx context.Context, result *GameResult) {
	// Update winner stats
	if result.WinnerID != "" {
		stats, err := ep.gameSvc.userRepo.GetStats(ctx, result.WinnerID)
		if err == nil {
			stats.UpdateStats(result.WinnerScore, true)
			ep.gameSvc.userRepo.UpdateStats(ctx, stats)
		}
	}
	
	// Update loser stats
	if result.LoserID != "" {
		stats, err := ep.gameSvc.userRepo.GetStats(ctx, result.LoserID)
		if err == nil {
			stats.UpdateStats(result.LoserScore, false)
			ep.gameSvc.userRepo.UpdateStats(ctx, stats)
		}
	}
}

// updateLeaderboards updates leaderboards after a game
func (ep *EventProcessor) updateLeaderboards(ctx context.Context, result *GameResult) {
	// Update global leaderboard
	globalLB, err := ep.gameSvc.leaderboardRepo.GetByName(ctx, "global")
	if err == nil {
		// Update winner
		if result.WinnerID != "" {
			user, _ := ep.gameSvc.userRepo.GetByID(ctx, result.WinnerID)
			if user != nil {
				ep.gameSvc.leaderboardRepo.AddEntry(ctx, globalLB.ID, &models.LeaderboardEntry{
					UserID:    result.WinnerID,
					Username:  user.Username,
					Score:     result.WinnerScore,
					UpdatedAt: time.Now(),
				})
			}
		}
	}
}
