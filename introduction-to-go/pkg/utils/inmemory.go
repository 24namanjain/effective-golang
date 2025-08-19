package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"effective-golang/internal/models"
)

// InMemoryUnitOfWork implements UnitOfWork with in-memory storage
type InMemoryUnitOfWork struct {
	userRepo        *InMemoryUserRepository
	gameRepo        *InMemoryGameRepository
	leaderboardRepo *InMemoryLeaderboardRepository
	cacheRepo       *InMemoryCacheRepository
	txManager       *InMemoryTransactionManager
}

// NewInMemoryUnitOfWork creates a new in-memory unit of work
func NewInMemoryUnitOfWork() models.UnitOfWork {
	userRepo := &InMemoryUserRepository{
		users:  make(map[string]*models.User),
		stats:  make(map[string]*models.UserStats),
		mutex:  sync.RWMutex{},
	}
	
	gameRepo := &InMemoryGameRepository{
		games:  make(map[string]*models.Game),
		events: make(map[string][]*models.GameEvent),
		mutex:  sync.RWMutex{},
	}
	
	leaderboardRepo := &InMemoryLeaderboardRepository{
		leaderboards: make(map[string]*models.Leaderboard),
		mutex:        sync.RWMutex{},
	}
	
	cacheRepo := &InMemoryCacheRepository{
		data:   make(map[string]*cacheEntry),
		mutex:  sync.RWMutex{},
	}
	
	txManager := &InMemoryTransactionManager{
		unitOfWork: nil, // Will be set below
	}
	
	unitOfWork := &InMemoryUnitOfWork{
		userRepo:        userRepo,
		gameRepo:        gameRepo,
		leaderboardRepo: leaderboardRepo,
		cacheRepo:       cacheRepo,
		txManager:       txManager,
	}
	
	txManager.unitOfWork = unitOfWork
	
	return unitOfWork
}

// Repository implementations
func (uow *InMemoryUnitOfWork) UserRepository() models.UserRepository {
	return uow.userRepo
}

func (uow *InMemoryUnitOfWork) GameRepository() models.GameRepository {
	return uow.gameRepo
}

func (uow *InMemoryUnitOfWork) LeaderboardRepository() models.LeaderboardRepository {
	return uow.leaderboardRepo
}

func (uow *InMemoryUnitOfWork) CacheRepository() models.CacheRepository {
	return uow.cacheRepo
}

func (uow *InMemoryUnitOfWork) TransactionManager() models.TransactionManager {
	return uow.txManager
}

func (uow *InMemoryUnitOfWork) Begin(ctx context.Context) (models.Transaction, error) {
	return uow.txManager.Begin(ctx)
}

func (uow *InMemoryUnitOfWork) Close() error {
	return nil
}

// InMemoryUserRepository implements UserRepository with in-memory storage
type InMemoryUserRepository struct {
	users map[string]*models.User
	stats map[string]*models.UserStats
	mutex sync.RWMutex
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, models.ErrUserNotFound
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, models.ErrUserNotFound
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[user.ID]; !exists {
		return models.ErrUserNotFound
	}
	
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.users[id]; !exists {
		return models.ErrUserNotFound
	}
	
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	// Simple pagination
	if offset >= len(users) {
		return []*models.User{}, nil
	}
	
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}
	
	return users[offset:end], nil
}

func (r *InMemoryUserRepository) GetStats(ctx context.Context, userID string) (*models.UserStats, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	stats, exists := r.stats[userID]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return stats, nil
}

func (r *InMemoryUserRepository) UpdateStats(ctx context.Context, stats *models.UserStats) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.stats[stats.UserID] = stats
	return nil
}

// InMemoryGameRepository implements GameRepository with in-memory storage
type InMemoryGameRepository struct {
	games  map[string]*models.Game
	events map[string][]*models.GameEvent
	mutex  sync.RWMutex
}

func (r *InMemoryGameRepository) Create(ctx context.Context, game *models.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.games[game.ID] = game
	r.events[game.ID] = make([]*models.GameEvent, 0)
	return nil
}

func (r *InMemoryGameRepository) GetByID(ctx context.Context, id string) (*models.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	game, exists := r.games[id]
	if !exists {
		return nil, models.ErrGameNotFound
	}
	return game, nil
}

func (r *InMemoryGameRepository) Update(ctx context.Context, game *models.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.games[game.ID]; !exists {
		return models.ErrGameNotFound
	}
	
	r.games[game.ID] = game
	return nil
}

func (r *InMemoryGameRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.games[id]; !exists {
		return models.ErrGameNotFound
	}
	
	delete(r.games, id)
	delete(r.events, id)
	return nil
}

