package tests

import (
	"strings"
	"testing"

	"effective-golang/internal/models"
)

// TestNewUser tests user creation with various inputs
func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid user",
			username: "testuser",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "username too short",
			username: "ab",
			email:    "test@example.com",
			password: "password123",
			wantErr:  true,
			errMsg:   "invalid username",
		},
		{
			name:     "username too long",
			username: "verylongusernamethatexceedslimit",
			email:    "test@example.com",
			password: "password123",
			wantErr:  true,
			errMsg:   "invalid username",
		},
		{
			name:     "invalid email",
			username: "testuser",
			email:    "invalid-email",
			password: "password123",
			wantErr:  true,
			errMsg:   "invalid email",
		},
		{
			name:     "password too short",
			username: "testuser",
			email:    "test@example.com",
			password: "123",
			wantErr:  true,
			errMsg:   "password too short",
		},
		{
			name:     "empty username",
			username: "",
			email:    "test@example.com",
			password: "password123",
			wantErr:  true,
			errMsg:   "invalid username",
		},
		{
			name:     "empty email",
			username: "testuser",
			email:    "",
			password: "password123",
			wantErr:  true,
			errMsg:   "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := models.NewUser(tt.username, tt.email, tt.password)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewUser() expected error but got none")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("NewUser() error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}
			
			if user == nil {
				t.Errorf("NewUser() returned nil user")
				return
			}
			
			// Verify user fields
			if user.Username != tt.username {
				t.Errorf("NewUser() username = %v, want %v", user.Username, tt.username)
			}
			
			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}
			
			if user.Password == "" {
				t.Errorf("NewUser() password should not be empty")
			}
			
			if user.ID == "" {
				t.Errorf("NewUser() ID should not be empty")
			}
			
			if !user.IsActive {
				t.Errorf("NewUser() user should be active by default")
			}
			
			// Verify timestamps are set
			if user.CreatedAt.IsZero() {
				t.Errorf("NewUser() CreatedAt should be set")
			}
			
			if user.UpdatedAt.IsZero() {
				t.Errorf("NewUser() UpdatedAt should be set")
			}
		})
	}
}

// TestUserStats tests user statistics functionality
func TestUserStats(t *testing.T) {
	stats := &models.UserStats{
		UserID:       "user123",
		TotalGames:   10,
		Wins:         7,
		Losses:       3,
		TotalScore:   1500,
		AverageScore: 150.0,
		Rank:         5,
	}
	
	t.Run("GetWinRate", func(t *testing.T) {
		expected := 70.0 // 7 wins / 10 total games * 100
		got := stats.GetWinRate()
		
		if got != expected {
			t.Errorf("GetWinRate() = %v, want %v", got, expected)
		}
	})
	
	t.Run("GetWinRate zero games", func(t *testing.T) {
		emptyStats := &models.UserStats{}
		expected := 0.0
		got := emptyStats.GetWinRate()
		
		if got != expected {
			t.Errorf("GetWinRate() with zero games = %v, want %v", got, expected)
		}
	})
	
	t.Run("GetAverageScore", func(t *testing.T) {
		expected := 150.0 // 1500 total score / 10 total games
		got := stats.GetAverageScore()
		
		if got != expected {
			t.Errorf("GetAverageScore() = %v, want %v", got, expected)
		}
	})
	
	t.Run("GetAverageScore zero games", func(t *testing.T) {
		emptyStats := &models.UserStats{}
		expected := 0.0
		got := emptyStats.GetAverageScore()
		
		if got != expected {
			t.Errorf("GetAverageScore() with zero games = %v, want %v", got, expected)
		}
	})
}

