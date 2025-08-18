package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"effective-golang/internal/models"
)

// AuthService handles authentication and authorization logic
type AuthService struct {
	userRepo models.UserRepository
	cacheRepo models.CacheRepository
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Custom errors for authentication
var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrSessionExpired     = fmt.Errorf("session expired")
	ErrSessionNotFound    = fmt.Errorf("session not found")
	ErrUserAlreadyExists  = fmt.Errorf("user already exists")
)

// NewAuthService creates a new authentication service
func NewAuthService(userRepo models.UserRepository, cacheRepo models.CacheRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already taken: %w", ErrUserAlreadyExists)
	}
	
	existingUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already registered: %w", ErrUserAlreadyExists)
	}
	
	// Create new user
	user, err := models.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	
	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	
	// Initialize user stats
	stats := &models.UserStats{
		UserID:       user.ID,
		TotalGames:   0,
		Wins:         0,
		Losses:       0,
		TotalScore:   0,
		AverageScore: 0,
		Rank:         0,
	}
	
	if err := s.userRepo.UpdateStats(ctx, stats); err != nil {
		return nil, fmt.Errorf("failed to initialize user stats: %w", err)
	}
	
	return user, nil
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*Session, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", ErrInvalidCredentials)
	}
	
	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}
	
	// Verify password (in real app, use bcrypt.CompareHashAndPassword)
	if user.Password != req.Password {
		return nil, fmt.Errorf("authentication failed: %w", ErrInvalidCredentials)
	}
	
	// Create session
	session, err := s.createSession(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	
	return session, nil
}

// Logout invalidates a user session
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	// Remove session from cache
	cacheKey := fmt.Sprintf("session:%s", sessionID)
	if err := s.cacheRepo.Delete(ctx, cacheKey); err != nil {
		return fmt.Errorf("failed to remove session: %w", err)
	}
	
	return nil
}

// ValidateSession validates a session and returns user information
func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (*Session, error) {
	cacheKey := fmt.Sprintf("session:%s", sessionID)
	
	var session Session
	if err := s.cacheRepo.Get(ctx, cacheKey, &session); err != nil {
		return nil, fmt.Errorf("session validation failed: %w", ErrSessionNotFound)
	}
	
	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		s.cacheRepo.Delete(ctx, cacheKey)
		return nil, fmt.Errorf("session validation failed: %w", ErrSessionExpired)
	}
	
	return &session, nil
}

// GetUserBySession retrieves user information from a session
func (s *AuthService) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	session, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}
	
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

// RefreshSession extends a session's expiration time
func (s *AuthService) RefreshSession(ctx context.Context, sessionID string) (*Session, error) {
	session, err := s.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}
	
	// Extend session expiration
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	
	// Update session in cache
	cacheKey := fmt.Sprintf("session:%s", sessionID)
	if err := s.cacheRepo.Set(ctx, cacheKey, session, 86400); err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}
	
	return session, nil
}

// createSession creates a new session for a user
func (s *AuthService) createSession(ctx context.Context, user *models.User) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}
	
	now := time.Now()
	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		Username:  user.Username,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
	}
	
	// Store session in cache
	cacheKey := fmt.Sprintf("session:%s", sessionID)
	if err := s.cacheRepo.Set(ctx, cacheKey, session, 86400); err != nil {
		return nil, fmt.Errorf("failed to store session: %w", err)
	}
	
	return session, nil
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GetUserStats retrieves user statistics
func (s *AuthService) GetUserStats(ctx context.Context, userID string) (*models.UserStats, error) {
	stats, err := s.userRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}
	
	return stats, nil
}

// UpdateUserStats updates user statistics
func (s *AuthService) UpdateUserStats(ctx context.Context, stats *models.UserStats) error {
	if err := s.userRepo.UpdateStats(ctx, stats); err != nil {
		return fmt.Errorf("failed to update user stats: %w", err)
	}
	
	return nil
}