func (r *InMemoryGameRepository) GetUserGames(ctx context.Context, userID string, limit int) ([]*models.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var games []*models.Game
	for _, game := range r.games {
		if game.Player1ID == userID || game.Player2ID == userID {
			games = append(games, game)
		}
	}
	
	if limit > 0 && len(games) > limit {
		games = games[:limit]
	}
	
	return games, nil
}

func (r *InMemoryGameRepository) GetActiveGames(ctx context.Context) ([]*models.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var games []*models.Game
	for _, game := range r.games {
		if game.State == models.GameStatePlaying {
			games = append(games, game)
		}
	}
	
	return games, nil
}

func (r *InMemoryGameRepository) AddEvent(ctx context.Context, event *models.GameEvent) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.events[event.GameID] = append(r.events[event.GameID], event)
	return nil
}

func (r *InMemoryGameRepository) GetGameEvents(ctx context.Context, gameID string) ([]*models.GameEvent, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	events, exists := r.events[gameID]
	if !exists {
		return []*models.GameEvent{}, nil
	}
	
	return events, nil
}

// InMemoryLeaderboardRepository implements LeaderboardRepository with in-memory storage
type InMemoryLeaderboardRepository struct {
	leaderboards map[string]*models.Leaderboard
	mutex        sync.RWMutex
}

func (r *InMemoryLeaderboardRepository) Create(ctx context.Context, leaderboard *models.Leaderboard) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.leaderboards[leaderboard.ID] = leaderboard
	return nil
}

func (r *InMemoryLeaderboardRepository) GetByID(ctx context.Context, id string) (*models.Leaderboard, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	leaderboard, exists := r.leaderboards[id]
	if !exists {
		return nil, models.ErrLeaderboardNotFound
	}
	return leaderboard, nil
}

func (r *InMemoryLeaderboardRepository) GetByName(ctx context.Context, name string) (*models.Leaderboard, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, leaderboard := range r.leaderboards {
		if leaderboard.Name == name {
			return leaderboard, nil
		}
	}
	return nil, models.ErrLeaderboardNotFound
}

func (r *InMemoryLeaderboardRepository) Update(ctx context.Context, leaderboard *models.Leaderboard) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.leaderboards[leaderboard.ID]; !exists {
		return models.ErrLeaderboardNotFound
	}
	
	r.leaderboards[leaderboard.ID] = leaderboard
	return nil
}

func (r *InMemoryLeaderboardRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.leaderboards[id]; !exists {
		return models.ErrLeaderboardNotFound
	}
	
	delete(r.leaderboards, id)
	return nil
}

func (r *InMemoryLeaderboardRepository) List(ctx context.Context, offset, limit int) ([]*models.Leaderboard, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	leaderboards := make([]*models.Leaderboard, 0, len(r.leaderboards))
	for _, leaderboard := range r.leaderboards {
		leaderboards = append(leaderboards, leaderboard)
	}
	
	if offset >= len(leaderboards) {
		return []*models.Leaderboard{}, nil
	}
	
	end := offset + limit
	if end > len(leaderboards) {
		end = len(leaderboards)
	}
	
	return leaderboards[offset:end], nil
}

func (r *InMemoryLeaderboardRepository) GetByType(ctx context.Context, leaderboardType models.LeaderboardType) ([]*models.Leaderboard, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var leaderboards []*models.Leaderboard
	for _, leaderboard := range r.leaderboards {
		if leaderboard.Type == leaderboardType {
			leaderboards = append(leaderboards, leaderboard)
		}
	}
	
	return leaderboards, nil
}

func (r *InMemoryLeaderboardRepository) AddEntry(ctx context.Context, leaderboardID string, entry *models.LeaderboardEntry) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	leaderboard, exists := r.leaderboards[leaderboardID]
	if !exists {
		return models.ErrLeaderboardNotFound
	}
	
	return leaderboard.AddEntry(entry.UserID, entry.Username, entry.Score)
}

func (r *InMemoryLeaderboardRepository) RemoveEntry(ctx context.Context, leaderboardID, userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	leaderboard, exists := r.leaderboards[leaderboardID]
	if !exists {
		return models.ErrLeaderboardNotFound
	}
	
	return leaderboard.RemoveUser(userID)
}

