package leaderboard

import (
	"context"
	"fmt"
	"sync"
	"time"

	"effective-golang/internal/models"
)

// LeaderboardService handles leaderboard operations and caching
type LeaderboardService struct {
	leaderboardRepo models.LeaderboardRepository
	userRepo        models.UserRepository
	cacheRepo       models.CacheRepository
	
	// Cache for leaderboard data
	cacheMutex      sync.RWMutex
	cacheTTL        int
	
	// Real-time updates
	updateChannels  map[string]chan *LeaderboardUpdate
	channelMutex    sync.RWMutex
}

// LeaderboardUpdate represents a leaderboard update
type LeaderboardUpdate struct {
	LeaderboardID string                    `json:"leaderboard_id"`
	Type          string                    `json:"type"`
	Entries       []models.LeaderboardEntry `json:"entries,omitempty"`
	UserID        string                    `json:"user_id,omitempty"`
	NewRank       int                       `json:"new_rank,omitempty"`
	OldRank       int                       `json:"old_rank,omitempty"`
	Timestamp     time.Time                 `json:"timestamp"`
}

// LeaderboardStats represents aggregated leaderboard statistics
type LeaderboardStats struct {
	TotalUsers     int     `json:"total_users"`
	AverageScore   float64 `json:"average_score"`
	HighestScore   int64   `json:"highest_score"`
	LowestScore    int64   `json:"lowest_score"`
	ScoreRange     int64   `json:"score_range"`
	LastUpdated    time.Time `json:"last_updated"`
}

// Custom errors for leaderboard operations
var (
	ErrLeaderboardNotFound = fmt.Errorf("leaderboard not found")
	ErrInvalidScore        = fmt.Errorf("invalid score")
	ErrUserNotFound        = fmt.Errorf("user not found")
	ErrCacheMiss           = fmt.Errorf("cache miss")
)

// NewLeaderboardService creates a new leaderboard service
func NewLeaderboardService(
	leaderboardRepo models.LeaderboardRepository,
	userRepo models.UserRepository,
	cacheRepo models.CacheRepository,
	cacheTTL int,
) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
		userRepo:        userRepo,
		cacheRepo:       cacheRepo,
		cacheTTL:        cacheTTL,
		updateChannels:  make(map[string]chan *LeaderboardUpdate),
	}
}

// CreateLeaderboard creates a new leaderboard
func (s *LeaderboardService) CreateLeaderboard(
	ctx context.Context,
	name string,
	leaderboardType models.LeaderboardType,
	maxEntries int,
) (*models.Leaderboard, error) {
	// Check if leaderboard already exists
	existing, err := s.leaderboardRepo.GetByName(ctx, name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("leaderboard already exists: %w", ErrLeaderboardNotFound)
	}
	
	// Create new leaderboard
	leaderboard := models.NewLeaderboard(name, leaderboardType, maxEntries)
	
	// Save to database
	if err := s.leaderboardRepo.Create(ctx, leaderboard); err != nil {
		return nil, fmt.Errorf("failed to create leaderboard: %w", err)
	}
	
	// Initialize cache
	s.cacheLeaderboard(ctx, leaderboard)
	
	// Create update channel
	s.channelMutex.Lock()
	s.updateChannels[leaderboard.ID] = make(chan *LeaderboardUpdate, 100)
	s.channelMutex.Unlock()
	
	return leaderboard, nil
}

// AddScore adds or updates a score in a leaderboard
func (s *LeaderboardService) AddScore(
	ctx context.Context,
	leaderboardID, userID string,
	score int64,
) error {
	if score < 0 {
		return ErrInvalidScore
	}
	
	// Get user information
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	
	// Get current rank (if any)
	oldRank, _ := s.leaderboardRepo.GetUserRank(ctx, leaderboardID, userID)
	
	// Add entry to leaderboard
	entry := &models.LeaderboardEntry{
		UserID:    userID,
		Username:  user.Username,
		Score:     score,
		UpdatedAt: time.Now(),
	}
	
	if err := s.leaderboardRepo.AddEntry(ctx, leaderboardID, entry); err != nil {
		return fmt.Errorf("failed to add entry: %w", err)
	}
	
	// Get new rank
	newRank, err := s.leaderboardRepo.GetUserRank(ctx, leaderboardID, userID)
	if err != nil {
		return fmt.Errorf("failed to get new rank: %w", err)
	}
	
	// Invalidate cache
	s.invalidateCache(ctx, leaderboardID)
	
	// Send real-time update
	s.sendUpdate(&LeaderboardUpdate{
		LeaderboardID: leaderboardID,
		Type:          "score_updated",
		UserID:        userID,
		NewRank:       newRank,
		OldRank:       oldRank,
		Timestamp:     time.Now(),
	})
	
	return nil
}

