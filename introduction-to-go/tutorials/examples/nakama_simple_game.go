//go:build nakama

// Package main - Nakama Game Server Module
// This file contains a complete multiplayer game implementation using Nakama
// It demonstrates RPC functions, hooks, storage, leaderboards, and real-time game logic
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

// =============================================================================
// GAME MODELS AND DATA STRUCTURES
// =============================================================================

// Player represents a user in the game
// This struct defines what data we store for each player
type Player struct {
	UserID    string   `json:"user_id"`    // Unique identifier for the user
	Username  string   `json:"username"`   // Display name of the player
	Position  Position `json:"position"`   // Current position in the game world
	Score     int      `json:"score"`      // Player's current score
	Health    int      `json:"health"`     // Player's health (0-100)
	Ready     bool     `json:"ready"`      // Whether player is ready to start
}

// Position represents 2D coordinates in the game world
// Used for tracking player movement and positioning
type Position struct {
	X float64 `json:"x"` // X coordinate (horizontal position)
	Y float64 `json:"y"` // Y coordinate (vertical position)
}

// GameState represents the complete state of a game session
// This is stored in Nakama's storage system and shared between all players
type GameState struct {
	GameID      string             `json:"game_id"`      // Unique identifier for this game
	Players     map[string]*Player `json:"players"`      // Map of user_id -> Player data
	GameStatus  string             `json:"game_status"`  // "waiting", "playing", "finished"
	StartTime   time.Time          `json:"start_time"`   // When the game started
	EndTime     time.Time          `json:"end_time"`     // When the game ended
	MaxPlayers  int                `json:"max_players"`  // Maximum players allowed
	MinPlayers  int                `json:"min_players"`  // Minimum players needed to start
}

// =============================================================================
// RPC (REMOTE PROCEDURE CALL) FUNCTIONS
// These are server functions that clients can call remotely
// =============================================================================