// TestUserStatsUpdate tests updating user statistics
func TestUserStatsUpdate(t *testing.T) {
	stats := &models.UserStats{
		UserID:       "user123",
		TotalGames:   5,
		Wins:         3,
		Losses:       2,
		TotalScore:   500,
		AverageScore: 100.0,
		Rank:         10,
	}
	
	t.Run("UpdateStats win", func(t *testing.T) {
		originalGames := stats.TotalGames
		originalWins := stats.Wins
		originalScore := stats.TotalScore
		
		stats.UpdateStats(200, true)
		
		if stats.TotalGames != originalGames+1 {
			t.Errorf("UpdateStats() TotalGames = %v, want %v", stats.TotalGames, originalGames+1)
		}
		
		if stats.Wins != originalWins+1 {
			t.Errorf("UpdateStats() Wins = %v, want %v", stats.Wins, originalWins+1)
		}
		
		if stats.TotalScore != originalScore+200 {
			t.Errorf("UpdateStats() TotalScore = %v, want %v", stats.TotalScore, originalScore+200)
		}
		
		expectedAvg := float64(stats.TotalScore) / float64(stats.TotalGames)
		if stats.AverageScore != expectedAvg {
			t.Errorf("UpdateStats() AverageScore = %v, want %v", stats.AverageScore, expectedAvg)
		}
	})
	
	t.Run("UpdateStats loss", func(t *testing.T) {
		originalGames := stats.TotalGames
		originalLosses := stats.Losses
		originalScore := stats.TotalScore
		
		stats.UpdateStats(150, false)
		
		if stats.TotalGames != originalGames+1 {
			t.Errorf("UpdateStats() TotalGames = %v, want %v", stats.TotalGames, originalGames+1)
		}
		
		if stats.Losses != originalLosses+1 {
			t.Errorf("UpdateStats() Losses = %v, want %v", stats.Losses, originalLosses+1)
		}
		
		if stats.TotalScore != originalScore+150 {
			t.Errorf("UpdateStats() TotalScore = %v, want %v", stats.TotalScore, originalScore+150)
		}
	})
}

// TestUserValidation tests user validation functions
func TestUserValidation(t *testing.T) {
	t.Run("validateUsername", func(t *testing.T) {
		tests := []struct {
			username string
			wantErr  bool
		}{
			{"valid", false},
			{"user123", false},
			{"a", true},      // too short
			{"verylongusernamethatexceedslimit", true}, // too long
			{"", true},       // empty
		}
		
		for _, tt := range tests {
			err := validateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUsername(%q) = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		}
	})
	
	t.Run("validateEmail", func(t *testing.T) {
		tests := []struct {
			email   string
			wantErr bool
		}{
			{"test@example.com", false},
			{"user@domain.org", false},
			{"invalid-email", true},
			{"@domain.com", true},
			{"user@", true},
			{"", true},
		}
		
		for _, tt := range tests {
			err := validateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail(%q) = %v, wantErr %v", tt.email, err, tt.wantErr)
			}
		}
	})
	
	t.Run("validatePassword", func(t *testing.T) {
		tests := []struct {
			password string
			wantErr  bool
		}{
			{"password123", false},
			{"123456", false},
			{"123", true},    // too short
			{"", true},       // empty
		}
		
		for _, tt := range tests {
			err := validatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePassword(%q) = %v, wantErr %v", tt.password, err, tt.wantErr)
			}
		}
	})
}

// Benchmark tests for performance testing
func BenchmarkNewUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := models.NewUser("testuser", "test@example.com", "password123")
		if err != nil {
			b.Fatalf("NewUser failed: %v", err)
		}
	}
}

func BenchmarkUserStatsUpdate(b *testing.B) {
	stats := &models.UserStats{
		UserID:       "user123",
		TotalGames:   100,
		Wins:         70,
		Losses:       30,
		TotalScore:   15000,
		AverageScore: 150.0,
		Rank:         5,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stats.UpdateStats(200, true)
	}
}

func BenchmarkGetWinRate(b *testing.B) {
	stats := &models.UserStats{
		UserID:     "user123",
		TotalGames: 100,
		Wins:       70,
		Losses:     30,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = stats.GetWinRate()
	}
}

// Helper functions for testing (these would be private in the actual implementation)
func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return models.ErrInvalidUsername
	}
	return nil
}

func validateEmail(email string) error {
	if len(email) < 5 || !contains(email, "@") {
		return models.ErrInvalidEmail
	}
	
	// Check for basic email format
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return models.ErrInvalidEmail
	}
	
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return models.ErrInvalidEmail
	}
	
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return models.ErrInvalidPassword
	}
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		contains(s[1:len(s)-1], substr))))
}