// GetTopEntries retrieves top entries from a leaderboard
func (s *LeaderboardService) GetTopEntries(
	ctx context.Context,
	leaderboardID string,
	count int,
) ([]models.LeaderboardEntry, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("leaderboard:%s:top:%d", leaderboardID, count)
	var entries []models.LeaderboardEntry
	
	if err := s.cacheRepo.Get(ctx, cacheKey, &entries); err == nil {
		return entries, nil
	}
	
	// Get from database
	repoEntries, err := s.leaderboardRepo.GetTopEntries(ctx, leaderboardID, count)
	if err != nil {
		return nil, fmt.Errorf("failed to get top entries: %w", err)
	}
	
	// Convert to value slice for caching
	entryValues := make([]models.LeaderboardEntry, len(repoEntries))
	for i, entry := range repoEntries {
		entryValues[i] = *entry
	}
	
	// Cache the result
	s.cacheRepo.Set(ctx, cacheKey, entryValues, s.cacheTTL)
	
	return entryValues, nil
}

// GetUserRank retrieves a user's rank in a leaderboard
func (s *LeaderboardService) GetUserRank(
	ctx context.Context,
	leaderboardID, userID string,
) (int, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("leaderboard:%s:rank:%s", leaderboardID, userID)
	var rank int
	
	if err := s.cacheRepo.Get(ctx, cacheKey, &rank); err == nil {
		return rank, nil
	}
	
	// Get from database
	rank, err := s.leaderboardRepo.GetUserRank(ctx, leaderboardID, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user rank: %w", err)
	}
	
	// Cache the result
	s.cacheRepo.Set(ctx, cacheKey, rank, s.cacheTTL)
	
	return rank, nil
}

