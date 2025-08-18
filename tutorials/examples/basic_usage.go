package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"effective-golang/internal/auth"
	"effective-golang/internal/game"
	"effective-golang/internal/leaderboard"
	"effective-golang/internal/models"
	"effective-golang/pkg/utils"
)

// This example demonstrates how to use the game leaderboard system
func demonstrateBasicUsage() {
	// Initialize in-memory repositories
	unitOfWork := utils.NewInMemoryUnitOfWork()
	
	// Initialize services
	authService := auth.NewAuthService(
		unitOfWork.UserRepository(),
		unitOfWork.CacheRepository(),
	)
	
	gameService := game.NewGameService(
		unitOfWork.GameRepository(),
		unitOfWork.UserRepository(),
		unitOfWork.LeaderboardRepository(),
		unitOfWork.CacheRepository(),
		5,  // max workers
		50, // queue size
	)
	
	leaderboardSvc := leaderboard.NewLeaderboardService(
		unitOfWork.LeaderboardRepository(),
		unitOfWork.UserRepository(),
		unitOfWork.CacheRepository(),
		3600, // cache TTL in seconds
	)
	
	// Create context
	ctx := context.Background()
	
	// Example 1: User Registration and Authentication
	fmt.Println("=== User Registration and Authentication ===")
	
	// Register users
	user1, err := authService.Register(ctx, &auth.RegisterRequest{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Failed to register user1: %v", err)
	}
	
	user2, err := authService.Register(ctx, &auth.RegisterRequest{
		Username: "bob",
		Email:    "bob@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Failed to register user2: %v", err)
	}
	
	fmt.Printf("Registered users: %s, %s\n", user1.Username, user2.Username)
	
	// Login user
	session, err := authService.Login(ctx, &auth.LoginRequest{
		Username: "alice",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	
	fmt.Printf("User logged in: %s (Session: %s)\n", session.Username, session.ID)
	
	// Example 2: Game Creation and Management
	fmt.Println("\n=== Game Creation and Management ===")
	
	// Create a game
	game, err := gameService.CreateGame(ctx, user1.ID, user2.ID)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	
	fmt.Printf("Created game: %s between %s and %s\n", game.ID, user1.Username, user2.Username)
	
	// Start the game
	err = gameService.StartGame(ctx, game.ID)
	if err != nil {
		log.Fatalf("Failed to start game: %v", err)
	}
	
	fmt.Println("Game started successfully")
	
	// Update scores
	err = gameService.UpdateScore(ctx, game.ID, user1.ID, 150)
	if err != nil {
		log.Fatalf("Failed to update score: %v", err)
	}
	
	err = gameService.UpdateScore(ctx, game.ID, user2.ID, 120)
	if err != nil {
		log.Fatalf("Failed to update score: %v", err)
	}
	
	fmt.Printf("Updated scores - %s: 150, %s: 120\n", user1.Username, user2.Username)
	
	// End the game
	result, err := gameService.EndGame(ctx, game.ID)
	if err != nil {
		log.Fatalf("Failed to end game: %v", err)
	}
	
	fmt.Printf("Game ended - Winner: %s, Duration: %v\n", result.WinnerID, result.Duration)
	
	// Example 3: Leaderboard Management
	fmt.Println("\n=== Leaderboard Management ===")
	
	// Create a leaderboard
	lb, err := leaderboardSvc.CreateLeaderboard(ctx, "Global Leaderboard", models.LeaderboardTypeGlobal, 100)
	if err != nil {
		log.Fatalf("Failed to create leaderboard: %v", err)
	}
	
	fmt.Printf("Created leaderboard: %s\n", lb.Name)
	
	// Add scores to leaderboard
	err = leaderboardSvc.AddScore(ctx, lb.ID, user1.ID, 150)
	if err != nil {
		log.Fatalf("Failed to add score: %v", err)
	}
	
	err = leaderboardSvc.AddScore(ctx, lb.ID, user2.ID, 120)
	if err != nil {
		log.Fatalf("Failed to add score: %v", err)
	}
	
	fmt.Println("Added scores to leaderboard")
	
	// Get top entries
	entries, err := leaderboardSvc.GetTopEntries(ctx, lb.ID, 10)
	if err != nil {
		log.Fatalf("Failed to get top entries: %v", err)
	}
	
	fmt.Println("Top entries:")
	for _, entry := range entries {
		fmt.Printf("  %d. %s - %d points\n", entry.Rank, entry.Username, entry.Score)
	}
	
	// Get user rank
	rank, err := leaderboardSvc.GetUserRank(ctx, lb.ID, user1.ID)
	if err != nil {
		log.Fatalf("Failed to get user rank: %v", err)
	}
	
	fmt.Printf("%s's rank: %d\n", user1.Username, rank)
	
	// Example 4: User Statistics
	fmt.Println("\n=== User Statistics ===")
	
	// Get user stats
	stats, err := authService.GetUserStats(ctx, user1.ID)
	if err != nil {
		log.Fatalf("Failed to get user stats: %v", err)
	}
	
	fmt.Printf("%s's statistics:\n", user1.Username)
	fmt.Printf("  Total games: %d\n", stats.TotalGames)
	fmt.Printf("  Wins: %d\n", stats.Wins)
	fmt.Printf("  Losses: %d\n", stats.Losses)
	fmt.Printf("  Win rate: %.1f%%\n", stats.GetWinRate())
	fmt.Printf("  Average score: %.1f\n", stats.GetAverageScore())
	
	// Example 5: Concurrent Operations
	fmt.Println("\n=== Concurrent Operations ===")
	
	// Create multiple games concurrently
	gameIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		game, err := gameService.CreateGame(ctx, user1.ID, user2.ID)
		if err != nil {
			log.Fatalf("Failed to create game %d: %v", i+1, err)
		}
		gameIDs[i] = game.ID
	}
	
	fmt.Printf("Created %d games concurrently\n", len(gameIDs))
	
	// Start all games
	for _, gameID := range gameIDs {
		err := gameService.StartGame(ctx, gameID)
		if err != nil {
			log.Printf("Failed to start game %s: %v", gameID, err)
		}
	}
	
	fmt.Println("Started all games")
	
	// Update scores concurrently
	for i, gameID := range gameIDs {
		score := int64(100 + i*50)
		err := gameService.UpdateScore(ctx, gameID, user1.ID, score)
		if err != nil {
			log.Printf("Failed to update score for game %s: %v", gameID, err)
		}
	}
	
	fmt.Println("Updated scores concurrently")
	
	// End all games
	for _, gameID := range gameIDs {
		_, err := gameService.EndGame(ctx, gameID)
		if err != nil {
			log.Printf("Failed to end game %s: %v", gameID, err)
		}
	}
	
	fmt.Println("Ended all games")
	
	// Example 6: Error Handling
	fmt.Println("\n=== Error Handling ===")
	
	// Try to create a game with invalid players
	_, err = gameService.CreateGame(ctx, "invalid-user", user2.ID)
	if err != nil {
		fmt.Printf("Expected error when creating game with invalid user: %v\n", err)
	}
	
	// Try to add invalid score
	err = leaderboardSvc.AddScore(ctx, lb.ID, user1.ID, -50)
	if err != nil {
		fmt.Printf("Expected error when adding negative score: %v\n", err)
	}
	
	// Try to get non-existent user
	_, err = authService.GetUserStats(ctx, "non-existent-user")
	if err != nil {
		fmt.Printf("Expected error when getting stats for non-existent user: %v\n", err)
	}
	
	// Example 7: Cleanup
	fmt.Println("\n=== Cleanup ===")
	
	// Logout user
	err = authService.Logout(ctx, session.ID)
	if err != nil {
		log.Printf("Failed to logout: %v", err)
	}
	
	fmt.Println("User logged out")
	
	// Close services
	err = gameService.Close()
	if err != nil {
		log.Printf("Failed to close game service: %v", err)
	}
	
	leaderboardSvc.Close()
	
	err = unitOfWork.Close()
	if err != nil {
		log.Printf("Failed to close unit of work: %v", err)
	}
	
	fmt.Println("All services closed successfully")
	
	fmt.Println("\n=== Example Completed Successfully ===")
}

// Helper function to demonstrate context usage
func demonstrateContext(ctx context.Context) {
	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	// Simulate some work
	select {
	case <-timeoutCtx.Done():
		fmt.Println("Context cancelled or timed out")
	case <-time.After(2 * time.Second):
		fmt.Println("Work completed successfully")
	}
}

// Helper function to demonstrate error wrapping
func demonstrateErrorWrapping() error {
	// Simulate an error
	originalErr := fmt.Errorf("database connection failed")
	
	// Wrap the error with context
	wrappedErr := fmt.Errorf("failed to create user: %w", originalErr)
	
	return wrappedErr
}
