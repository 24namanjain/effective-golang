package models

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// LeaderboardType represents different types of leaderboards
type LeaderboardType string

const (
	LeaderboardTypeGlobal    LeaderboardType = "global"
	LeaderboardTypeWeekly    LeaderboardType = "weekly"
	LeaderboardTypeMonthly   LeaderboardType = "monthly"
	LeaderboardTypeSeasonal  LeaderboardType = "seasonal"
)

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	UserID    string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Score     int64     `json:"score" db:"score"`
	Rank      int       `json:"rank" db:"rank"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Leaderboard represents a leaderboard with entries
type Leaderboard struct {
	ID          string           `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Type        LeaderboardType  `json:"type" db:"type"`
	Entries     []LeaderboardEntry `json:"entries" db:"entries"`
	MaxEntries  int              `json:"max_entries" db:"max_entries"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	
	// Thread-safe access to leaderboard data
	mu sync.RWMutex
}

// LeaderboardStats contains statistics about the leaderboard
type LeaderboardStats struct {
	TotalEntries    int     `json:"total_entries"`
	AverageScore    float64 `json:"average_score"`
	HighestScore    int64   `json:"highest_score"`
	LowestScore     int64   `json:"lowest_score"`
	LastUpdated     time.Time `json:"last_updated"`
}

// Custom errors for leaderboard operations
var (
	ErrLeaderboardNotFound = errors.New("leaderboard not found")
	ErrInvalidScore        = errors.New("invalid score")
	ErrUserNotFoundInLeaderboard = errors.New("user not found in leaderboard")
	ErrLeaderboardFull     = errors.New("leaderboard is full")
)

// NewLeaderboard creates a new leaderboard
func NewLeaderboard(name string, leaderboardType LeaderboardType, maxEntries int) *Leaderboard {
	now := time.Now()
	return &Leaderboard{
		ID:         generateLeaderboardID(),
		Name:       name,
		Type:       leaderboardType,
		Entries:    make([]LeaderboardEntry, 0),
		MaxEntries: maxEntries,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// AddEntry adds or updates an entry in the leaderboard
func (l *Leaderboard) AddEntry(userID, username string, score int64) error {
	if score < 0 {
		return ErrInvalidScore
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// Check if user already exists
	for i, entry := range l.Entries {
		if entry.UserID == userID {
			// Update existing entry
			l.Entries[i].Score = score
			l.Entries[i].Username = username
			l.Entries[i].UpdatedAt = time.Now()
			l.sortAndUpdateRanks()
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	
	// Add new entry
	if len(l.Entries) >= l.MaxEntries {
		// Check if new score is higher than lowest score
		if len(l.Entries) > 0 && score <= l.Entries[len(l.Entries)-1].Score {
			return ErrLeaderboardFull
		}
		
		// Remove lowest score entry
		l.Entries = l.Entries[:len(l.Entries)-1]
	}
	
	newEntry := LeaderboardEntry{
		UserID:    userID,
		Username:  username,
		Score:     score,
		UpdatedAt: time.Now(),
	}
	
	l.Entries = append(l.Entries, newEntry)
	l.sortAndUpdateRanks()
	l.UpdatedAt = time.Now()
	
	return nil
}

// GetUserRank returns the rank of a user in the leaderboard
func (l *Leaderboard) GetUserRank(userID string) (int, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	for _, entry := range l.Entries {
		if entry.UserID == userID {
			return entry.Rank, nil
		}
	}
	
	return 0, ErrUserNotFoundInLeaderboard
}

// GetTopEntries returns the top N entries from the leaderboard
func (l *Leaderboard) GetTopEntries(count int) []LeaderboardEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if count <= 0 {
		return []LeaderboardEntry{}
	}
	
	if count > len(l.Entries) {
		count = len(l.Entries)
	}
	
	result := make([]LeaderboardEntry, count)
	copy(result, l.Entries[:count])
	return result
}

// GetUserEntry returns the entry for a specific user
func (l *Leaderboard) GetUserEntry(userID string) (*LeaderboardEntry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	for _, entry := range l.Entries {
		if entry.UserID == userID {
			return &entry, nil
		}
	}
	
	return nil, ErrUserNotFoundInLeaderboard
}

// GetStats returns statistics about the leaderboard
func (l *Leaderboard) GetStats() *LeaderboardStats {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	if len(l.Entries) == 0 {
		return &LeaderboardStats{
			TotalEntries: 0,
			AverageScore: 0,
			HighestScore: 0,
			LowestScore:  0,
			LastUpdated:  l.UpdatedAt,
		}
	}
	
	var totalScore int64
	highestScore := l.Entries[0].Score
	lowestScore := l.Entries[len(l.Entries)-1].Score
	
	for _, entry := range l.Entries {
		totalScore += entry.Score
	}
	
	return &LeaderboardStats{
		TotalEntries: len(l.Entries),
		AverageScore: float64(totalScore) / float64(len(l.Entries)),
		HighestScore: highestScore,
		LowestScore:  lowestScore,
		LastUpdated:  l.UpdatedAt,
	}
}

// Clear removes all entries from the leaderboard
func (l *Leaderboard) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	l.Entries = make([]LeaderboardEntry, 0)
	l.UpdatedAt = time.Now()
}

// RemoveUser removes a user from the leaderboard
func (l *Leaderboard) RemoveUser(userID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	for i, entry := range l.Entries {
		if entry.UserID == userID {
			// Remove entry
			l.Entries = append(l.Entries[:i], l.Entries[i+1:]...)
			l.sortAndUpdateRanks()
			l.UpdatedAt = time.Now()
			return nil
		}
	}
	
	return ErrUserNotFoundInLeaderboard
}

// sortAndUpdateRanks sorts entries by score (descending) and updates ranks
func (l *Leaderboard) sortAndUpdateRanks() {
	// Sort by score in descending order
	sort.Slice(l.Entries, func(i, j int) bool {
		return l.Entries[i].Score > l.Entries[j].Score
	})
	
	// Update ranks
	for i := range l.Entries {
		l.Entries[i].Rank = i + 1
	}
}

// Helper function to generate leaderboard ID
func generateLeaderboardID() string {
	return "lb_" + time.Now().Format("20060102150405")
}
