package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"effective-golang/internal/auth"
	"effective-golang/internal/game"
	"effective-golang/internal/leaderboard"
	"effective-golang/internal/models"
	"effective-golang/pkg/utils"
)

// Auth handlers

func registerHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		user, err := authService.Register(r.Context(), &req)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.CreatedResponse(w, user)
	}
}

func loginHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req auth.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		session, err := authService.Login(r.Context(), &req)
		if err != nil {
			utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		
		utils.SuccessResponse(w, session)
	}
}

func logoutHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Authorization")
		if sessionID == "" {
			utils.ErrorResponse(w, http.StatusUnauthorized, "No session provided")
			return
		}
		
		if err := authService.Logout(r.Context(), sessionID); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		utils.SuccessResponse(w, map[string]string{"message": "Logged out successfully"})
	}
}

// Game handlers

func createGameHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Player1ID string `json:"player1_id"`
			Player2ID string `json:"player2_id"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		game, err := gameService.CreateGame(r.Context(), req.Player1ID, req.Player2ID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.CreatedResponse(w, game)
	}
}

func startGameHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["gameID"]
		
		if err := gameService.StartGame(r.Context(), gameID); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.SuccessResponse(w, map[string]string{"message": "Game started successfully"})
	}
}

func updateScoreHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["gameID"]
		
		var req struct {
			PlayerID string `json:"player_id"`
			Score    int64  `json:"score"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		if err := gameService.UpdateScore(r.Context(), gameID, req.PlayerID, req.Score); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.SuccessResponse(w, map[string]string{"message": "Score updated successfully"})
	}
}

func endGameHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["gameID"]
		
		result, err := gameService.EndGame(r.Context(), gameID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.SuccessResponse(w, result)
	}
}

func getActiveGamesHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		games, err := gameService.GetActiveGames(r.Context())
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		utils.SuccessResponse(w, games)
	}
}

func getGameHandler(gameService *game.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["gameID"]
		
		game, err := gameService.GetGame(r.Context(), gameID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, game)
	}
}

// Leaderboard handlers

func createLeaderboardHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name        string                    `json:"name"`
			Type        models.LeaderboardType    `json:"type"`
			MaxEntries  int                       `json:"max_entries"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		leaderboard, err := leaderboardSvc.CreateLeaderboard(r.Context(), req.Name, req.Type, req.MaxEntries)
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.CreatedResponse(w, leaderboard)
	}
}

func addScoreHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardID := vars["leaderboardID"]
		
		var req struct {
			UserID string `json:"user_id"`
			Score  int64  `json:"score"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		
		if err := leaderboardSvc.AddScore(r.Context(), leaderboardID, req.UserID, req.Score); err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		utils.SuccessResponse(w, map[string]string{"message": "Score added successfully"})
	}
}

func getTopEntriesHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardID := vars["leaderboardID"]
		
		countStr := r.URL.Query().Get("count")
		count := 10 // default
		if countStr != "" {
			if parsed, err := strconv.Atoi(countStr); err == nil {
				count = parsed
			}
		}
		
		entries, err := leaderboardSvc.GetTopEntries(r.Context(), leaderboardID, count)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, entries)
	}
}

func getUserRankHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardID := vars["leaderboardID"]
		userID := vars["userID"]
		
		rank, err := leaderboardSvc.GetUserRank(r.Context(), leaderboardID, userID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, map[string]interface{}{
			"user_id": userID,
			"rank":    rank,
		})
	}
}

func getLeaderboardStatsHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardID := vars["leaderboardID"]
		
		stats, err := leaderboardSvc.GetStats(r.Context(), leaderboardID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, stats)
	}
}

func getLeaderboardHandler(leaderboardSvc *leaderboard.LeaderboardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		leaderboardID := vars["leaderboardID"]
		
		leaderboard, err := leaderboardSvc.GetLeaderboard(r.Context(), leaderboardID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, leaderboard)
	}
}

// User handlers

func getUserStatsHandler(authService *auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userID"]
		
		stats, err := authService.GetUserStats(r.Context(), userID)
		if err != nil {
			utils.ErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		
		utils.SuccessResponse(w, stats)
	}
}
