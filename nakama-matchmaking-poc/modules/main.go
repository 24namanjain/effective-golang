package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

// MatchState represents the state of a match
// This is like a "game room" that tracks the state of the match
type MatchState struct {
	Players map[string]*Player `json:"players"` // map of user_id to player
	State   string             `json:"state"`   // "waiting", "playing", "finished"
	Created time.Time          `json:"created"`
}

// Player represents a player in the match
type Player struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Ready    bool   `json:"ready"`
	Score    int    `json:"score"`
}

// MatchData represents match information
type MatchData struct {
	MatchID string    `json:"match_id"`
	Players []*Player `json:"players"`
	State   string    `json:"state"`
	Created time.Time `json:"created"`
}

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	// Register match handler
	initializer.RegisterMatch("matchmaking_match", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return &MatchmakingMatch{}, nil
	})

	// Register RPC functions
	initializer.RegisterRpc("get_match_info", GetMatchInfoRPC)
	initializer.RegisterRpc("player_ready", PlayerReadyRPC)
	initializer.RegisterRpc("get_matchmaking_stats", GetMatchmakingStatsRPC)

	logger.Info("Matchmaking POC Module initialized successfully!")
	return nil
}

// MatchmakingMatch implements the runtime.Match interface
type MatchmakingMatch struct{}

func (m *MatchmakingMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("Match initialized")

	state := &MatchState{
		Players: make(map[string]*Player),
		State:   "waiting",
		Created: time.Now(),
	}

	return state, 2, "" // min 2 players, no label
}

func (m *MatchmakingMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	matchState := state.(*MatchState)

	logger.Info("Player attempting to join match", "user_id", presence.GetUserId(), "username", presence.GetUsername())

	// Check if player is already in the match
	if _, exists := matchState.Players[presence.GetUserId()]; exists {
		return matchState, true, ""
	}

	// Add player to match
	player := &Player{
		UserID:   presence.GetUserId(),
		Username: presence.GetUsername(),
		Ready:    false,
		Score:    0,
	}

	matchState.Players[presence.GetUserId()] = player

	// Notify all players about the new player
	dispatcher.BroadcastMessage(
		1, // priority
		// message to send
		[]byte(fmt.Sprintf(`{"type":"player_joined","player":{"user_id":"%s","username":"%s"}}`, presence.GetUserId(), presence.GetUsername())),
		nil, // sender user_id
		nil, // sender username
		// send to all players
		true,
	)

	return matchState, true, ""
}

func (m *MatchmakingMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	matchState := state.(*MatchState)

	for _, presence := range presences {
		logger.Info("Player joined match", "user_id", presence.GetUserId(), "username", presence.GetUsername())

		// Notify all players about the join
		dispatcher.BroadcastMessage(1, []byte(fmt.Sprintf(`{"type":"player_joined","player":{"user_id":"%s","username":"%s"}}`, presence.GetUserId(), presence.GetUsername())), nil, nil, true)
	}

	return matchState
}

func (m *MatchmakingMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	matchState := state.(*MatchState)

	for _, presence := range presences {
		logger.Info("Player left match", "user_id", presence.GetUserId(), "username", presence.GetUsername())

		// Remove player from match
		delete(matchState.Players, presence.GetUserId())

		// Notify remaining players
		dispatcher.BroadcastMessage(1, []byte(fmt.Sprintf(`{"type":"player_left","user_id":"%s","username":"%s"}`, presence.GetUserId(), presence.GetUsername())), nil, nil, true)
	}

	return matchState
}

/*
*
This is the main loop that runs every tick.
It processes incoming messages from players.
*/
func (m *MatchmakingMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	matchState := state.(*MatchState)

	// Process incoming messages
	for _, message := range messages {
		userID := message.GetUserId()
		username := message.GetUsername()

		// Unmarshal the message data
		var data map[string]interface{}
		if err := json.Unmarshal(message.GetData(), &data); err != nil {
			logger.Error("Failed to unmarshal message data", "error", err)
			continue
		}

		// Get the message type
		messageType, ok := data["type"].(string)
		if !ok {
			continue
		}

		switch messageType {
		// Player marked as ready
		case "player_ready":
			// Check if the player exists in the match
			if player, exists := matchState.Players[userID]; exists {
				player.Ready = true
				logger.Info("Player marked as ready", "user_id", userID, "username", username)

				// Notify all players
				dispatcher.BroadcastMessage(
					1,
					[]byte(fmt.Sprintf(`{"type":"player_ready","user_id":"%s","username":"%s"}`, userID, username)),
					nil,
					nil,
					true,
				)

				// Check if all players are ready and we have enough players
				if len(matchState.Players) >= 2 && allPlayersReady(matchState) {
					// Set the match state to playing
					matchState.State = "playing"
					logger.Info("Match started - all players ready")

					// Notify all players that match has started
					dispatcher.BroadcastMessage(
						1,
						[]byte(`{"type":"match_started"}`),
						nil,
						nil,
						true,
					)
				}
			}

		// Player score updated
		case "player_score":
			if player, exists := matchState.Players[userID]; exists {
				if score, ok := data["score"].(float64); ok {
					player.Score = int(score)
					logger.Info("Player score updated", "user_id", userID, "score", player.Score)

					// Notify all players about score update
					dispatcher.BroadcastMessage(1, []byte(fmt.Sprintf(`{"type":"score_update","user_id":"%s","score":%d}`, userID, player.Score)), nil, nil, true)
				}
			}
		}
	}

	return matchState
}

func (m *MatchmakingMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	matchState := state.(*MatchState)

	logger.Info("Match terminating", "grace_seconds", graceSeconds)

	// Notify all players that match is ending
	dispatcher.BroadcastMessage(1, []byte(`{"type":"match_ending"}`), nil, nil, true)

	return matchState
}

// MatchSignal is called when a signal is sent to the match
// It is used to send data to the match
func (m *MatchmakingMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, ""
}

// Helper function to check if all players are ready
func allPlayersReady(state *MatchState) bool {
	for _, player := range state.Players {
		if !player.Ready {
			return false
		}
	}
	return true
}

// RPC Functions

func GetMatchInfoRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	// Get user's current match
	matches, err := nk.MatchList(ctx, 1, true, "", nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to list matches: %v", err)
	}

	response := map[string]interface{}{
		"user_id": userID,
		"matches": matches,
	}

	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

func PlayerReadyRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, _ := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	response := map[string]interface{}{
		"user_id": userID,
		"status":  "ready",
		"message": "Player marked as ready",
	}

	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}

func GetMatchmakingStatsRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// Get active matches
	matches, err := nk.MatchList(ctx, 10, true, "", nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to list matches: %v", err)
	}

	response := map[string]interface{}{
		"active_matches": len(matches),
		"timestamp":      time.Now().Unix(),
	}

	responseBytes, _ := json.Marshal(response)
	return string(responseBytes), nil
}
