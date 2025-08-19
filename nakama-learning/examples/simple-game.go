//go:build nakama

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

// Simple game data structure
type Game struct {
	ID      string            `json:"id"`
	Players map[string]Player `json:"players"`
	State   string            `json:"state"` // "waiting", "playing", "finished"
}

type Player struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Score  int    `json:"score"`
	Health int    `json:"health"`
}

// Store active games in memory
var games = make(map[string]*Game)

// This function is called when Nakama starts up
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("🎮 Game module loaded!")

	// Register our game functions so players can call them
	initializer.RegisterRpc("create_game", createGame)
	initializer.RegisterRpc("join_game", joinGame)
	initializer.RegisterRpc("submit_score", submitScore)
	initializer.RegisterRpc("get_leaderboard", getLeaderboard)

	// Create a leaderboard for scores
	nk.LeaderboardCreate(ctx, "game_scores", "desc", "best", "weekly", 0, false)

	logger.Info("✅ All game functions registered")
	return nil
}

// Function 1: Create a new game
func createGame(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Get the player's info
	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	// Create a new game
	gameID := fmt.Sprintf("game_%s", userID)
	game := &Game{
		ID: gameID,
		Players: map[string]Player{
			userID: {
				ID:     userID,
				Name:   username,
				Score:  0,
				Health: 100,
			},
		},
		State: "waiting",
	}

	// Save the game
	games[gameID] = game

	logger.Info("🎮 Game created: %s by %s", gameID, username)

	// Send response back to player
	response := map[string]interface{}{
		"success": true,
		"game_id": gameID,
		"message": "Game created!",
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// Function 2: Join an existing game
func joinGame(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the request
	var request struct {
		GameID string `json:"game_id"`
	}
	json.Unmarshal([]byte(payload), &request)

	// Get player info
	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	// Find the game
	game, exists := games[request.GameID]
	if !exists {
		return "", fmt.Errorf("game not found")
	}

	// Add player to game
	game.Players[userID] = Player{
		ID:     userID,
		Name:   username,
		Score:  0,
		Health: 100,
	}

	// Start game if we have 2+ players
	if len(game.Players) >= 2 {
		game.State = "playing"
	}

	logger.Info("👥 %s joined game %s", username, request.GameID)

	response := map[string]interface{}{
		"success": true,
		"message": "Joined game!",
		"players": len(game.Players),
		"state":   game.State,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// Function 3: Submit a score
func submitScore(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Parse the score
	var request struct {
		Score int `json:"score"`
	}
	json.Unmarshal([]byte(payload), &request)

	// Get player info
	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	username := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)

	// Save score to leaderboard
	nk.LeaderboardRecordWrite(ctx, "game_scores", userID, username, int64(request.Score), nil, nil)

	logger.Info("🏆 %s scored %d points", username, request.Score)

	response := map[string]interface{}{
		"success": true,
		"message": "Score saved!",
		"score":   request.Score,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// Function 4: Get leaderboard
func getLeaderboard(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Get top 10 scores
	records, _, _ := nk.LeaderboardRecordsList(ctx, "game_scores", []string{}, 10, 10)

	// Format the response
	var leaderboard []map[string]interface{}
	for _, record := range records {
		leaderboard = append(leaderboard, map[string]interface{}{
			"rank":   record.Rank,
			"player": record.Username,
			"score":  record.Score,
		})
	}

	response := map[string]interface{}{
		"success":     true,
		"leaderboard": leaderboard,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// Called when a new player registers
func AfterAuthenticate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, session *api.Session, username string, created bool) {
	if created {
		logger.Info("🎉 New player: %s", username)
	} else {
		logger.Info("👋 Player logged in: %s", username)
	}
}
