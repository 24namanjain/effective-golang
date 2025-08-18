package models

import (
	"context"
	"fmt"
)

// Repository interfaces demonstrate the repository pattern
// This allows for easy testing and swapping implementations

// UserRepository defines operations for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*User, error)
	
	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *User) error
	
	// Delete removes a user
	Delete(ctx context.Context, id string) error
	
	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*User, error)
	
	// GetStats retrieves user statistics
	GetStats(ctx context.Context, userID string) (*UserStats, error)
	
	// UpdateStats updates user statistics
	UpdateStats(ctx context.Context, stats *UserStats) error
}

// GameRepository defines operations for game data access
type GameRepository interface {
	// Create creates a new game
	Create(ctx context.Context, game *Game) error
	
	// GetByID retrieves a game by ID
	GetByID(ctx context.Context, id string) (*Game, error)
	
	// Update updates an existing game
	Update(ctx context.Context, game *Game) error
	
	// Delete removes a game
	Delete(ctx context.Context, id string) error
	
	// GetUserGames retrieves games for a specific user
	GetUserGames(ctx context.Context, userID string, limit int) ([]*Game, error)
	
	// GetActiveGames retrieves active games
	GetActiveGames(ctx context.Context) ([]*Game, error)
	
	// AddEvent adds a game event
	AddEvent(ctx context.Context, event *GameEvent) error
	
	// GetGameEvents retrieves events for a game
	GetGameEvents(ctx context.Context, gameID string) ([]*GameEvent, error)
}

// LeaderboardRepository defines operations for leaderboard data access
type LeaderboardRepository interface {
	// Create creates a new leaderboard
	Create(ctx context.Context, leaderboard *Leaderboard) error
	
	// GetByID retrieves a leaderboard by ID
	GetByID(ctx context.Context, id string) (*Leaderboard, error)
	
	// GetByName retrieves a leaderboard by name
	GetByName(ctx context.Context, name string) (*Leaderboard, error)
	
	// Update updates an existing leaderboard
	Update(ctx context.Context, leaderboard *Leaderboard) error
	
	// Delete removes a leaderboard
	Delete(ctx context.Context, id string) error
	
	// List retrieves leaderboards with pagination
	List(ctx context.Context, offset, limit int) ([]*Leaderboard, error)
	
	// GetByType retrieves leaderboards by type
	GetByType(ctx context.Context, leaderboardType LeaderboardType) ([]*Leaderboard, error)
	
	// AddEntry adds an entry to a leaderboard
	AddEntry(ctx context.Context, leaderboardID string, entry *LeaderboardEntry) error
	
	// RemoveEntry removes an entry from a leaderboard
	RemoveEntry(ctx context.Context, leaderboardID, userID string) error
	
	// GetTopEntries retrieves top entries from a leaderboard
	GetTopEntries(ctx context.Context, leaderboardID string, count int) ([]*LeaderboardEntry, error)
	
	// GetUserRank retrieves a user's rank in a leaderboard
	GetUserRank(ctx context.Context, leaderboardID, userID string) (int, error)
}

// Custom errors for cache operations
var (
	ErrCacheMiss = fmt.Errorf("cache miss")
)

// CacheRepository defines operations for caching
type CacheRepository interface {
	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	
	// Get retrieves a value from cache
	Get(ctx context.Context, key string, dest interface{}) error
	
	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)
	
	// SetNX sets a value only if it doesn't exist
	SetNX(ctx context.Context, key string, value interface{}, ttl int) (bool, error)
	
	// Increment increments a numeric value
	Increment(ctx context.Context, key string, value int64) (int64, error)
	
	// Expire sets expiration for a key
	Expire(ctx context.Context, key string, ttl int) error
}

// TransactionManager defines operations for database transactions
type TransactionManager interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error
	
	// Rollback rolls back the transaction
	Rollback() error
	
	// UserRepository returns a user repository within the transaction
	UserRepository() UserRepository
	
	// GameRepository returns a game repository within the transaction
	GameRepository() GameRepository
	
	// LeaderboardRepository returns a leaderboard repository within the transaction
	LeaderboardRepository() LeaderboardRepository
}

// UnitOfWork pattern for managing multiple repositories
type UnitOfWork interface {
	// UserRepository returns the user repository
	UserRepository() UserRepository
	
	// GameRepository returns the game repository
	GameRepository() GameRepository
	
	// LeaderboardRepository returns the leaderboard repository
	LeaderboardRepository() LeaderboardRepository
	
	// CacheRepository returns the cache repository
	CacheRepository() CacheRepository
	
	// TransactionManager returns the transaction manager
	TransactionManager() TransactionManager
	
	// Begin starts a new transaction
	Begin(ctx context.Context) (Transaction, error)
	
	// Close closes all connections
	Close() error
}