func (r *InMemoryLeaderboardRepository) GetTopEntries(ctx context.Context, leaderboardID string, count int) ([]*models.LeaderboardEntry, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	leaderboard, exists := r.leaderboards[leaderboardID]
	if !exists {
		return nil, models.ErrLeaderboardNotFound
	}
	
	entries := leaderboard.GetTopEntries(count)
	result := make([]*models.LeaderboardEntry, len(entries))
	for i := range entries {
		result[i] = &entries[i]
	}
	return result, nil
}

func (r *InMemoryLeaderboardRepository) GetUserRank(ctx context.Context, leaderboardID, userID string) (int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	leaderboard, exists := r.leaderboards[leaderboardID]
	if !exists {
		return 0, models.ErrLeaderboardNotFound
	}
	
	return leaderboard.GetUserRank(userID)
}

// InMemoryCacheRepository implements CacheRepository with in-memory storage
type InMemoryCacheRepository struct {
	data  map[string]*cacheEntry
	mutex sync.RWMutex
}

type cacheEntry struct {
	value      interface{}
	expiration time.Time
}

func (r *InMemoryCacheRepository) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	expiration := time.Now().Add(time.Duration(ttl) * time.Second)
	r.data[key] = &cacheEntry{
		value:      value,
		expiration: expiration,
	}
	return nil
}

func (r *InMemoryCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	entry, exists := r.data[key]
	if !exists {
		return models.ErrCacheMiss
	}
	
	if time.Now().After(entry.expiration) {
		delete(r.data, key)
		return models.ErrCacheMiss
	}
	
	// Simple JSON marshaling/unmarshaling for deep copy
	data, err := json.Marshal(entry.value)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, dest)
}

func (r *InMemoryCacheRepository) Delete(ctx context.Context, key string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	delete(r.data, key)
	return nil
}

func (r *InMemoryCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	entry, exists := r.data[key]
	if !exists {
		return false, nil
	}
	
	if time.Now().After(entry.expiration) {
		delete(r.data, key)
		return false, nil
	}
	
	return true, nil
}

func (r *InMemoryCacheRepository) SetNX(ctx context.Context, key string, value interface{}, ttl int) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.data[key]; exists {
		return false, nil
	}
	
	expiration := time.Now().Add(time.Duration(ttl) * time.Second)
	r.data[key] = &cacheEntry{
		value:      value,
		expiration: expiration,
	}
	return true, nil
}

func (r *InMemoryCacheRepository) Increment(ctx context.Context, key string, value int64) (int64, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	entry, exists := r.data[key]
	if !exists {
		r.data[key] = &cacheEntry{
			value:      value,
			expiration: time.Now().Add(24 * time.Hour),
		}
		return value, nil
	}
	
	if time.Now().After(entry.expiration) {
		entry.value = value
		entry.expiration = time.Now().Add(24 * time.Hour)
		return value, nil
	}
	
	if current, ok := entry.value.(int64); ok {
		newValue := current + value
		entry.value = newValue
		return newValue, nil
	}
	
	return 0, fmt.Errorf("value is not an integer")
}

func (r *InMemoryCacheRepository) Expire(ctx context.Context, key string, ttl int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	entry, exists := r.data[key]
	if !exists {
		return models.ErrCacheMiss
	}
	
	entry.expiration = time.Now().Add(time.Duration(ttl) * time.Second)
	return nil
}

// InMemoryTransactionManager implements TransactionManager with in-memory storage
type InMemoryTransactionManager struct {
	unitOfWork *InMemoryUnitOfWork
}

func (tm *InMemoryTransactionManager) Begin(ctx context.Context) (models.Transaction, error) {
	// For in-memory implementation, we just return the same unit of work
	// In a real implementation, this would create a transaction context
	return &InMemoryTransaction{
		unitOfWork: tm.unitOfWork,
	}, nil
}

// InMemoryTransaction implements Transaction with in-memory storage
type InMemoryTransaction struct {
	unitOfWork *InMemoryUnitOfWork
	committed  bool
	rolledBack bool
}

func (tx *InMemoryTransaction) Commit() error {
	if tx.rolledBack {
		return fmt.Errorf("transaction already rolled back")
	}
	tx.committed = true
	return nil
}

func (tx *InMemoryTransaction) Rollback() error {
	if tx.committed {
		return fmt.Errorf("transaction already committed")
	}
	tx.rolledBack = true
	return nil
}

func (tx *InMemoryTransaction) UserRepository() models.UserRepository {
	return tx.unitOfWork.userRepo
}

func (tx *InMemoryTransaction) GameRepository() models.GameRepository {
	return tx.unitOfWork.gameRepo
}

func (tx *InMemoryTransaction) LeaderboardRepository() models.LeaderboardRepository {
	return tx.unitOfWork.leaderboardRepo
}
