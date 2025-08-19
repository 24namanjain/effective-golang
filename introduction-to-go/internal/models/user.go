package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// User represents a game user with authentication and profile information
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // "-" means this field won't be included in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
}

// UserStats contains user's game statistics
type UserStats struct {
	UserID       string  `json:"user_id" db:"user_id"`
	TotalGames   int     `json:"total_games" db:"total_games"`
	Wins         int     `json:"wins" db:"wins"`
	Losses       int     `json:"losses" db:"losses"`
	TotalScore   int64   `json:"total_score" db:"total_score"`
	AverageScore float64 `json:"average_score" db:"average_score"`
	Rank         int     `json:"rank" db:"rank"`
}

// Custom error types for better error handling
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPassword   = errors.New("password too short")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// NewUser creates a new user with validation
func NewUser(username, email, password string) (*User, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	
	now := time.Now()
	return &User{
		ID:        generateUserID(),
		Username:  username,
		Email:     email,
		Password:  hashPassword(password), // In real app, use bcrypt
		CreatedAt: now,
		UpdatedAt: now,
		IsActive:  true,
	}, nil
}

// GetWinRate calculates and returns the user's win rate
func (u *UserStats) GetWinRate() float64 {
	if u.TotalGames == 0 {
		return 0.0
	}
	return float64(u.Wins) / float64(u.TotalGames) * 100
}

// GetAverageScore calculates and returns the user's average score
func (u *UserStats) GetAverageScore() float64 {
	if u.TotalGames == 0 {
		return 0.0
	}
	return float64(u.TotalScore) / float64(u.TotalGames)
}

// UpdateStats updates user statistics after a game
func (u *UserStats) UpdateStats(score int64, won bool) {
	u.TotalGames++
	u.TotalScore += score
	
	if won {
		u.Wins++
	} else {
		u.Losses++
	}
	
	u.AverageScore = u.GetAverageScore()
}

// Validation functions
func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	return nil
}

func validateEmail(email string) error {
	// Simple email validation - in real app, use regex
	if len(email) < 5 || !contains(email, "@") {
		return ErrInvalidEmail
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return ErrInvalidPassword
	}
	return nil
}

// Helper functions (in real app, these would be more sophisticated)
func generateUserID() string {
	// Use high-resolution time and random bytes to avoid collisions
	raw := make([]byte, 6)
	_, _ = rand.Read(raw)
	return "user_" + time.Now().Format("20060102T150405.000000000") + "_" + hex.EncodeToString(raw)
}

func hashPassword(password string) string {
	// In real app, use bcrypt or similar
	return password // Simplified for learning
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
		s[len(s)-len(substr):] == substr ||
		contains(s[1:len(s)-1], substr))))
}