// CreateGameRPC - Creates a new game session
// Called by a client to start a new multiplayer game
// Parameters: max_players, min_players
// Returns: game_id and game state
func CreateGameRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the JSON request from the client
	var request struct {
		MaxPlayers int `json:"max_players"` // Maximum number of players allowed
		MinPlayers int `json:"min_players"` // Minimum players needed to start
	}
	
	// Convert JSON string to struct
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	// Get the user ID from the context (Nakama provides this automatically)
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	
	// Validate input parameters
	if request.MaxPlayers < 2 || request.MaxPlayers > 8 {
		return "", errors.New("max players must be between 2 and 8")
	}
	if request.MinPlayers < 2 || request.MinPlayers > request.MaxPlayers {
		return "", errors.New("invalid min players")
	}
	
	// Create a new game state
	gameState := &GameState{
		GameID:     generateGameID(),           // Generate unique game ID
		Players:    make(map[string]*Player),   // Initialize empty players map
		GameStatus: "waiting",                  // Game starts in waiting state
		MaxPlayers: request.MaxPlayers,
		MinPlayers: request.MinPlayers,
	}
	
	// Get user information from Nakama
	user, err := nk.UserGetId(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// Add the creator as the first player
	gameState.Players[userID] = &Player{
		UserID:   userID,
		Username: user.Username,
		Position: Position{X: 0, Y: 0}, // Start at origin
		Score:    0,
		Health:   100, // Full health
		Ready:    false,
	}
	
	// Save game state to Nakama's storage system
	// This makes the game persistent and accessible to other players
	gameStateBytes, _ := json.Marshal(gameState)
	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "games",        // Storage collection name
			Key:        gameState.GameID, // Unique key for this game
			UserID:     "",             // Empty = global storage (not user-specific)
			Value:      string(gameStateBytes), // JSON string of game state
		},
	})
	
	if err != nil {
		return "", err
	}
	
	// Return success response to client
	response := map[string]interface{}{
		"success": true,
		"game_id": gameState.GameID,
		"game":    gameState,
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// JoinGameRPC - Join an existing game
// Called by a client to join a game created by another player
// Parameters: game_id
// Returns: updated game state
func JoinGameRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the request
	var request struct {
		GameID string `json:"game_id"` // ID of the game to join
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	// Get authenticated user ID
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	
	// Read the game state from storage
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "games",
			Key:        request.GameID,
			UserID:     "", // Global storage
		},
	})
	
	if err != nil {
		return "", err
	}
	
	// Check if game exists
	if len(objects) == 0 {
		return "", errors.New("game not found")
	}
	
	// Parse the stored game state
	var gameState GameState
	if err := json.Unmarshal([]byte(objects[0].Value), &gameState); err != nil {
		return "", err
	}
	
	// Validate join conditions
	if len(gameState.Players) >= gameState.MaxPlayers {
		return "", errors.New("game is full")
	}
	
	if gameState.GameStatus == "playing" {
		return "", errors.New("game already in progress")
	}
	
	// Check if player is already in the game
	if _, exists := gameState.Players[userID]; exists {
		return "", errors.New("already in game")
	}
	
	// Get user information
	user, err := nk.UserGetId(ctx, userID)
	if err != nil {
		return "", err
	}
	
	// Add player to the game
	gameState.Players[userID] = &Player{
		UserID:   userID,
		Username: user.Username,
		Position: Position{X: 0, Y: 0},
		Score:    0,
		Health:   100,
		Ready:    false,
	}
	
	// Save the updated game state
	gameStateBytes, _ := json.Marshal(gameState)
	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "games",
			Key:        gameState.GameID,
			UserID:     "",
			Value:      string(gameStateBytes),
		},
	})
	
	if err != nil {
		return "", err
	}
	
	// Return success response
	response := map[string]interface{}{
		"success": true,
		"game":    gameState,
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// StartGameRPC - Start the game if enough players are ready
// Called when players want to begin the actual gameplay
// Parameters: game_id
// Returns: updated game state
func StartGameRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	
	// Get current game state
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "games",
			Key:        request.GameID,
			UserID:     "",
		},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(objects) == 0 {
		return "", errors.New("game not found")
	}
	
	var gameState GameState
	if err := json.Unmarshal([]byte(objects[0].Value), &gameState); err != nil {
		return "", err
	}
	
	// Validate user is in the game
	if _, exists := gameState.Players[userID]; !exists {
		return "", errors.New("not in game")
	}
	
	// Check if we have enough players
	if len(gameState.Players) < gameState.MinPlayers {
		return "", errors.New("not enough players")
	}
	
	// Check if all players are ready
	readyCount := 0
	for _, player := range gameState.Players {
		if player.Ready {
			readyCount++
		}
	}
	
	if readyCount < len(gameState.Players) {
		return "", errors.New("not all players are ready")
	}
	
	// Start the game
	gameState.GameStatus = "playing"
	gameState.StartTime = time.Now()
	
	// Save the updated state
	gameStateBytes, _ := json.Marshal(gameState)
	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "games",
			Key:        gameState.GameID,
			UserID:     "",
			Value:      string(gameStateBytes),
		},
	})
	
	if err != nil {
		return "", err
	}
	
	response := map[string]interface{}{
		"success": true,
		"game":    gameState,
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// UpdatePlayerPositionRPC - Update player position during game
// Called frequently during gameplay to track player movement
// Parameters: game_id, position (x, y coordinates)
// Returns: updated position
func UpdatePlayerPositionRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request struct {
		GameID   string   `json:"game_id"`
		Position Position `json:"position"` // New position coordinates
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	
	// Get current game state
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "games",
			Key:        request.GameID,
			UserID:     "",
		},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(objects) == 0 {
		return "", errors.New("game not found")
	}
	
	var gameState GameState
	if err := json.Unmarshal([]byte(objects[0].Value), &gameState); err != nil {
		return "", err
	}
	
	// Only allow position updates during active gameplay
	if gameState.GameStatus != "playing" {
		return "", errors.New("game not in progress")
	}
	
	// Find the player in the game
	player, exists := gameState.Players[userID]
	if !exists {
		return "", errors.New("not in game")
	}
	
	// Update the player's position
	player.Position = request.Position
	
	// Save the updated game state
	gameStateBytes, _ := json.Marshal(gameState)
	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "games",
			Key:        gameState.GameID,
			UserID:     "",
			Value:      string(gameStateBytes),
		},
	})
	
	if err != nil {
		return "", err
	}
	
	response := map[string]interface{}{
		"success": true,
		"position": player.Position,
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// EndGameRPC - End the game and save results
// Called when the game is finished to save scores and clean up
// Parameters: game_id
// Returns: final game state
func EndGameRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("user not authenticated")
	}
	
	// Get current game state
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "games",
			Key:        request.GameID,
			UserID:     "",
		},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(objects) == 0 {
		return "", errors.New("game not found")
	}
	
	var gameState GameState
	if err := json.Unmarshal([]byte(objects[0].Value), &gameState); err != nil {
		return "", err
	}
	
	// Validate user is in the game
	if _, exists := gameState.Players[userID]; !exists {
		return "", errors.New("not in game")
	}
	
	// Mark game as finished
	gameState.GameStatus = "finished"
	gameState.EndTime = time.Now()
	
	// Save all player scores to the leaderboard
	// This creates competitive rankings across all games
	for playerID, player := range gameState.Players {
		_, err := nk.LeaderboardRecordWrite(ctx, "game_scores", playerID, int64(player.Score), nil, nil)
		if err != nil {
			logger.Error("Failed to save score for player %s: %v", playerID, err)
		}
	}
	
	// Save the final game state
	gameStateBytes, _ := json.Marshal(gameState)
	_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
		{
			Collection: "games",
			Key:        gameState.GameID,
			UserID:     "",
			Value:      string(gameStateBytes),
		},
	})
	
	if err != nil {
		return "", err
	}
	
	response := map[string]interface{}{
		"success": true,
		"game":    gameState,
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// GetGameStateRPC - Get current game state
// Called by clients to sync their local state with the server
// Parameters: game_id
// Returns: current game state
func GetGameStateRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request struct {
		GameID string `json:"game_id"`
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	// Read game state from storage
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "games",
			Key:        request.GameID,
			UserID:     "",
		},
	})
	
	if err != nil {
		return "", err
	}
	
	if len(objects) == 0 {
		return "", errors.New("game not found")
	}
	
	// Return the raw game state JSON
	response := map[string]interface{}{
		"success": true,
		"game":    objects[0].Value, // Raw JSON string
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// GetLeaderboardRPC - Get game leaderboard
// Called to show competitive rankings
// Parameters: limit (optional, defaults to 10)
// Returns: leaderboard records
func GetLeaderboardRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var request struct {
		Limit int `json:"limit"` // Number of records to return
	}
	
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return "", err
	}
	
	// Default to 10 records if not specified
	if request.Limit == 0 {
		request.Limit = 10
	}
	
	// Get leaderboard records from Nakama
	// This returns the top scores from all games
	records, err := nk.LeaderboardRecordsList(ctx, "game_scores", []string{}, request.Limit, "")
	if err != nil {
		return "", err
	}
	
	response := map[string]interface{}{
		"success": true,
		"records": records.Records,       // Global leaderboard
		"owner_records": records.OwnerRecords, // Current user's records
	}
	
	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

