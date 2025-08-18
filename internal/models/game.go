package models

import (
	"errors"
	"sync"
	"time"
)

// GameState represents the current state of a game
type GameState string

const (
	GameStateWaiting   GameState = "waiting"
	GameStatePlaying   GameState = "playing"
	GameStateFinished  GameState = "finished"
	GameStateCancelled GameState = "cancelled"
)

// Game represents a game session
type Game struct {
	ID          string    `json:"id" db:"id"`
	Player1ID   string    `json:"player1_id" db:"player1_id"`
	Player2ID   string    `json:"player2_id" db:"player2_id"`
	State       GameState `json:"state" db:"state"`
	Score1      int64     `json:"score1" db:"score1"`
	Score2      int64     `json:"score2" db:"score2"`
	WinnerID    *string   `json:"winner_id,omitempty" db:"winner_id"`
	StartedAt   time.Time `json:"started_at" db:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty" db:"finished_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	
	// Thread-safe access to game state
	mu sync.RWMutex
}

// GameEvent represents an event that occurred during the game
type GameEvent struct {
	ID        string    `json:"id" db:"id"`
	GameID    string    `json:"game_id" db:"game_id"`
	PlayerID  string    `json:"player_id" db:"player_id"`
	EventType string    `json:"event_type" db:"event_type"`
	Score     int64     `json:"score" db:"score"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Data      string    `json:"data" db:"data"` // JSON string for additional data
}

// Custom errors for game operations
var (
	ErrGameNotFound     = errors.New("game not found")
	ErrGameAlreadyEnded = errors.New("game already ended")
	ErrInvalidPlayer    = errors.New("invalid player")
	ErrGameNotStarted   = errors.New("game not started")
)

// NewGame creates a new game between two players
func NewGame(player1ID, player2ID string) (*Game, error) {
	if player1ID == "" || player2ID == "" {
		return nil, ErrInvalidPlayer
	}
	
	if player1ID == player2ID {
		return nil, errors.New("player cannot play against themselves")
	}
	
	now := time.Now()
	return &Game{
		ID:        generateGameID(),
		Player1ID: player1ID,
		Player2ID: player2ID,
		State:     GameStateWaiting,
		Score1:    0,
		Score2:    0,
		CreatedAt: now,
	}, nil
}

// Start begins the game
func (g *Game) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.State != GameStateWaiting {
		return errors.New("game cannot be started from current state")
	}
	
	g.State = GameStatePlaying
	g.StartedAt = time.Now()
	return nil
}

// UpdateScore updates the score for a player
func (g *Game) UpdateScore(playerID string, score int64) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.State != GameStatePlaying {
		return ErrGameNotStarted
	}
	
	switch playerID {
	case g.Player1ID:
		g.Score1 = score
	case g.Player2ID:
		g.Score2 = score
	default:
		return ErrInvalidPlayer
	}
	
	return nil
}

// End finishes the game and determines the winner
func (g *Game) End() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.State != GameStatePlaying {
		return ErrGameNotStarted
	}
	
	g.State = GameStateFinished
	now := time.Now()
	g.FinishedAt = &now
	
	// Determine winner
	if g.Score1 > g.Score2 {
		g.WinnerID = &g.Player1ID
	} else if g.Score2 > g.Score1 {
		g.WinnerID = &g.Player2ID
	}
	// If scores are equal, it's a tie (WinnerID remains nil)
	
	return nil
}

// Cancel cancels the game
func (g *Game) Cancel() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	if g.State == GameStateFinished {
		return ErrGameAlreadyEnded
	}
	
	g.State = GameStateCancelled
	now := time.Now()
	g.FinishedAt = &now
	return nil
}

// GetWinner returns the winner ID or empty string if tie
func (g *Game) GetWinner() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	if g.WinnerID != nil {
		return *g.WinnerID
	}
	return ""
}

// IsPlayer checks if the given user ID is a player in this game
func (g *Game) IsPlayer(userID string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	return userID == g.Player1ID || userID == g.Player2ID
}

// GetOpponent returns the opponent's ID for a given player
func (g *Game) GetOpponent(playerID string) (string, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	switch playerID {
	case g.Player1ID:
		return g.Player2ID, nil
	case g.Player2ID:
		return g.Player1ID, nil
	default:
		return "", ErrInvalidPlayer
	}
}

// GetDuration returns the duration of the game
func (g *Game) GetDuration() time.Duration {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	if g.StartedAt.IsZero() {
		return 0
	}
	
	endTime := g.StartedAt
	if g.FinishedAt != nil {
		endTime = *g.FinishedAt
	}
	
	return endTime.Sub(g.StartedAt)
}

// GetScore returns the current score for a player
func (g *Game) GetScore(playerID string) (int64, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	switch playerID {
	case g.Player1ID:
		return g.Score1, nil
	case g.Player2ID:
		return g.Score2, nil
	default:
		return 0, ErrInvalidPlayer
	}
}

// Helper function to generate game ID
func generateGameID() string {
	return "game_" + time.Now().Format("20060102150405")
}