// GetLeaderboard retrieves a complete leaderboard
func (s *LeaderboardService) GetLeaderboard(
	ctx context.Context,
	leaderboardID string,
) (*models.Leaderboard, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("leaderboard:%s", leaderboardID)
	var leaderboard models.Leaderboard
	
	if err := s.cacheRepo.Get(ctx, cacheKey, &leaderboard); err == nil {
		return &leaderboard, nil
	}
	
	// Get from database
	repoLeaderboard, err := s.leaderboardRepo.GetByID(ctx, leaderboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	
	// Cache the result
	s.cacheLeaderboard(ctx, repoLeaderboard)
	
	return repoLeaderboard, nil
}

// GetStats retrieves leaderboard statistics
func (s *LeaderboardService) GetStats(
	ctx context.Context,
	leaderboardID string,
) (*LeaderboardStats, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("leaderboard:%s:stats", leaderboardID)
	var stats LeaderboardStats
	
	if err := s.cacheRepo.Get(ctx, cacheKey, &stats); err == nil {
		return &stats, nil
	}
	
	// Get leaderboard
	leaderboard, err := s.GetLeaderboard(ctx, leaderboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	
	// Calculate statistics
	stats = s.calculateStats(leaderboard)
	
	// Cache the result
	s.cacheRepo.Set(ctx, cacheKey, stats, s.cacheTTL)
	
	return &stats, nil
}

// SubscribeToUpdates subscribes to real-time leaderboard updates
func (s *LeaderboardService) SubscribeToUpdates(leaderboardID string) (<-chan *LeaderboardUpdate, error) {
	s.channelMutex.RLock()
	channel, exists := s.updateChannels[leaderboardID]
	s.channelMutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("leaderboard not found: %w", ErrLeaderboardNotFound)
	}
	
	return channel, nil
}

// UnsubscribeFromUpdates unsubscribes from real-time updates
func (s *LeaderboardService) UnsubscribeFromUpdates(leaderboardID string) {
	s.channelMutex.Lock()
	if channel, exists := s.updateChannels[leaderboardID]; exists {
		close(channel)
		delete(s.updateChannels, leaderboardID)
	}
	s.channelMutex.Unlock()
}

// RefreshLeaderboard refreshes leaderboard data from database
func (s *LeaderboardService) RefreshLeaderboard(ctx context.Context, leaderboardID string) error {
	// Get fresh data from database
	leaderboard, err := s.leaderboardRepo.GetByID(ctx, leaderboardID)
	if err != nil {
		return fmt.Errorf("failed to get leaderboard: %w", err)
	}
	
	// Update cache
	s.cacheLeaderboard(ctx, leaderboard)
	
	// Send refresh update
	s.sendUpdate(&LeaderboardUpdate{
		LeaderboardID: leaderboardID,
		Type:          "refreshed",
		Entries:       leaderboard.Entries,
		Timestamp:     time.Now(),
	})
	
	return nil
}

// cacheLeaderboard caches leaderboard data
func (s *LeaderboardService) cacheLeaderboard(ctx context.Context, leaderboard *models.Leaderboard) {
	cacheKey := fmt.Sprintf("leaderboard:%s", leaderboard.ID)
	s.cacheRepo.Set(ctx, cacheKey, leaderboard, s.cacheTTL)
}

// invalidateCache invalidates cached data for a leaderboard
func (s *LeaderboardService) invalidateCache(ctx context.Context, leaderboardID string) {
	// Remove main leaderboard cache
	cacheKey := fmt.Sprintf("leaderboard:%s", leaderboardID)
	s.cacheRepo.Delete(ctx, cacheKey)
	
	// Remove stats cache
	statsKey := fmt.Sprintf("leaderboard:%s:stats", leaderboardID)
	s.cacheRepo.Delete(ctx, statsKey)
	
	// Remove top entries cache (pattern matching would be better in real implementation)
	for i := 1; i <= 100; i++ {
		topKey := fmt.Sprintf("leaderboard:%s:top:%d", leaderboardID, i)
		s.cacheRepo.Delete(ctx, topKey)
	}
}

// sendUpdate sends a real-time update to subscribers
func (s *LeaderboardService) sendUpdate(update *LeaderboardUpdate) {
	s.channelMutex.RLock()
	channel, exists := s.updateChannels[update.LeaderboardID]
	s.channelMutex.RUnlock()
	
	if exists {
		select {
		case channel <- update:
			// Update sent successfully
		default:
			// Channel is full, skip update
		}
	}
}

// calculateStats calculates leaderboard statistics
func (s *LeaderboardService) calculateStats(leaderboard *models.Leaderboard) LeaderboardStats {
	if len(leaderboard.Entries) == 0 {
		return LeaderboardStats{
			TotalUsers:   0,
			AverageScore: 0,
			HighestScore: 0,
			LowestScore:  0,
			ScoreRange:   0,
			LastUpdated:  leaderboard.UpdatedAt,
		}
	}
	
	var totalScore int64
	highestScore := leaderboard.Entries[0].Score
	lowestScore := leaderboard.Entries[len(leaderboard.Entries)-1].Score
	
	for _, entry := range leaderboard.Entries {
		totalScore += entry.Score
	}
	
	return LeaderboardStats{
		TotalUsers:   len(leaderboard.Entries),
		AverageScore: float64(totalScore) / float64(len(leaderboard.Entries)),
		HighestScore: highestScore,
		LowestScore:  lowestScore,
		ScoreRange:   highestScore - lowestScore,
		LastUpdated:  leaderboard.UpdatedAt,
	}
}

// GetLeaderboardsByType retrieves leaderboards by type
func (s *LeaderboardService) GetLeaderboardsByType(
	ctx context.Context,
	leaderboardType models.LeaderboardType,
) ([]*models.Leaderboard, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("leaderboards:type:%s", leaderboardType)
	var leaderboards []*models.Leaderboard
	
	if err := s.cacheRepo.Get(ctx, cacheKey, &leaderboards); err == nil {
		return leaderboards, nil
	}
	
	// Get from database
	leaderboards, err := s.leaderboardRepo.GetByType(ctx, leaderboardType)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboards: %w", err)
	}
	
	// Cache the result
	s.cacheRepo.Set(ctx, cacheKey, leaderboards, s.cacheTTL)
	
	return leaderboards, nil
}

// Close closes the leaderboard service
func (s *LeaderboardService) Close() {
	s.channelMutex.Lock()
	defer s.channelMutex.Unlock()
	
	for _, channel := range s.updateChannels {
		close(channel)
	}
	s.updateChannels = make(map[string]chan *LeaderboardUpdate)
}