// =============================================================================
// HOOKS - AUTOMATIC SERVER EVENTS
// These functions are called automatically by Nakama when certain events occur
// =============================================================================

// AfterAuthenticateHook - Track user login and create profile if needed
// This hook runs automatically every time a user logs in
// It ensures every user has a profile and tracks login activity
func AfterAuthenticateHook(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateRequest) error {
	userID := out.UserId
	
	// Check if user already has a profile
	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "user_profiles", // User-specific storage collection
			Key:        "profile",       // Key for the profile data
			UserID:     userID,          // This user's storage
		},
	})
	
	if err != nil {
		logger.Error("Failed to read user profile: %v", err)
		return err
	}
	
	// Create profile if it doesn't exist
	if len(objects) == 0 {
		// Initialize new user profile with default values
		profile := map[string]interface{}{
			"user_id":      userID,
			"games_played": 0,           // Track total games
			"total_score":  0,           // Track total score
			"best_score":   0,           // Track best score
			"created_at":   time.Now(),  // When profile was created
			"last_login":   time.Now(),  // Last login time
		}
		
		// Save the new profile
		profileBytes, _ := json.Marshal(profile)
		_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
			{
				Collection: "user_profiles",
				Key:        "profile",
				UserID:     userID,
				Value:      string(profileBytes),
			},
		})
		
		if err != nil {
			logger.Error("Failed to create user profile: %v", err)
			return err
		}
		
		logger.Info("Created profile for user: %s", userID)
	} else {
		// Update last login time for existing users
		var profile map[string]interface{}
		if err := json.Unmarshal([]byte(objects[0].Value), &profile); err == nil {
			profile["last_login"] = time.Now()
			profileBytes, _ := json.Marshal(profile)
			_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
				{
					Collection: "user_profiles",
					Key:        "profile",
					UserID:     userID,
					Value:      string(profileBytes),
				},
			})
		}
	}
	
	return nil
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// generateGameID - Creates a unique identifier for games
// Uses timestamp to ensure uniqueness
func generateGameID() string {
	return fmt.Sprintf("game_%d", time.Now().UnixNano())
}

// =============================================================================
// MODULE INITIALIZATION
// This function is called when Nakama loads your module
// =============================================================================

// InitModule - Initialize the game module
// This is the main entry point that registers all RPC functions and hooks
// Called automatically by Nakama when the server starts
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// Register all RPC functions so clients can call them
	initializer.RegisterRpc("create_game", CreateGameRPC)           // Create new game
	initializer.RegisterRpc("join_game", JoinGameRPC)               // Join existing game
	initializer.RegisterRpc("start_game", StartGameRPC)             // Start gameplay
	initializer.RegisterRpc("update_position", UpdatePlayerPositionRPC) // Update player position
	initializer.RegisterRpc("end_game", EndGameRPC)                 // End game
	initializer.RegisterRpc("get_game_state", GetGameStateRPC)      // Get current state
	initializer.RegisterRpc("get_leaderboard", GetLeaderboardRPC)   // Get rankings
	
	// Register hooks for automatic events
	initializer.RegisterAfterAuthenticate(AfterAuthenticateHook)    // User login tracking
	
	// Create the leaderboard for storing scores
	// This creates a global leaderboard that persists across server restarts
	_, err := nk.LeaderboardCreate(ctx, "game_scores", "Game High Scores", "desc", "best", "0 0 * * 1", false, "")
	if err != nil {
		logger.Error("Failed to create leaderboard: %v", err)
		return err
	}
	
	logger.Info("Simple Game Module initialized successfully!")
	return nil
}
